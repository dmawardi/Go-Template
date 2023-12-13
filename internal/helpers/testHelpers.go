package helpers

import (
	"fmt"
	"reflect"
)

// Update a struct field dynamically
func UpdateStructField(structPtr interface{}, fieldName string, fieldValue interface{}) error {
	value := reflect.ValueOf(structPtr)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("invalid struct pointer")
	}

	structValue := value.Elem()
	if !structValue.CanSet() {
		return fmt.Errorf("cannot set struct field value")
	}

	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("invalid struct field name")
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set struct field value")
	}

	fieldValueRef := reflect.ValueOf(fieldValue)
	if !fieldValueRef.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("field value type mismatch")
	}

	field.Set(fieldValueRef)
	return nil
}
