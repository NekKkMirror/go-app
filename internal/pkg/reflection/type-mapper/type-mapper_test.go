package typemapper

import (
	"reflect"
	"testing"
)

type Test struct {
	Field1 string
	Field2 int
}

func TestTypeByName(t *testing.T) {
	types["typemapper.Test"] = reflect.TypeOf(Test{})

	typ := TypeByName("typemapper.Test")
	if typ == nil {
		t.Errorf("Expected type not found")
	}

	if typ.Kind() != reflect.Struct {
		t.Errorf("Expected a struct type, got %v", typ.Kind())
	}
}

func TestGetTypeName(t *testing.T) {
	name := GetTypeName(Test{})
	expected := "typemapper.Test"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestTypeByPackageName(t *testing.T) {
	packages["typemapper"] = map[string]reflect.Type{
		"Test": reflect.TypeOf(Test{}),
	}

	typ := TypeByPackageName("typemapper", "Test")
	if typ == nil {
		t.Errorf("Expected type not found")
	}

	if typ.Kind() != reflect.Struct {
		t.Errorf("Expected a struct type, got %v", typ.Kind())
	}
}

func TestInstanceByTypeName(t *testing.T) {
	types["typemapper.Test"] = reflect.TypeOf(Test{})

	instance := InstanceByTypeName("typemapper.Test")
	if instance == nil {
		t.Errorf("Expected instance not created")
	}

	if _, ok := instance.(Test); !ok {
		t.Errorf("Expected instance of type Test, got %T", instance)
	}
}

func TestInstancePointerByTypeName(t *testing.T) {
	types["*typemapper.Test"] = reflect.TypeOf(&Test{})

	instance := InstancePointerByTypeName("*typemapper.Test")
	if instance == nil {
		t.Errorf("Expected instance pointer not created")
	}

	if _, ok := instance.(*Test); !ok {
		t.Errorf("Expected instance pointer of type *Test, got %T", instance)
	}
}

func TestInstanceByPackageName(t *testing.T) {
	packages["typemapper"] = map[string]reflect.Type{
		"Test": reflect.TypeOf(Test{}),
	}

	instance := InstanceByPackageName("typemapper", "Test")
	if instance == nil {
		t.Errorf("Expected instance not created")
	}

	if _, ok := instance.(Test); !ok {
		t.Errorf("Expected instance of type Test, got %T", instance)
	}
}

func TestGenericInstanceByTypeName(t *testing.T) {
	types["typemapper.Test"] = reflect.TypeOf(Test{})

	typ := TypeByName("typemapper.Test")
	if typ == nil {
		t.Fatalf("Type not found for typemapper.Test")
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

// Successfully creates an instance of a type when the type name is valid
func TestGenericInstanceByTypeNameValidType(t *testing.T) {
	// Assuming 'int' is a valid type name in the types map
	types["int"] = reflect.TypeOf(0)

	instance := GenericInstanceByTypeName[int]("int")

	if instance != 0 {
		t.Errorf("Expected instance of type int, got %v", instance)
	}
}

// Handles nil input gracefully without causing a panic
func TestGenericInstanceByTypeNameNilInput(t *testing.T) {
	// Assuming 'nonexistent' is not a valid type name in the types map
	instance := GenericInstanceByTypeName[int]("nonexistent")

	if instance != 0 {
		t.Errorf("Expected zero value for int, got %v", instance)
	}
}
