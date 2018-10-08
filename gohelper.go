package gohelper

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

// sendGet sends HTTP GET request to the specified URI. Args can be url.Values and http.Header type.
func sendGet(method string, uri string, args ...interface{}) (responseBody []byte, statusCode int, err error) {
	client := &http.Client{}

	req, _ := http.NewRequest(method, uri, nil)

	for _, arg := range args {
		switch val := arg.(type) {
		case url.Values:
			// Assign http query to the request URL
			req.URL.RawQuery = val.Encode()

		case http.Header:
			// Assign request header
			req.Header = val
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Response error. Error:  %+v\n", err.Error())
		return nil, 500, err
	}

	defer resp.Body.Close()

	// Read response body
	statusCode = resp.StatusCode

	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error when converting response body to []byte: %+v\n", err.Error())
		return nil, statusCode, err
	}

	return result, statusCode, nil
}

// sendPost sends HTTP POST request to the specified URI. Args can be url.Values, http.Header, and []byte type.
func sendPost(method string, uri string, args ...interface{}) (responseBody []byte, statusCode int, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, uri, nil)

	for _, arg := range args {
		switch val := arg.(type) {
		case url.Values:
			// Assign request body
			req.Body = ioutil.NopCloser(strings.NewReader(val.Encode()))

		case http.Header:
			// Assign request header
			req.Header = val

		case []byte:
			// Assign request body
			data := bytes.NewBuffer(val)
			req.Body = ioutil.NopCloser(data)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Response error. Error:  %+v\n", err.Error())

		return nil, 500, err
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode

	// Read response body
	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error when reading response body: %+v\n", err.Error())
		return nil, statusCode, err
	}

	return result, statusCode, nil
}

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

// SendRequest sends HTTP GET, POST, or PUT request to the specified URI. Args can be url.Values, http.Header, and []byte type.
func SendRequest(method string, uri string, args ...interface{}) (responseBody []byte, statusCode int, err error) {
	method = strings.ToUpper(method)

	switch method {
	case "GET":
		resp, statusCode, err := sendGet(method, uri, args...)
		return resp, statusCode, err

	case "POST", "PUT":
		resp, statusCode, err := sendPost(method, uri, args...)
		return resp, statusCode, err

	default:
		err := errors.New("Request method is not supported")
		return nil, http.StatusMethodNotAllowed, err
	}
}

// GetEnv gets the environment variable. If environment variable is not set,
// it returns the fallback.
func GetEnv(key string, fallback string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		env = fallback
	}

	return env
}

// InArray checks whether needle is in haystack.
func InArray(needle interface{}, haystack interface{}) (bool, int, error) {
	haystackValue := reflect.ValueOf(haystack)
	haystackType := haystackValue.Type()

	if haystackType.Kind() != reflect.Array && haystackType.Kind() != reflect.Slice {
		err := errors.New("Parameter 2 is not an array or slice")
		return false, -1, err
	}

	switch haystackType.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < haystackValue.Len(); i++ {
			hayVal := haystackValue.Index(i).Interface()

			if reflect.DeepEqual(hayVal, needle) {
				return true, i, nil
			}
		}
	}

	return false, -1, nil
}
