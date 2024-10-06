// Package mapper provides utilities for mapping between different types using reflection.
// It supports creating mappings for struct types and applying custom mapping functions.
package mapper

import (
	"flag"
	"fmt"
	"reflect"

	reflectionHelper "github.com/NekKkMirror/go-app/internal/pkg/reflection/reflection-helper"
	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Error definitions for the mapper package.
var (
	ErrNilFunction = errors.New("mapper: nil function")

	ErrMapNotExist = errors.New("mapper: map does not exist")

	ErrMapAlreadyExists = errors.New("mapper: map already exists")

	ErrUnsupportedMap = errors.New("mapper: unsupported map")
)

// Constants for indexing source and destination keys.
const (
	SrcKeyIndex = iota
	DestKeyIndex
)

// Config holds configuration options for the mapper.
type Config struct {
	MapUnexportedFields bool // Determines if unexported fields should be mapped.
}

// mappingsEntry represents a mapping between source and destination types.
type mappingsEntry struct {
	SourceType      reflect.Type
	DestinationType reflect.Type
}

// typeMeta holds metadata for mapping keys to tags and vice versa.
type typeMeta struct {
	keysToTags map[string]string
	tagsToKeys map[string]string
}

// mapFunc defines a function type for custom mapping logic.
type mapFunc[TSrc any, TDst any] func(TSrc) TDst

// profiles and maps store mapping profiles and functions.
var profiles = map[string][][2]string{}
var maps = map[mappingsEntry]interface{}{}
var mapperConfig *Config

// init initializes the default mapper configuration.
func init() {
	mapperConfig = &Config{
		MapUnexportedFields: false,
	}
}

// Configure sets the mapper configuration.
func Configure(config *Config) {
	mapperConfig = config
}

// CreateMap creates a mapping between source and destination types.
// It returns an error if the mapping already exists or if the types are unsupported.
func CreateMap[TSrc any, TDst any]() error {
	var src TSrc
	var dst TDst
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	if (srcType.Kind() != reflect.Struct && !(srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct)) ||
		(desType.Kind() != reflect.Struct && !(desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Struct)) {
		return ErrUnsupportedMap
	}

	if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
		pointerStructTypeKey := mappingsEntry{
			SourceType:      srcType,
			DestinationType: desType,
		}
		nonePointerStructTypeKey := mappingsEntry{
			SourceType:      srcType.Elem(),
			DestinationType: desType.Elem(),
		}

		if _, ok := maps[pointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}
		if _, ok := maps[nonePointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}

		maps[pointerStructTypeKey] = nil
		maps[nonePointerStructTypeKey] = nil
	} else {
		pointerStructTypeKey := mappingsEntry{
			SourceType:      srcType,
			DestinationType: desType,
		}
		nonePointerStructTypeKey := mappingsEntry{
			SourceType:      reflect.New(srcType).Type(),
			DestinationType: reflect.New(desType).Type(),
		}

		if _, ok := maps[pointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}
		if _, ok := maps[nonePointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}

		maps[pointerStructTypeKey] = nil
		maps[nonePointerStructTypeKey] = nil
	}

	if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
		srcType = srcType.Elem()
	}
	if desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Struct {
		desType = desType.Elem()
	}

	configProfile(srcType, desType)

	return nil
}

// configProfile generates a mapping profile between source and destination types.
// It iterates through the fields and methods of the source and destination types,
// and creates a mapping profile based on the field and method names.
//
// Parameters:
// - srcType: The reflect.Type of the source type.
// - desType: The reflect.Type of the destination type.
//
// Return:
// - None. The function modifies the global 'profiles' map directly.
func configProfile(srcType reflect.Type, desType reflect.Type) {
	flag.Parse()

	if srcType.Kind() != reflect.Struct {
		log.Errorf("expected reflect.Struct kind for type %s, but got %s", srcType.String(), srcType.Kind().String())
	}

	if desType.Kind() != reflect.Struct {
		log.Errorf("expected reflect.Struct kind for type %s, but got %s", desType.String(), desType.Kind().String())
	}

	var profile [][2]string

	srcMeta := getTypeMeta(srcType)
	destMeta := getTypeMeta(desType)
	srcMethods := getTypeMethods(srcType)

	for srcKey, srcTag := range srcMeta.keysToTags {
		if _, ok := destMeta.keysToTags[strcase.ToCamel(srcKey)]; ok {
			profile = append(profile, [2]string{srcKey, strcase.ToCamel(srcKey)})
		}

		if _, ok := destMeta.keysToTags[srcKey]; ok {
			profile = append(profile, [2]string{srcKey, srcKey})
			continue
		}

		if destKey, ok := destMeta.tagsToKeys[srcKey]; ok {
			profile = append(profile, [2]string{srcKey, destKey})
			continue
		}

		if _, ok := destMeta.keysToTags[srcTag]; ok {
			profile = append(profile, [2]string{srcKey, srcTag})
			continue
		}

		if destKey, ok := destMeta.tagsToKeys[srcTag]; ok {
			profile = append(profile, [2]string{srcKey, destKey})
			continue
		}

	}

	for _, method := range srcMethods {
		if _, ok := destMeta.keysToTags[method]; ok {
			profile = append(profile, [2]string{method, method})
			continue
		}
	}

	profiles[getProfileKey(srcType, desType)] = profile
}

// getTypeMeta generates a mapping between field names and their corresponding tags.
// It also generates a reverse mapping between tags and field names.
//
// Parameters:
// - val: The reflect.Type of the struct for which the mapping needs to be generated.
//
// Return:
// - typeMeta: A struct containing two maps: keysToTags and tagsToKeys.
//   - keysToTags: A map where the keys are field names and the values are their corresponding tags.
//   - tagsToKeys: A map where the keys are tags and the values are their corresponding field names.
func getTypeMeta(val reflect.Type) typeMeta {
	fieldsNum := val.NumField()

	keysToTags := make(map[string]string)
	tagsToKeys := make(map[string]string)

	for i := 0; i < fieldsNum; i++ {
		field := val.Field(i)
		fieldName := field.Name
		fieldTag := field.Tag.Get("mapper")

		keysToTags[fieldName] = fieldTag
		if fieldTag != "" {
			tagsToKeys[fieldTag] = fieldName
		}
	}

	return typeMeta{
		keysToTags: keysToTags,
		tagsToKeys: tagsToKeys,
	}
}

// getTypeMethods retrieves the names of all methods defined for a given struct type.
//
// Parameters:
// - val: reflect.Type representing the struct type for which the method names need to be retrieved.
//
// Return:
// - []string: A slice of strings containing the names of all methods defined for the given struct type.
//
// The function iterates through all methods of the given struct type and appends their names to a slice.
// Finally, it returns the slice containing the names of all methods.
func getTypeMethods(val reflect.Type) []string {
	methodsNum := val.NumMethod()
	var keys []string

	for i := 0; i < methodsNum; i++ {
		methodName := val.Method(i).Name
		keys = append(keys, methodName)
	}

	return keys
}

func getProfileKey(srcType reflect.Type, desType reflect.Type) string {
	return fmt.Sprintf("%s_%s", srcType.Name(), desType.Name())
}

// Map performs the mapping from the source to the destination type.
// It returns the mapped destination type and an error if the mapping does not exist.
func Map[TSrc any, TDes any](src TSrc) (TDes, error) {
	var des TDes
	srcType := reflect.TypeOf(src)
	desType := reflect.TypeOf(des)
	srcIsArray := false
	desIsArray := false

	if srcType.Kind() == reflect.Array || (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Array) || srcType.Kind() == reflect.Slice || (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Slice) {
		srcType = srcType.Elem()
		srcIsArray = true
	}

	if desType.Kind() == reflect.Array || (desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Array) || desType.Kind() == reflect.Slice || (desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Slice) {
		desType = desType.Elem()
		desIsArray = true
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	fn, ok := maps[k]
	if !ok {
		return *new(TDes), ErrMapNotExist
	}
	if fn != nil {
		fnReflect := reflect.ValueOf(fn)

		if desIsArray && srcIsArray {
			linq.From(src).Select(func(x interface{}) interface{} {
				return fnReflect.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
			}).ToSlice(&des)
		} else {
			return fnReflect.Call([]reflect.Value{reflect.ValueOf(src)})[0].Interface().(TDes), nil
		}

		desTypeValue := reflect.ValueOf(&des).Elem()
		err := processValues[TDes, TSrc](reflect.ValueOf(src), desTypeValue)
		if err != nil {
			return *new(TDes), err
		}
	}

	return des, nil
}

// processValues is a helper function that performs the mapping from the source to the destination type.
// It handles different kinds of source and destination values, including structs, slices, maps, and pointers.
//
// Parameters:
// - src: The reflect.Value representing the source value to be mapped.
// - dest: The reflect.Value representing the destination value where the mapped value will be stored.
//
// Return:
// - error: An error if the mapping fails, or nil if the mapping is successful.
func processValues[TDes any, TSrc any](src reflect.Value, dest reflect.Value) error {
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	if dest.Kind() == reflect.Interface {
		dest = dest.Elem()
	}

	srcKind := src.Kind()
	destKind := dest.Kind()

	if srcKind == reflect.Invalid || destKind == reflect.Invalid {
		return nil
	}

	if srcKind != destKind {
		return nil
	}

	if src.Type() != dest.Type() {
		reflectionHelper.SetFieldValue(dest, src.Interface())
		return nil
	}

	switch src.Kind() {
	case reflect.Struct:

		mapStructs[TDes, TSrc](src, dest)
	case reflect.Slice:

		mapSlices[TDes, TSrc](src, dest)
	case reflect.Map:

		mapMaps[TDes, TSrc](src, dest)
	case reflect.Ptr:
		mapPointers[TDes, TSrc](src, dest)
	default:
		dest.Set(src)
	}

	return nil
}

// mapStructs maps the fields of a source struct to a destination struct.
// It uses a mapping profile to determine which fields to map and how to map them.
// If the source and destination types have different field names, it uses the field tags to match them.
// If the source and destination types have different field tags, it uses the field names to match them.
//
// Parameters:
// - src: The reflect.Value representing the source struct to be mapped.
// - dest: The reflect.Value representing the destination struct where the mapped values will be stored.
//
// Return:
// - None. The function modifies the destination struct directly.
func mapStructs[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	profile, ok := profiles[getProfileKey(src.Type(), dest.Type())]
	if !ok {
		log.Errorf("no conversion specified for types %s and %s", src.Type().String(), dest.Type().String())
		return
	}

	for _, keys := range profile {
		destinationField := dest.FieldByName(keys[DestKeyIndex])
		sourceField := src.FieldByName(keys[SrcKeyIndex])
		var sourceFieldValue reflect.Value

		if sourceField.Kind() != reflect.Invalid {
			if !sourceField.CanInterface() {
				if mapperConfig.MapUnexportedFields {
					sourceFieldValue = reflectionHelper.GetFieldValue(sourceField)
				} else {
					sourceFieldValue = reflectionHelper.GetFieldValueFromMethodAndReflectValue(src.Addr(), strcase.ToCamel(keys[SrcKeyIndex]))
				}
			} else {
				if mapperConfig.MapUnexportedFields {
					sourceFieldValue = reflectionHelper.GetFieldValue(sourceField)
				} else {
					sourceFieldValue = sourceField
				}
			}
		} else {
			sourceFieldValue = reflectionHelper.GetFieldValueFromMethodAndReflectValue(src.Addr(), strcase.ToCamel(keys[SrcKeyIndex]))
		}

		_ = processValues[TDes, TSrc](sourceFieldValue, destinationField)
	}
}

// mapSlices maps the elements of a source slice to a destination slice.
// It creates a new slice in the destination with the same length and capacity as the source.
// Then, it iterates over each element in the source slice, processes the values using the processValues function,
// and assigns the processed values to the corresponding index in the destination slice.
//
// Parameters:
// - src: The source reflect.Value representing the slice to be mapped.
// - dest: The destination reflect.Value representing the slice where the mapped values will be stored.
//
// Return:
// - None. The function modifies the destination slice directly.
func mapSlices[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Cap()))

	for i := 0; i < src.Len(); i++ {
		srcVal := src.Index(i)
		destVal := dest.Index(i)

		_ = processValues[TDes, TSrc](srcVal, destVal)
	}
}

// mapMaps maps the elements of a source map to a destination map.
// It creates a new map in the destination with the same size as the source.
// Then, it iterates over each key-value pair in the source map, processes the keys and values using the processValues function,
// and assigns the processed keys and values to the corresponding key-value pair in the destination map.
//
// Parameters:
// - src: The source reflect.Value representing the map to be mapped.
// - dest: The destination reflect.Value representing the map where the mapped key-value pairs will be stored.
//
// Return:
// - None. The function modifies the destination map directly.
func mapMaps[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	dest.Set(reflect.MakeMapWithSize(dest.Type(), src.Len()))

	srcMapIter := src.MapRange()
	destMapIter := dest.MapRange()

	for destMapIter.Next() && srcMapIter.Next() {
		destKey := reflect.New(destMapIter.Key().Type()).Elem()
		destValue := reflect.New(destMapIter.Key().Type()).Elem()
		_ = processValues[TDes, TSrc](srcMapIter.Key(), destKey)
		_ = processValues[TDes, TSrc](srcMapIter.Value(), destValue)

		dest.SetMapIndex(destKey, destValue)
	}
}

// mapPointers maps the elements of a source pointer to a destination pointer.
// It creates a new value in the destination with the same type as the source element.
// Then, it processes the source element using the processValues function,
// and assigns the processed value to the new value in the destination.
// Finally, it sets the destination pointer to the address of the new value.
//
// Parameters:
// - src: The source reflect.Value representing the pointer to be mapped.
// - dest: The destination reflect.Value representing the pointer where the mapped value will be stored.
//
// Return:
// - None. The function modifies the destination pointer directly.
func mapPointers[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	val := reflect.New(dest.Type().Elem()).Elem()

	_ = processValues[TDes, TSrc](src.Elem(), val)

	dest.Set(val.Addr())
}

// CreateCustomMap creates a custom mapping function between source and destination types.
// It returns an error if the function is nil or if the mapping already exists.
func CreateCustomMap[TSrc any, TDes any](fn mapFunc[TSrc, TDes]) error {
	if fn == nil {
		return ErrNilFunction
	}
	var src TSrc
	var des TDes
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&des).Elem()

	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) || (desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
		return ErrUnsupportedMap
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	if _, ok := maps[k]; ok {
		return ErrMapAlreadyExists
	}

	return nil
}
