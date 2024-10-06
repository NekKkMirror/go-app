package typeregistry

import (
	"reflect"
)

// Global type registry mapping type names to their reflect.Type.
var typeRegistry = make(map[string]reflect.Type)

type MyString string
type myString string

func init() {
	initializeTypeRegistry()
}

// initializeTypeRegistry initializes the type registry by registering the types MyString and myString.
func initializeTypeRegistry() {
	registerType((*MyString)(nil))
	registerType((*myString)(nil))
}

// registerType registers a given type into the typeRegistry map.
func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeName := t.PkgPath() + "." + t.Name()
	typeRegistry[typeName] = t
}

// makeInstance creates a new instance of a type registered in the typeRegistry map.
func makeInstance(name string) (interface{}, bool) {
	typ, exists := typeRegistry[name]
	if !exists {
		return nil, false
	}
	return reflect.New(typ).Elem().Interface(), true
}
