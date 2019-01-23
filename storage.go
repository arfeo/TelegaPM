package main

import (
	"encoding/json"
	"log"
	"strconv"
)

/**
 * Function tries to get storage data from encrypted content with the given hash;
 * returns storage date in case of success, otherwise returns void
 */
func GetStorage(data []byte, hash [32]byte) (storage []Storage) {
	if len(data) > 0 {
		decr := Decrypt(hash[:], data)
		err := json.Unmarshal([]byte(decr), &storage)

		if err != nil {
			log.Println(err)
		}
	}

	return
}

/**
 * Function pushes a new element to storage with the given hash;
 * returns encrypted storage data
 */
func PushToStorage(storage []Storage, current Current, hash [32]byte) []byte {
	storage = append(storage, current.Entry)
	j, err := json.Marshal(storage)

	if err != nil {
		log.Println(err)
	}

	return Encrypt(hash[:], j)
}

/**
 * Function tries to get elements list from storage;
 * returns the list as the first value and true as the second one if the list is not empty;
 * returns the corresponding message as the first value and false as the second one if the list is empty
 */
func StorageList(storage []Storage) (string, bool) {
	result := ""

	if len(storage) > 1 {
		for i, el := range storage[1:] {
			result += strconv.Itoa(i+1) + ".\t" + el.Title + "\n"
		}

		return result, true
	}

	return "No elements.", false
}

/**
 * Function tries to get an element info from storage by the given element number;
 * returns element's info as the first value, password as the second value, and true as the third one if element exists;
 * returns the corresponding message as the first value, empty string as the second value, and false if element not found
 */
func GetStorageElementInfo(storage []Storage, num int) (string, string, bool) {
	result := ""

	if num > 0 && num <= len(storage) {
		result += "Title:\t<code>" + storage[num].Title + "</code>\n" + "Login:\t<code>" + storage[num].Login + "</code>\n" + "URL:\t<code>" + storage[num].Url + "</code>\n"

		return result, storage[num].Pass, true
	}

	return "Element not found.", "", false
}

/**
 * Function tries to remove an element from storage by the given number with the given hash;
 * returns encrypted storage date  as the first value and true as the second one in case of success;
 * otherwise returns empty string as the first value and false as the second one
 */
func RemoveElementFromStorage(storage []Storage, num int, hash [32]byte) ([]byte, bool) {
	if num > 0 && num <= len(storage) {
		storage = append(storage[:num], storage[num+1:]...)
		j, err := json.Marshal(storage)

		if err != nil {
			log.Println(err)
		}

		return Encrypt(hash[:], j), true
	} else {
		return []byte(""), false
	}
}
