package typemapper

import (
	"reflect"
	"unsafe"
)

// GetTypeName returns the fully qualified name of the type of the input value.
func GetTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	return t.String()
}

// InstancePointerByTypeName retrieves an instance of a type by its name and returns a pointer to the instance.
func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Interface()
}

// InstanceByPackageName creates an instance by package path and type name.
func InstanceByPackageName(pkgPath string, name string) interface{} {
	typ := TypeByPackageName(pkgPath, name)
	return getInstanceFromType(typ)
}

// getInstanceFromType creates an instance from a given reflect.Type.
func getInstanceFromType(typ reflect.Type) interface{} {
	if typ == nil {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Elem().Interface()
}

// GenericInstanceByTypeName creates an instance for a type name using generics.
func GenericInstanceByTypeName[T any](typeName string) T {
	instance := InstanceByTypeName(typeName)
	if instance == nil {
		var zero T
		return zero
	}
	return instance.(T)
}

// TypeByPackageName retrieves a reflect.Type by package path and type name.
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

// InstanceByTypeName creates an instance by type name.
func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)
	return getInstanceFromType(typ)
}

// TypeByName retrieves a type by its name.
func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ
	}
	return nil
}

var (
	types    map[string]reflect.Type
	packages map[string]map[string]reflect.Type
)

// init initializes the types and packages maps, and calls discoverTypes function.
func init() {
	types = make(map[string]reflect.Type)
	packages = make(map[string]map[string]reflect.Type)

	discoverTypes()
}

// discoverTypes is a helper function that discovers and registers all types in the current program.
func discoverTypes() {
	// Dummy implementation of typelinks2 and resolveTypeOff
	typelinks2 := func() ([][]byte, [][]int32) { return nil, nil }
	resolveTypeOff := func(_ []byte, _ int32) unsafe.Pointer { return nil }

	zeroValueType := reflect.TypeOf(0)
	sections, offsets := typelinks2()

	for i, offs := range offsets {
		rodata := sections[i]
		for _, off := range offs {
			interfaceType := (*emptyInterface)(unsafe.Pointer(&zeroValueType))
			interfaceType.data = resolveTypeOff(rodata, off)

			if zeroValueType.Kind() == reflect.Ptr && zeroValueType.Elem().Kind() == reflect.Struct {
				ptrType := zeroValueType
				elemType := zeroValueType.Elem()

				registerType(ptrType, elemType)
			}
		}
	}
}

// registerType registers pointer and element types into the maps.
func registerType(ptrType, elemType reflect.Type) {
	pkgPath := elemType.PkgPath()
	ptrPkgPath := ptrType.PkgPath()

	if packages[pkgPath] == nil {
		packages[pkgPath] = make(map[string]reflect.Type)
	}
	if packages[ptrPkgPath] == nil {
		packages[ptrPkgPath] = make(map[string]reflect.Type)
	}

	types[elemType.String()] = elemType
	types[ptrType.String()] = ptrType
	packages[pkgPath][elemType.Name()] = elemType
	packages[ptrPkgPath][ptrType.Name()] = ptrType
}

// emptyInterface is a placeholder for the internal representation of an empty interface.
type emptyInterface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}
