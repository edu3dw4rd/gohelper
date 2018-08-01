package gohelper

import (
	"errors"
	"reflect"
)

// RemoveAtIndex removes element from array data at index.
func RemoveAtIndex(data interface{}, index int) (interface{}, error) {
	// Get concrete value of data
	value := reflect.ValueOf(data)

	// Get the type of value
	valueType := value.Type()

	if valueType.Kind() != reflect.Array && valueType.Kind() != reflect.Slice {
		err := errors.New("Data parameter is not an array or slice")
		return nil, err
	}

	if index >= value.Len() {
		err := errors.New("Index is greater than data length")
		return nil, err
	}

	// Create slice from value
	resultSlice := reflect.AppendSlice(value.Slice(0, index), value.Slice(index+1, value.Len()))

	return resultSlice.Interface(), nil
}
