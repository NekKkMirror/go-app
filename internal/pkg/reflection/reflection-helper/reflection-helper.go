package reflectionhelper

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// GetFieldValueByIndex retrieves the value of a field at the specified index in the given object.
// The function supports both pointer and non-pointer types.
// If the object is a pointer, it dereferences it before accessing the field.
// If the object is a struct, it directly accesses the field.
// If the field is not accessible or does not exist, the function returns nil.
//
// Parameters:
// - object: The object containing the field.
// - index: The index of the field to retrieve.
//
// Returns:
// - The value of the field at the specified index, or nil if the field is not accessible or does not exist.
func GetFieldValueByIndex[T any](object T, index int) interface{} {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.Field(index)
		if field.CanInterface() {
			return field.Interface()
		}
		return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
	} else if v.Kind() == reflect.Struct {
		val := v
		field := val.Field(index)
		if field.CanInterface() {
			return field.Interface()
		}
		rs2 := reflect.New(val.Type()).Elem()
		rs2.Set(val)
		val = rs2.Field(index)
		return reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem().Interface()
	}

	return nil
}

// GetFieldValueBuName retrieves the value of a field with the specified name in the given object.
// The function supports both pointer and non-pointer types.
// If the object is a pointer, it dereferences it before accessing the field.
// If the object is a struct, it directly accesses the field.
// If the field is not accessible or does not exist, the function returns nil.
//
// Parameters:
// - object: The object containing the field.
// - name: The name of the field to retrieve.
//
// Returns:
// - The value of the field with the specified name, or nil if the field is not accessible or does not exist.
func GetFieldValueBuName[T any](object T, name string) interface{} {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.FieldByName(name)
		if field.CanInterface() {
			return field.Interface()
		}
		return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
	} else if v.Kind() == reflect.Struct {
		val := v
		field := val.FieldByName(name)
		if field.CanInterface() {
			return field.Interface()
		}
		rs2 := reflect.New(val.Type()).Elem()
		rs2.Set(val)
		val = rs2.FieldByName(name)
		return reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem().Interface()
	}

	return nil
}

// SetFieldValueByIndex sets the value of a field at the specified index in the given object.
// The function supports both pointer and non-pointer types.
// If the object is a pointer, it dereferences it before accessing the field.
// If the object is a struct, it directly accesses the field.
// If the field is not accessible or does not exist, the function does nothing.
//
// Parameters:
// - object: The object containing the field. It must be of type T.
// - index: The index of the field to set.
// - value: The new value to assign to the field.
//
// Returns:
// - None
func SetFieldValueByIndex[T any](object T, index int, value interface{}) {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.Field(index)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
		}
	} else if v.Kind() == reflect.Struct {
		val := v
		field := val.Field(index)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
			object = val.Interface().(T)
		} else {
			rs2 := reflect.New(val.Type()).Elem()
			rs2.Set(val)
			val = rs2.Field(index)
			val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()

			val.Set(reflect.ValueOf(value))
		}
	}
}

// SetFieldValueByName sets the value of a field with the specified name in the given object.
// The function supports both pointer and non-pointer types.
// If the object is a pointer, it dereferences it before accessing the field.
// If the object is a struct, it directly accesses the field.
// If the field is not accessible or does not exist, the function does nothing.
//
// Parameters:
// - object: The object containing the field. It must be of type T.
// - name: The name of the field to set.
// - value: The new value to assign to the field.
//
// Returns:
// - None
func SetFieldValueByName[T any](object T, name string, value interface{}) {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.FieldByName(name)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
		}
	} else if v.Kind() == reflect.Struct {
		val := v
		field := val.FieldByName(name)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
			object = val.Interface().(T)
		} else {
			rs2 := reflect.New(val.Type()).Elem()
			rs2.Set(val)
			val = rs2.FieldByName(name)
			val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()

			val.Set(reflect.ValueOf(value))
		}
	}
}

// GetFieldValue retrieves the value of a given reflect.Value.
// If the field can be directly accessed, it returns the field value.
// If the field is not accessible, it creates a new value at the field's memory address and returns it.
//
// Parameters:
// - field: The reflect.Value representing the field to retrieve the value from.
//
// Returns:
// - reflect.Value: The value of the given field.
func GetFieldValue(field reflect.Value) reflect.Value {
	if field.CanInterface() {
		return field
	}
	res := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()

	return res
}

// SetFieldValue sets the value of a given reflect.Value.
// If the field can be directly accessed, it sets the field value.
// If the field is not accessible, it creates a new value at the field's memory address and sets it.
//
// Parameters:
// - field: The reflect.Value representing the field to set the value for.
// - value: The new value to assign to the field. It must be assignable to the field's type.
//
// Returns:
// - None
func SetFieldValue(field reflect.Value, value interface{}) {
	if field.CanInterface() && field.CanAddr() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
	}
}

// GetFieldValueFromMethodAndObject retrieves the value of a method with the specified name from the given object.
// The function supports both pointer and non-pointer types.
// If the object is a pointer, it dereferences it before accessing the method.
// If the object is a struct, it directly accesses the method.
// If the method is not accessible or does not exist, the function returns a zero reflect.Value.
//
// Parameters:
// - object: The object containing the method. It must be of type T.
// - name: The name of the method to retrieve the value from.
//
// Returns:
// - reflect.Value: The value of the given method. If the method does not exist or is not accessible, it returns a zero reflect.Value.
func GetFieldValueFromMethodAndObject[T interface{}](object T, name string) reflect.Value {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v
		method := val.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		}
	} else if v.Kind() == reflect.Struct {
		method := v.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		} else {
			pointerType := v.Addr()
			method := pointerType.MethodByName(name)
			res := method.Call(nil)
			return res[0]
		}
	}

	return *new(reflect.Value)
}

// GetFieldValueFromMethodAndReflectValue retrieves the value of a method with the specified name from the given reflect.Value.
// The function supports both pointer and non-pointer types.
// If the reflect.Value is a pointer, it dereferences it before accessing the method.
// If the reflect.Value is a struct, it directly accesses the method.
// If the method is not accessible or does not exist, the function returns a zero reflect.Value.
//
// Parameters:
// - val: The reflect.Value containing the method. It must be of type reflect.Value.
// - name: The name of the method to retrieve the value from.
//
// Returns:
// - reflect.Value: The value of the given method. If the method does not exist or is not accessible, it returns a zero reflect.Value.
func GetFieldValueFromMethodAndReflectValue(val reflect.Value, name string) reflect.Value {
	if val.Kind() != reflect.Ptr {
		method := val.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		}
	} else if val.Kind() == reflect.Struct {
		method := val.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		} else {
			pointerType := val.Addr()
			method = pointerType.MethodByName(name)
			res := method.Call(nil)
			return res[0]
		}
	}

	return *new(reflect.Value)
}

// SetValue sets the value of a given variable of type T to the provided value.
// It supports both pointer and non-pointer types. If the input data is a pointer,
// it dereferences it before setting the value. If the input data is a non-pointer,
// it directly sets the value.
//
// Parameters:
// - data: The variable of type T to set the value for. It can be a pointer or a non-pointer.
// - value: The new value to assign to the variable. It can be of any type.
//
// Returns:
// - None
func SetValue[T any](data T, value interface{}) {
	var inputValue reflect.Value
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		inputValue = reflect.ValueOf(data).Elem()
	} else {
		inputValue = reflect.ValueOf(data)
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() == reflect.Ptr {
		inputValue.Set(valueReflect.Elem())
	} else {
		inputValue.Set(valueReflect)
	}
}

// ObjectTypePath returns the fully qualified type path of the given object.
// It extracts the package path and type name from the object's type and returns them as a string.
//
// Parameters:
// - obj: The object for which to retrieve the type path. It can be of any type.
//
// Returns:
// - string: The fully qualified type path of the given object.
//
// Example:
//
//	type MyStruct struct {
//		Field1 string
//		Field2 int
//	}
//
//	myObj := MyStruct{}
//	typePath := ObjectTypePath(myObj)
//	fmt.Println(typePath) // Output: "main.MyStruct"
func ObjectTypePath(obj any) string {
	objType := reflect.TypeOf(obj).Elem()
	path := fmt.Sprintf("%s.%s", objType.PkgPath(), objType.Name())
	return path
}

// TypePath returns the fully qualified type path of the given generic type T.
// It extracts the package path and type name from the type's reflection and returns them as a string.
//
// Parameters:
// - T: A generic type for which to retrieve the type path. It can be of any type.
//
// Returns:
// - string: The fully qualified type path of the given generic type T.
//
// Example:
//
//	type MyStruct struct {
//		Field1 string
//		Field2 int
//	}
//
//	typePath := TypePath[MyStruct]()
//	fmt.Println(typePath) // Output: "main.MyStruct"
func TypePath[T any]() string {
	var msg T
	return ObjectTypePath(msg)
}

// CreateInstance creates a new instance of the given generic type T.
// It uses reflection to dynamically create a new instance of the type.
//
// Parameters:
// - T: A generic type for which to create a new instance. It can be of any type.
//
// Returns:
// - T: A new instance of the given generic type T.
func CreateInstance[T any]() T {
	var msg T
	typ := reflect.TypeOf(msg).Elem()
	instance := reflect.New(typ).Interface()
	return instance.(T)
}

// MethodPath returns the fully qualified method path of the given function.
// It extracts the package path, receiver type, and method name from the function's reflection and returns them as a string.
//
// Parameters:
// - f: The function for which to retrieve the method path. It can be of any type.
//
// Returns:
// - string: The fully qualified method path of the given function.
//
// Example:
//
//	func (s *MyStruct) MyMethod() {
//		// ...
//	}
//
//	func main() {
//		myStruct := &MyStruct{}
//		methodPath := MethodPath(myStruct.MyMethod)
//		fmt.Println(methodPath) // Output: "main.(*MyStruct).MyMethod"
//	}
func MethodPath(f interface{}) string {
	pointerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lastSlashIndex := strings.LastIndex(pointerName, "/")
	methodPath := pointerName[lastSlashIndex+1:]
	if methodPath[len(methodPath)-3:] == "-fm" {
		methodPath = methodPath[:len(methodPath)-3]
	}
	methodPath = strings.ReplaceAll(methodPath, ".", ":")
	methodPath = strings.ReplaceAll(methodPath, "(", "")
	methodPath = strings.ReplaceAll(methodPath, ")", "")
	methodPath = strings.ReplaceAll(methodPath, "*", "")
	return methodPath
}
