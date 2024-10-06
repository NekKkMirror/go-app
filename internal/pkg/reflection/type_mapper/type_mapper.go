package type_mapper

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var types map[string]reflect.Type
var packages map[string]map[string]reflect.Type

func init() {
	types = make(map[string]reflect.Type)
	packages = make(map[string]map[string]reflect.Type)

	discoverTypes()
}

func discoverTypes() {
	typ := reflect.TypeOf(0)
	sections, offset := typelinks2()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			emptyInterface := (*emptyInterface)(unsafe.Pointer(&typ))
			emptyInterface.data = resolveTypeOff(rodata, off)
			if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {

				loadedTypePtr := typ
				loadedType := typ.Elem()

				pkgTypes := packages[loadedType.PkgPath()]
				pkgTypesPtr := packages[loadedTypePtr.PkgPath()]

				if pkgTypes == nil {
					pkgTypes = map[string]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				if pkgTypesPtr == nil {
					pkgTypesPtr = map[string]reflect.Type{}
					packages[loadedTypePtr.PkgPath()] = pkgTypesPtr
				}
				f := strings.Contains(loadedType.String(), "Test")
				if f {
					fmt.Println(loadedType.String())
				}

				types[loadedType.String()] = loadedType
				types[loadedTypePtr.String()] = loadedTypePtr
				pkgTypes[loadedType.Name()] = loadedType
				pkgTypesPtr[loadedTypePtr.Name()] = loadedTypePtr
			}
		}
	}
}

func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ
	}
	return nil
}

func GetTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	return t.String()
}

func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)
	return getInstanceFromType(typ)
}

func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Interface()
}

func InstanceByPackageName(pkgPath string, name string) interface{} {
	typ := TypeByPackageName(pkgPath, name)
	return getInstanceFromType(typ)
}

func getInstanceFromType(typ reflect.Type) interface{} {
	if typ == nil {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Elem().Interface()
}

func GenericInstanceByTypeName[T any](typeName string) T {
	instance := InstanceByTypeName(typeName)
	if instance == nil {
		var zero T
		return zero
	}
	return instance.(T)
}
