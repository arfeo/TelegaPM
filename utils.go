package main

import (
	"reflect"
)

/**
 * Function checks if the given value exists in the given array;
 * returns true as the first value and value's index as the second one if value exists,
 * otherwise returns false as the first value and -1 as the second one
 */
func ValueInArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice :
		{
			s := reflect.ValueOf(array)

			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
					index = i
					exists = true

					return
				}
			}
		}
	}

	return
}