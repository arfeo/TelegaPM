package main

import (
	"encoding/json"
	"log"
	"strconv"
)

/**
 *
 *	Get storage data from encrypted content
 *
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
 *
 *	Push a new element to storage
 *
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
 *
 *	Get elements list from storage
 *
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
 *
 *	Get element info
 *
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
 *
 *	Remove an element from storage
 *
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
