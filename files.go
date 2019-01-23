package main

import (
	"io/ioutil"
	"log"
	"os"
)

/**
 *
 *	Check whether a file or directory exists
 *
 */

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

/**
 *
 *	Create a new file or directory
 *
 */

func CreateFile(name string) bool {
	_, err := os.Create(name)

	if err != nil {
		log.Println(err)

		return false
	}

	return true
}

/**
 *
 *	Remove a file
 *
 */

func RemoveFile(name string) bool {
	err := os.Remove(name)

	if err != nil {
		log.Println(err)

		return false
	}

	return true
}

/**
 *
 *	Read data from a file
 *
 */

func ReadFromFile(name string) ([]byte, bool) {
	file, err := ioutil.ReadFile(name)

	if err != nil {
		log.Println(err)

		return nil, false
	}

	return file, true
}

/**
 *
 *	Write data to a file
 *
 */

func WriteToFile(name string, data []byte) bool {
	err := ioutil.WriteFile(name, data, 0644)

	if err != nil {
		log.Println(err)

		return false
	}

	return true
}
