package reflectionhelper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PersonPublic struct {
	Name string
	Age  int
}

type PersonPrivate struct {
	name string
	age  int
}

func (p *PersonPrivate) Name() string {
	return p.name
}

func (p *PersonPrivate) Age() int {
	return p.age
}

func Test_Field_Values_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{Name: "John", Age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Field_Values_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{Name: "John", Age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Field_Values_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Field_Values_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Set_Field_Value_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{}

	SetFieldValueByIndex(p, 0, "John")
	SetFieldValueByIndex(p, 1, 20)

	assert.Equal(t, "John", p.Name)
	assert.Equal(t, 20, p.Age)
}

func Test_Set_Field_Value_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{}

	SetFieldValueByIndex(&p, 0, "John")
	SetFieldValueByIndex(&p, 1, 20)

	assert.Equal(t, "John", p.Name)
	assert.Equal(t, 20, p.Age)
}

func Test_Set_Field_Value_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{}

	SetFieldValueByIndex(p, 0, "John")
	SetFieldValueByIndex(p, 1, 20)

	assert.Equal(t, "John", p.name)
	assert.Equal(t, 20, p.age)
}

func Test_Set_Field_Value_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{}

	SetFieldValueByIndex(&p, 0, "John")
	SetFieldValueByIndex(&p, 1, 20)

	assert.Equal(t, "John", p.name)
	assert.Equal(t, 20, p.age)
}

func Test_Get_Field_Value_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{Name: "John", Age: 20}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("Name")).Interface()
	age := GetFieldValue(v.FieldByName("Age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func Test_Get_Field_Value_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 30}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("name")).Interface()
	age := GetFieldValue(v.FieldByName("age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 30, age)
}

func Test_Get_Field_Value_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{Name: "John", Age: 20}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("Name")).Interface()
	age := GetFieldValue(v.FieldByName("Age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func Test_Get_Field_Value_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 20}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("name")).Interface()
	age := GetFieldValue(v.FieldByName("age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func Test_Set_Field_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("Name"))
	age := GetFieldValue(v.FieldByName("Age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Set_Field_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("name"))
	age := GetFieldValue(v.FieldByName("age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Set_Field_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("Name"))
	age := GetFieldValue(v.FieldByName("Age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Set_Field_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("name"))
	age := GetFieldValue(v.FieldByName("age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Convert_NoPointer_Type_To_Pointer_Type_With_Addr(t *testing.T) {

	p := PersonPrivate{name: "John", age: 20}
	v := reflect.ValueOf(&p).Elem()
	pointerType := v.Addr()
	name := pointerType.MethodByName("Name").Call(nil)[0].Interface()
	age := pointerType.MethodByName("Age").Call(nil)[0].Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func TestGetFieldValueByNameWithValidFieldName(t *testing.T) {
	obj := PersonPublic{Name: "John", Age: 30}
	fieldName := "Name"
	expectedValue := "John"

	result := GetFieldValueByName(obj, fieldName)

	assert.Equal(t, expectedValue, result)
}

func TestSetFieldValueByName_Success(t *testing.T) {
	obj := PersonPublic{}
	SetFieldValueByName(&obj, "Name", "new value")

	if obj.Name != "new value" {
		t.Errorf("Expected Field1 to be 'new value', got %s", obj.Name)
	}
}

func TestSetFieldValueByName_NonExistentField(t *testing.T) {
	obj := PersonPublic{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when setting non-existent field, but did not panic")
		}
	}()

	SetFieldValueByName(&obj, "NonExistentField", "value")
}

func TestSetValueAssignsNonPointerValue(t *testing.T) {
	var data int
	value := 42

	SetValue(&data, value)

	if data != value {
		t.Errorf("Expected data to be %d, but got %d", value, data)
	}
}

func TestSetValueHandlesNilPointer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetValue panicked with nil pointer: %v", r)
		}
	}()

	var data *int
	SetValue(&data, nil)

	if data != nil {
		t.Errorf("Expected data to be nil, but got %v", data)
	}
}

// Returns correct type path for a basic type
func TestTypePathBasicType(t *testing.T) {
	type MyStruct struct{}
	expected := "github.com/NekKkMirror/go-app/internal/pkg/reflection/reflection-helper.MyStruct"
	result := TypePath[MyStruct]()
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

// Handles nil pointer gracefully
func TestTypePathNilPointer(t *testing.T) {
	type MyStruct struct{}
	var _ *MyStruct
	expected := "github.com/NekKkMirror/go-app/internal/pkg/reflection/reflection-helper.MyStruct"
	result := TypePath[*MyStruct]()
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestCreateInstanceReturnsNewInstance(t *testing.T) {
	instance := CreateInstance[PersonPublic]()
	if reflect.DeepEqual(instance, &PersonPublic{}) {
		t.Errorf("Expected a new instance of PersonPublic, but got zero value")
	}

	// Test for pointer instance
	ptrInstance := CreateInstance[*PersonPublic]()
	if ptrInstance == nil {
		t.Errorf("Expected a new instance of *PersonPublic, but got nil")
	}

	// Check that the pointer indeed points to a properly initialized value
	if ptrInstance != nil && (ptrInstance.Name != "" || ptrInstance.Age != 0) {
		t.Errorf("Expected initialized fields in *PersonPublic, got %v", ptrInstance)
	}
}

// Correctly extracts method path from function pointer
func TestExtractMethodPathFromFunctionPointer(t *testing.T) {
	testFunc := func() {}
	expected := "reflection-helper:TestExtractMethodPathFromFunctionPointer:func1"
	result := MethodPath(testFunc)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

// Handles functions with no slashes in their names
func TestHandleFunctionsWithNoSlashes(t *testing.T) {
	testFunc := func() {}
	expected := "reflection-helper:TestHandleFunctionsWithNoSlashes:func1"
	result := MethodPath(testFunc)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
func (p *PersonPublic) GetName() string {
	return p.Name
}

func TestGetFieldValueFromMethodAndObjectWithStruct(t *testing.T) {
	obj := &PersonPublic{Name: "home"}
	result := GetFieldValueFromMethodAndObject[PersonPublic](*obj, "GetName")
	if result.Kind() != reflect.String || result.String() != "home" {
		t.Errorf("Expected name 'home', got %v", result)
	}
}
