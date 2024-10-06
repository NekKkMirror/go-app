package mapper

import (
	"flag"
	"fmt"
	"reflect"

	reflectionHelper "github.com/NekKkMirror/go-app/internal/pkg/reflection/reflection_helper"
	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNilFunction = errors.New("mapper: nil function")

	ErrMapNotExist = errors.New("mapper: map does not exist")

	ErrMapAlreadyExists = errors.New("mapper: map already exists")

	ErrUnsupportedMap = errors.New("mapper: unsupported map")
)

const (
	SrcKeyIndex = iota
	DestKeyIndex
)

type Config struct {
	MapUnexportedFields bool
}

type mappingsEntry struct {
	SourceType      reflect.Type
	DestinationType reflect.Type
}

type typeMeta struct {
	keysToTags map[string]string
	tagsToKeys map[string]string
}

type mapFunc[TSrc any, TDst any] func(TSrc) TDst

var profiles = map[string][][2]string{}
var maps = map[mappingsEntry]interface{}{}
var mapperConfig *Config

func init() {
	mapperConfig = &Config{
		MapUnexportedFields: false,
	}
}

func Configure(config *Config) {
	mapperConfig = config
}

func CreateMap[TSrc any, TDst any]() error {
	var src TSrc
	var dst TDst
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) ||
		(desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
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
			})
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

func mapSlices[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Cap()))

	for i := 0; i < src.Len(); i++ {
		srcVal := src.Index(i)
		destVal := dest.Index(i)

		_ = processValues[TDes, TSrc](srcVal, destVal)
	}
}

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

func mapPointers[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	val := reflect.New(dest.Type().Elem()).Elem()

	_ = processValues[TDes, TSrc](src.Elem(), val)

	dest.Set(val.Addr())
}

func CreateCustomMap[TSrc any, TDes any](fn mapFunc[TSrc, TDes]) error {
	if fn == nil {
		return ErrNilFunction
	}
	var src TSrc
	var des TDes
	srcType := reflect.TypeOf(src).Elem()
	desType := reflect.TypeOf(des).Elem()

	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) || (desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
		return ErrUnsupportedMap
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	if _, ok := maps[k]; ok {
		return ErrMapAlreadyExists
	}

	return nil
}
