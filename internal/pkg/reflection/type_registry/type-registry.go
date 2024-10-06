package type_registry

import (
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = t
}

type MyString string
type myString string

func init() {
	registerType((*MyString)(nil))
	registerType((*myString)(nil))
}

func makeInstance(name string) interface{} {
	return reflect.New(typeRegistry[name]).Elem().Interface()
}
