package mapper

import (
	"fmt"
	"reflect"

	reflectionHelper "github.com/NekKkMirror/go-app/internal/pkg/reflection/reflection-helper"
	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

var (
	ErrNilFunction = errors.New("mapper: nil function")

	ErrMapNotExist = errors.New("mapper: map does not exist")

	ErrMapAlreadyExists = errors.New("mapper: map already exists")

	ErrUnsupportedMap = errors.New("mapper: unsupported map")

	ErrInvalidStructType = errors.New("mapper: expected reflect.Struct kind for type")
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

// CreateMap registers a mapping configuration between two types.
func CreateMap[TSrc any, TDst any]() error {
	var src TSrc
	var dst TDst

	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	// Check if types are valid for mapping
	if !isStructOrPointerToStruct(srcType) || !isStructOrPointerToStruct(desType) {
		return ErrUnsupportedMap
	}

	// Define mappings for both pointer and non-pointer struct types
	pointerStructTypeKey := mappingsEntry{SourceType: reflect.PointerTo(srcType), DestinationType: reflect.PointerTo(desType)}
	nonePointerStructTypeKey := mappingsEntry{SourceType: srcType, DestinationType: desType}

	// Check for existing mappings
	if _, exists := maps[pointerStructTypeKey]; exists {
		return ErrMapAlreadyExists
	}
	if _, exists := maps[nonePointerStructTypeKey]; exists {
		return ErrMapAlreadyExists
	}

	// Register new mappings
	maps[pointerStructTypeKey] = nil
	maps[nonePointerStructTypeKey] = nil

	// Configure profile between the base types
	err := configProfile(getBaseType(srcType), getBaseType(desType))
	if err != nil {
		return err
	}

	return nil
}

// isStructOrPointerToStruct checks if the given type is a struct or a pointer to a struct.
func isStructOrPointerToStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

// getBaseType returns the underlying type if it's a pointer, otherwise returns the type itself.
func getBaseType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return t.Elem()
	}
	return t
}

// configProfile configures a profile for mapping source and destination types.
// It returns an error if the types are not structs.
func configProfile(srcType reflect.Type, desType reflect.Type) error {
	if srcType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %s, but got %s", ErrInvalidStructType, srcType.String(), srcType.Kind().String())
	}

	if desType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %s, but got %s", ErrInvalidStructType, desType.String(), desType.Kind().String())
	}

	profile := createProfile(srcType, desType)
	profiles[getProfileKey(srcType, desType)] = profile
	return nil
}

// createProfile creates a mapping profile between the source and destination types.
func createProfile(srcType, desType reflect.Type) [][2]string {
	var profile [][2]string
	srcMeta := getTypeMeta(srcType)
	destMeta := getTypeMeta(desType)
	srcMethods := getTypeMethods(srcType)

	for srcKey, srcTag := range srcMeta.keysToTags {
		if destKey, ok := getDestinationKey(srcKey, srcTag, destMeta); ok {
			profile = append(profile, [2]string{srcKey, destKey})
		}
	}

	for _, method := range srcMethods {
		if _, ok := destMeta.keysToTags[method]; ok {
			profile = append(profile, [2]string{method, method})
		}
	}

	return profile
}

// getTypeMeta returns the key and tag mappings for a struct type.
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

// getTypeMethods returns the method names of a struct type.
func getTypeMethods(val reflect.Type) []string {
	methodsNum := val.NumMethod()
	var keys []string

	for i := 0; i < methodsNum; i++ {
		methodName := val.Method(i).Name
		keys = append(keys, methodName)
	}

	return keys
}

// getDestinationKey finds the corresponding destination key for a given source key or tag.
func getDestinationKey(srcKey, srcTag string, destMeta typeMeta) (string, bool) {
	if _, ok := destMeta.keysToTags[strcase.ToCamel(srcKey)]; ok {
		return strcase.ToCamel(srcKey), true
	}

	if _, ok := destMeta.keysToTags[srcKey]; ok {
		return srcKey, true
	}

	if destKey, ok := destMeta.tagsToKeys[srcKey]; ok {
		return destKey, true
	}

	if _, ok := destMeta.keysToTags[srcTag]; ok {
		return srcTag, true
	}

	if destKey, ok := destMeta.tagsToKeys[srcTag]; ok {
		return destKey, true
	}

	return "", false
}

// getProfileKey returns a unique key for the source and destination types.
func getProfileKey(srcType, desType reflect.Type) string {
	return fmt.Sprintf("%s_%s", srcType.Name(), desType.Name())
}

// Map is a generic function that maps a source value to a destination value of different types.
func Map[TSrc any, TDes any](src TSrc) (TDes, error) {
	var des TDes
	srcType, srcIsArray := getElementType(reflect.TypeOf(src))
	desType, desIsArray := getElementType(reflect.TypeOf(des))

	fn, err := getMappingFunction(srcType, desType)
	if err != nil {
		return des, err
	}

	if fn != nil {
		fnReflect := reflect.ValueOf(fn)
		if desIsArray && srcIsArray {
			mapArray(fnReflect, src, &des)
		} else {
			mappedValue := fnReflect.Call([]reflect.Value{reflect.ValueOf(src)})[0].Interface()
			return mappedValue.(TDes), nil
		}
	}

	err = processValues[TSrc, TDes](reflect.ValueOf(src), reflect.ValueOf(&des).Elem())
	if err != nil {
		return des, err
	}

	return des, nil
}

// getElementType determines if the given type is an array, pointer to an array, or slice, and returns the element type.
func getElementType(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		return t.Elem(), true
	}
	if t.Kind() == reflect.Ptr && (t.Elem().Kind() == reflect.Array || t.Elem().Kind() == reflect.Slice) {
		return t.Elem().Elem(), true
	}
	return t, false
}

// getMappingFunction retrieves the mapping function for the given source and destination types.
func getMappingFunction(srcType, desType reflect.Type) (interface{}, error) {
	key := mappingsEntry{SourceType: srcType, DestinationType: desType}
	fn, ok := maps[key]
	if !ok {
		return nil, ErrMapNotExist
	}
	return fn, nil
}

// mapArray applies a mapping function to each element of a source array and stores the results in the destination array.
func mapArray[TSrc any, TDes any](fn reflect.Value, src TSrc, des *TDes) {
	linq.From(src).Select(func(x interface{}) interface{} {
		return fn.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
	}).ToSlice(des)
}

// processValues is a generic function that handles the mapping of values from a source to a destination.
func processValues[TSrc any, TDes any](src reflect.Value, dest reflect.Value) error {
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	switch src.Kind() {
	case reflect.Struct:
		mapStructs[TSrc, TDes](src, dest)
	case reflect.Slice:
		mapSlices[TSrc, TDes](src, dest)
	case reflect.Map:
		mapMaps[TSrc, TDes](src, dest)
	case reflect.Ptr:
		mapPointers[TSrc, TDes](src, dest)
	default:
		dest.Set(src)
	}

	return nil
}

func mapStructs[TSrc any, TDes any](src reflect.Value, dest reflect.Value) {
	profileKey := getProfileKey(src.Type(), dest.Type())
	profile, exists := profiles[profileKey]
	if !exists {
		return
	}

	for _, keys := range profile {
		sourceField := retrieveSourceFieldValue(src, keys[SrcKeyIndex])
		destinationField := dest.FieldByName(keys[DestKeyIndex])
		_ = processValues[TSrc, TDes](sourceField, destinationField)
	}
}

// retrieveSourceFieldValue retrieves the value of a field from a source reflect.Value.
func retrieveSourceFieldValue(src reflect.Value, fieldName string) reflect.Value {
	field := src.FieldByName(fieldName)
	if field.Kind() != reflect.Invalid {
		if field.CanInterface() || !mapperConfig.MapUnexportedFields {
			return field
		}
		return reflectionHelper.GetFieldValue(field)
	}
	return reflectionHelper.GetFieldValueFromMethodAndReflectValue(src.Addr(), strcase.ToCamel(fieldName))
}

func mapSlices[TSrc any, TDes any](src reflect.Value, dest reflect.Value) {
	dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Cap()))

	for i := 0; i < src.Len(); i++ {
		_ = processValues[TSrc, TDes](src.Index(i), dest.Index(i))
	}
}

func mapMaps[TSrc any, TDes any](src reflect.Value, dest reflect.Value) {
	dest.Set(reflect.MakeMapWithSize(dest.Type(), src.Len()))
	srcMapIter := src.MapRange()

	for srcMapIter.Next() {
		destKey := reflect.New(dest.Type().Key()).Elem()
		destValue := reflect.New(dest.Type().Elem()).Elem()
		_ = processValues[TSrc, TDes](srcMapIter.Key(), destKey)
		_ = processValues[TSrc, TDes](srcMapIter.Value(), destValue)
		dest.SetMapIndex(destKey, destValue)
	}
}

func mapPointers[TSrc any, TDes any](src reflect.Value, dest reflect.Value) {
	if src.IsNil() {
		dest.Set(reflect.Zero(dest.Type()))
		return
	}

	dest.Set(reflect.New(dest.Type().Elem()))
	_ = processValues[TSrc, TDes](src.Elem(), dest.Elem())
}

// CreateCustomMap registers a custom mapping function between two types.
func CreateCustomMap[TSrc any, TDes any](fn mapFunc[TSrc, TDes]) error {
	if fn == nil {
		return ErrNilFunction
	}

	var src TSrc
	var des TDes
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&des).Elem()

	// Check if the source type and destination type are structs or pointers to structs.
	if !isStructOrPointerToStruct(srcType) || !isStructOrPointerToStruct(desType) {
		return ErrUnsupportedMap
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	if _, exists := maps[k]; exists {
		return ErrMapAlreadyExists
	}

	maps[k] = struct{}{}
	return nil
}
