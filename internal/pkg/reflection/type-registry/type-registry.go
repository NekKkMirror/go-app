package typeregistry

import (
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

type MyString string
type myString string

// init initializes the type registry by registering the types MyString and myString.
// It uses the registerType function to store the reflect.Type information of these types in the typeRegistry map.
func init() {
	registerType((*MyString)(nil))
	registerType((*myString)(nil))
}

// registerType is a function that registers a given type into the typeRegistry map.
// The function accepts an interface{} parameter, typedNil, which is expected to be a pointer to a type.
// The function extracts the type information from the interface{} using reflection and stores it in the typeRegistry map.
// The key in the map is a string composed of the package path and the type name, and the value is the reflect.Type of the type.
func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = t
}

// makeInstance creates a new instance of a type registered in the typeRegistry map.
// The function accepts a string parameter, name, which represents the key in the typeRegistry map.
// The function retrieves the reflect.Type associated with the given name from the typeRegistry map.
// It then creates a new instance of the type using reflection and returns it as an interface{}.
//
// If the given name does not exist in the typeRegistry map, the function will panic.
//
// Example usage:
//
//	type MyString string
//	registerType((*MyString)(nil))
//
//	myInstance := makeInstance("github.com/mypackage.MyString").(MyString)
func makeInstance(name string) interface{} {
	return reflect.New(typeRegistry[name]).Elem().Interface()
}
