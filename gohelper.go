package gohelper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
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

/**
 * SendRequest is used to send http request
 * Params:
 * 		string 				method 			http methods (GET | POST | PUT | DELETE)
 * 		string				uri				http url
 * 		...interface{}		args			optional http parameters (headers, request body)
 * 	Return:
 * 		[]byte				responseBody
 * 		int					statusCode
 * 		error				err
 */
func SendRequest(method string, uri string, args ...interface{}) (responseBody []byte, statusCode int, err error) {
	method = strings.ToUpper(method)

	switch method {
	case "GET":
		resp, statusCode, err := sendGet(method, uri, args...)
		return resp, statusCode, err

	case "POST":
		resp, statusCode, err := sendPost(method, uri, args...)
		return resp, statusCode, err

	default:
		err := errors.New("Request method is not supported")
		return nil, http.StatusMethodNotAllowed, err
	}

}

/**
 * sendGet is used to send http GET request
 * Params:
 * 		string 				method 			GET http methods
 * 		string				uri				http url
 * 		...interface{}		args			optional http parameters (headers, request body)
 * 	Return:
 * 		[]byte				responseBody
 * 		int					statusCode
 * 		error				err
 */
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
		return nil, resp.StatusCode, err
	}

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error when reading response body: %+v\n", err.Error())
		return nil, resp.StatusCode, err
	}

	return bodyResp, resp.StatusCode, nil
}

/**
 * sendPost is used to send http POST request
 * Params:
 * 		string 				method 			POST http methods
 * 		string				uri				http url
 * 		...interface{}		args			optional http parameters (headers, request body)
 * 	Return:
 * 		[]byte				responseBody
 * 		int					statusCode
 * 		error				err
 */
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
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Response error. Error:  %+v\n", err.Error())
		return nil, resp.StatusCode, err
	}

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error when reading response body: %+v\n", err.Error())
		return nil, resp.StatusCode, err
	}

	return bodyResp, resp.StatusCode, nil
}

/**
 * GetEnv gets the environment variable. If environment variable is not set,
 * it returns the fallback.
 *
 * Params:
 * 		string 		key
 * 		string		fallback		Default env variable if env with "key" is not set
 *
 * Returns:
 * 		string		env
 */
func GetEnv(key string, fallback string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		env = fallback
	}

	return env
}
