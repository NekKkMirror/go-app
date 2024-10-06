package type_mapper

import (
	"reflect"
	"testing"
)

type Test struct {
	Field1 string
	Field2 int
}

func TestTypeByName(t *testing.T) {
	types["type_mapper.Test"] = reflect.TypeOf(Test{})

	typ := TypeByName("type_mapper.Test")
	if typ == nil {
		t.Errorf("Expected type not found")
	}

	if typ.Kind() != reflect.Struct {
		t.Errorf("Expected a struct type, got %v", typ.Kind())
	}
}

func TestGetTypeName(t *testing.T) {
	name := GetTypeName(Test{})
	expected := "type_mapper.Test"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestTypeByPackageName(t *testing.T) {
	packages["type_mapper"] = map[string]reflect.Type{
		"Test": reflect.TypeOf(Test{}),
	}

	typ := TypeByPackageName("type_mapper", "Test")
	if typ == nil {
		t.Errorf("Expected type not found")
	}

	if typ.Kind() != reflect.Struct {
		t.Errorf("Expected a struct type, got %v", typ.Kind())
	}
}

func TestInstanceByTypeName(t *testing.T) {
	types["type_mapper.Test"] = reflect.TypeOf(Test{})

	instance := InstanceByTypeName("type_mapper.Test")
	if instance == nil {
		t.Errorf("Expected instance not created")
	}

	if _, ok := instance.(Test); !ok {
		t.Errorf("Expected instance of type Test, got %T", instance)
	}
}

func TestInstancePointerByTypeName(t *testing.T) {
	types["*type_mapper.Test"] = reflect.TypeOf(&Test{})

	instance := InstancePointerByTypeName("*type_mapper.Test")
	if instance == nil {
		t.Errorf("Expected instance pointer not created")
	}

	if _, ok := instance.(*Test); !ok {
		t.Errorf("Expected instance pointer of type *Test, got %T", instance)
	}
}

func TestInstanceByPackageName(t *testing.T) {
	packages["type_mapper"] = map[string]reflect.Type{
		"Test": reflect.TypeOf(Test{}),
	}

	instance := InstanceByPackageName("type_mapper", "Test")
	if instance == nil {
		t.Errorf("Expected instance not created")
	}

	if _, ok := instance.(Test); !ok {
		t.Errorf("Expected instance of type Test, got %T", instance)
	}
}

func TestGenericInstanceByTypeName(t *testing.T) {
	types["type_mapper.Test"] = reflect.TypeOf(Test{})

	typ := TypeByName("type_mapper.Test")
	if typ == nil {
		t.Fatalf("Type not found for type_mapper.Test")
	}

	instance := getInstanceFromType(typ)

	if instance == nil {
		t.Errorf("Expected non-nil instance")
	}

	testInstance, ok := instance.(Test)
	if !ok {
		t.Errorf("Expected instance of type Test, got %T", instance)
	}

	if testInstance != (Test{}) {
		t.Errorf("Expected zero instance of type Test, got %+v", testInstance)
	}
}
