package typeregistry

import (
	"reflect"
	"testing"
)

const (
	pubKey string = "github.com/NekKkMirror/go-app/internal/pkg/reflection/type-registry.MyString"
	priKey string = "github.com/NekKkMirror/go-app/internal/pkg/reflection/type-registry.myString"
)

func TestRegisterType(t *testing.T) {
	typeRegistry = make(map[string]reflect.Type)

	registerType((*MyString)(nil))
	if _, exists := typeRegistry[pubKey]; !exists {
		t.Errorf("Expected typeregistry.MyString to be registered")
	}

	registerType((*myString)(nil))

	if _, exists := typeRegistry[priKey]; !exists {
		t.Errorf("Expected typeregistry.myString to be registered")
	}
}

func TestMakeInstance(t *testing.T) {
	typeRegistry = make(map[string]reflect.Type)

	registerType((*MyString)(nil))

	instance, _ := makeInstance(pubKey)

	if _, ok := instance.(MyString); !ok {
		t.Errorf("Expected instance of type MyString, got %T", instance)
	}

	registerType((*myString)(nil))

	instance, _ = makeInstance(priKey)

	if _, ok := instance.(myString); !ok {
		t.Errorf("Expected instance of type myString, got %T", instance)
	}
}
