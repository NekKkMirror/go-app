package mapper

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestCreateMapForStructTypes(t *testing.T) {
	type Source struct {
		Field1 string
	}
	type Destination struct {
		Field1 string
	}

	err := CreateMap[Source, Destination]()

	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	srcType := reflect.TypeOf(Source{})
	desType := reflect.TypeOf(Destination{})

	pointerStructTypeKey := mappingsEntry{
		SourceType:      srcType,
		DestinationType: desType,
	}

	if _, ok := maps[pointerStructTypeKey]; !ok {
		t.Errorf("expected map to be created for struct types")
	}
}

func TestCreateMapForNonStructTypes(t *testing.T) {
	type Source int
	type Destination int

	err := CreateMap[Source, Destination]()

	if !errors.Is(err, ErrUnsupportedMap) {
		t.Errorf("expected error %v, but got %v", ErrUnsupportedMap, err)
	}
}

func TestCreateMapWithValidStructTypes(t *testing.T) {
	type Source struct{}
	type Destination struct{}

	err := CreateMap[Source, Destination]()

	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	srcType := reflect.TypeOf(Source{})
	desType := reflect.TypeOf(Destination{})

	pointerStructTypeKey := mappingsEntry{
		SourceType:      reflect.PointerTo(srcType),
		DestinationType: reflect.PointerTo(desType),
	}
	nonePointerStructTypeKey := mappingsEntry{
		SourceType:      srcType,
		DestinationType: desType,
	}

	if _, ok := maps[pointerStructTypeKey]; !ok {
		t.Errorf("expected map entry for pointer struct types")
	}
	if _, ok := maps[nonePointerStructTypeKey]; !ok {
		t.Errorf("expected map entry for non-pointer struct types")
	}
}

func TestCreateMapWithNonStructTypes(t *testing.T) {
	type Source int
	type Destination int

	err := CreateMap[Source, Destination]()

	if !errors.Is(err, ErrUnsupportedMap) {
		t.Errorf("expected ErrUnsupportedMap, but got %v", err)
	}
}

func TestMappingBetweenCompatibleTypesSucceeds(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestReturnsErrMapNotExistWhenNoMappingFunctionFound(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}

	_, err := Map[Source, Destination](src)

	if !errors.Is(err, ErrMapNotExist) {
		t.Fatalf("expected error %v, got %v", ErrMapNotExist, err)
	}
}

func TestMappingWithExistingMapping(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestMappingWithCompatibleTypes(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestMappingProcessValuesCorrectly(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestMappingWithComplexNestedStructures(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestMapUsingLinqForArrayAndSliceTransformations(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestMaintainsImmutabilityOfSourceDataDuringMapping(t *testing.T) {
	type Source struct {
		Name string
	}
	type Destination struct {
		Name string
	}

	src := Source{Name: "Test"}
	expected := Destination{Name: "Test"}

	maps[mappingsEntry{SourceType: reflect.TypeOf(src), DestinationType: reflect.TypeOf(expected)}] = func(s Source) Destination {
		return Destination{Name: s.Name}
	}

	result, err := Map[Source, Destination](src)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestCreateCustomMapWithValidFunction(t *testing.T) {
	type Source struct{}
	type Destination struct{}

	fn := func(src Source) Destination {
		return Destination{}
	}

	err := CreateCustomMap[Source, Destination](fn)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateCustomMapWithNilFunction(t *testing.T) {
	err := CreateCustomMap[struct{}, struct{}](nil)

	if !errors.Is(err, ErrNilFunction) {
		t.Errorf("expected ErrNilFunction, got %v", err)
	}
}
