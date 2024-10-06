package reflectionhelper

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// Helper function to get the addressable value of a field.
func getAddressableValue(field reflect.Value) reflect.Value {
	if field.CanInterface() {
		return field
	}
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}

// GetFieldValueByIndex retrieves the value of a field by its index from the given object.
func GetFieldValueByIndex[T any](object T, index int) interface{} {
	val := reflect.ValueOf(&object).Elem()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.Field(index)
	return getAddressableValue(field).Interface()
}

// GetFieldValueByName retrieves the value of a field by its name from the given object.
func GetFieldValueByName[T any](object T, name string) interface{} {
	val := reflect.ValueOf(&object).Elem()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(name)
	return getAddressableValue(field).Interface()
}

// SetFieldValueByIndex sets the value of a field by its index in the given object.
func SetFieldValueByIndex[T any](object T, index int, value interface{}) {
	val := reflect.ValueOf(&object).Elem()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.Field(index)
	if field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		getAddressableValue(field).Set(reflect.ValueOf(value))
	}
}

// SetFieldValueByName sets the value of a field by its name in the given object.
func SetFieldValueByName[T any](object T, name string, value interface{}) {
	val := reflect.ValueOf(&object).Elem()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(name)
	if field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		getAddressableValue(field).Set(reflect.ValueOf(value))
	}
}

// GetFieldValue retrieves the value of a field that might not be directly accessible.
func GetFieldValue(field reflect.Value) reflect.Value {
	return getAddressableValue(field)
}

// SetFieldValue sets the value of a field.
func SetFieldValue(field reflect.Value, value interface{}) {
	if field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		getAddressableValue(field).Set(reflect.ValueOf(value))
	}
}

// GetFieldValueFromMethodAndObject retrieves the value by invoking the method from the given object.
func GetFieldValueFromMethodAndObject[T interface{}](object T, name string) reflect.Value {
	v := reflect.ValueOf(&object).Elem()

	var method reflect.Value

	if v.Kind() == reflect.Ptr {
		method = v.MethodByName(name)
	} else if v.Kind() == reflect.Struct {
		method = v.MethodByName(name)
		if method.Kind() != reflect.Func {
			vPtr := v.Addr()
			method = vPtr.MethodByName(name)
		}
	}

	if method.Kind() == reflect.Func {
		results := method.Call(nil)
		if len(results) > 0 {
			return results[0]
		}
	}

	// Return an invalid reflect.Value if method not found or doesn't return anything
	return reflect.Value{}
}

// GetFieldValueFromMethodAndReflectValue retrieves the value by invoking the method from the given reflect value.
func GetFieldValueFromMethodAndReflectValue(val reflect.Value, name string) reflect.Value {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return getFieldValueFromMethod(val, name)
}

// Helper function for method invocation and obtaining result.
func getFieldValueFromMethod(val reflect.Value, name string) reflect.Value {
	method := val.MethodByName(name)
	if method.Kind() == reflect.Func {
		res := method.Call(nil)
		if len(res) > 0 {
			return res[0]
		}
	}
	return reflect.Value{}
}

// SetValue assigns a value to a given data reference.
func SetValue[T any](data T, value interface{}) {
	dataVal := reflect.ValueOf(&data).Elem() // Use a pointer to the data for setting its value
	if dataVal.Kind() == reflect.Ptr {
		if dataVal.IsNil() {
			// Initialize nil pointer
			dataVal.Set(reflect.New(dataVal.Type().Elem()))
		}
		dataVal = dataVal.Elem()
	}

	valueVal := reflect.ValueOf(value)
	if valueVal.Kind() == reflect.Ptr && !valueVal.IsNil() {
		valueVal = valueVal.Elem()
	}

	if valueVal.IsValid() {
		dataVal.Set(valueVal)
	} else {
		// Set to zero value if the value is not valid (e.g., nil)
		dataVal.Set(reflect.Zero(dataVal.Type()))
	}
}

// TypePath returns the type path of a given generic type T.
func TypePath[T any]() string {
	var msg T
	return ObjectTypePath(msg)
}

// ObjectTypePath returns the type path of an object.
func ObjectTypePath(obj any) string {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}
	return fmt.Sprintf("%s.%s", objType.PkgPath(), objType.Name())
}

// CreateInstance creates a new instance of the given generic type T.
func CreateInstance[T any]() T {
	var t T
	tType := reflect.TypeOf(t)

	if tType.Kind() == reflect.Ptr {
		return reflect.New(tType.Elem()).Interface().(T)
	}

	// Directly create an instance for non-pointer types.
	instance := reflect.New(tType).Elem().Interface()
	return instance.(T)
}

// MethodPath returns the method path for a given function.
func MethodPath(f interface{}) string {
	pointerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lastSlashIndex := strings.LastIndex(pointerName, "/")
	methodPath := pointerName[lastSlashIndex+1:]
	if strings.HasSuffix(methodPath, "-fm") {
		methodPath = methodPath[:len(methodPath)-3]
	}
	replacer := strings.NewReplacer(".", ":", "(", "", ")", "", "*", "")
	return replacer.Replace(methodPath)
}
