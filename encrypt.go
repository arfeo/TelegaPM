package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

/**
 *
 *	Encrypt with specified key (32 bytes)
 *
 */

func Encrypt(key, text []byte) []byte {
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	b := EncodeBase64(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return ciphertext
}

/**
 *
 *	Decrypt with specified key (32 bytes)
 *
 */

func Decrypt(key, text []byte) string {
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	if len(text) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	return string(DecodeBase64(string(text)))
}

/**
 *
 *	Base 64 encode
 *
 */

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

/**
 *
 *	Base 64 decode
 *
 */

func DecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		log.Println(err)
	}

	return data
}

/**
 *
 *	Calculate the sha1 hash of a string
 *
 */

func Sha1(b []byte) string {
	hash := sha1.New()
	hash.Write(b)

	return fmt.Sprintf("%x", hash.Sum(nil))
}
