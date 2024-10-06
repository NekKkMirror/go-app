package typemapper

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var types map[string]reflect.Type
var packages map[string]map[string]reflect.Type

// init initializes the type and package maps and calls discoverTypes function.
func init() {
	types = make(map[string]reflect.Type)
	packages = make(map[string]map[string]reflect.Type)

	discoverTypes()
}

// discoverTypes is a helper function that populates the types and packages maps with all the
// available types in the current Go program. It uses reflection to iterate through the type
// links and resolve the types. It also filters out test types and prints their names.
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

// TypeByName retrieves a reflect.Type from the internal type map by its name.
//
// The function takes a string parameter `typeName` which represents the name of the type to be retrieved.
// It searches for the type in the `types` map and returns the corresponding reflect.Type if found.
// If the type is not found in the map, the function returns nil.
//
// Return value:
// - reflect.Type: The reflect.Type corresponding to the given typeName if found, otherwise nil.
func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ
	}
	return nil
}

// GetTypeName returns the fully qualified name of the type of the input value.
//
// The function uses reflection to determine the type of the input value and returns its fully qualified name.
// The fully qualified name includes the package path and the type name.
//
// Parameters:
// - input: The input value for which the type name needs to be determined. It can be of any type.
//
// Return:
// - A string representing the fully qualified name of the type of the input value.
func GetTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	return t.String()
}

// TypeByPackageName retrieves a reflect.Type from the internal package map by its package path and name.
//
// The function takes two string parameters:
// - pkgPath: The package path of the type to be retrieved.
// - name: The name of the type to be retrieved.
//
// It searches for the type in the `packages` map using the provided package path and name.
// If the type is found, it returns the corresponding reflect.Type.
// If the type is not found in the map, the function returns nil.
//
// Return value:
// - reflect.Type: The reflect.Type corresponding to the given package path and name if found, otherwise nil.
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

// InstanceByTypeName retrieves an instance of a type by its name.
//
// The function takes a string parameter `name` which represents the name of the type to be retrieved.
// It first calls the TypeByName function to get the corresponding reflect.Type.
// Then, it calls the getInstanceFromType function to create an instance of the type.
//
// Parameters:
// - name: A string representing the name of the type to be retrieved.
//
// Return:
// - An interface{} containing the instance of the type. If the type is not found, the function returns nil.
func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)
	return getInstanceFromType(typ)
}

// InstancePointerByTypeName retrieves an instance of a type by its name and returns a pointer to the instance.
// If the type is not a pointer, it wraps the instance in a new pointer before returning.
//
// Parameters:
// - name: A string representing the name of the type to be retrieved.
//
// Return:
//   - An interface{} containing a pointer to the instance of the type.
//     If the type is not found, the function returns nil.
func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Interface()
}

// InstanceByPackageName retrieves an instance of a type by its package path and name.
// If the type is not found, it returns nil.
// If the type is a pointer, it returns a pointer to the instance.
// If the type is not a pointer, it wraps the instance in a new pointer before returning.
//
// Parameters:
// - pkgPath: A string representing the package path of the type to be retrieved.
// - name: A string representing the name of the type to be retrieved.
//
// Return:
//   - An interface{} containing a pointer to the instance of the type.
//     If the type is not found, the function returns nil.
func InstanceByPackageName(pkgPath string, name string) interface{} {
	typ := TypeByPackageName(pkgPath, name)
	return getInstanceFromType(typ)
}

// getInstanceFromType creates a new instance of a given reflect.Type.
// If the type is a pointer, it returns a pointer to the instance.
// If the type is not a pointer, it wraps the instance in a new pointer before returning.
//
// Parameters:
// - typ: The reflect.Type for which to create an instance.
//
// Return:
//   - An interface{} containing a pointer to the instance of the type.
//     If the input type is nil, the function returns nil.
func getInstanceFromType(typ reflect.Type) interface{} {
	if typ == nil {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Elem().Interface()
}

// GenericInstanceByTypeName retrieves an instance of a type by its name and returns it as the specified generic type.
// If the type is not found, it returns a zero value of the specified generic type.
//
// The function uses reflection to dynamically create an instance of the specified type by its name.
// It first calls the InstanceByTypeName function to get the instance as an interface{}.
// If the instance is nil, it returns a zero value of the specified generic type.
// Otherwise, it performs a type assertion to convert the instance to the specified generic type and returns it.
//
// Parameters:
// - typeName: A string representing the name of the type to be retrieved.
//
// Return:
//   - T: The instance of the specified type as the generic type T.
//     If the type is not found, it returns a zero value of the specified generic type.
func GenericInstanceByTypeName[T any](typeName string) T {
	instance := InstanceByTypeName(typeName)
	if instance == nil {
		var zero T
		return zero
	}
	return instance.(T)
}
