package lib

import (
	"fmt"
	"os"
	"reflect"
)

// CheckErr is a re-usable function for checking errors
func CheckErr(msg string, e error) {
	if e != nil {
		fmt.Printf("%v: %v\n", msg, e)
		os.Exit(1)
	}
}

// Contains function uses reflection to iterate over a given slice (s) and see if a value (v) is contained within it
func Contains(s interface{}, v interface{}) bool {
	arr := reflect.ValueOf(s)
	if arr.Kind() == reflect.Slice {
		for i := 0; i < arr.Len(); i++ {
			if arr.Index(i).Interface() == v {
				return true
			}
		}
	}
	return false
}
