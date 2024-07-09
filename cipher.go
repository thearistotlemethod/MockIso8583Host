/*
 *   Copyright (c) 2024 thearistotlemethod@gmail.com
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

const msgSessionKey = "3D138FE95179326243A14FE5BFDADACE"
const onlinePinKey = "F249EFA2945B8A3B32BA4C9E4929451F"

func DesEncryption(plainText []byte) ([]byte, error) {
	var iv = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	key, _ := hex.DecodeString(msgSessionKey)
	if len(key) == 16 {
		key = append(key, key[0:8]...)
	}

	block, err := des.NewTripleDESCipher(key)

	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	plainText = PKCS5Padding(plainText, bs)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(plainText))
	blockMode.CryptBlocks(cryted, plainText)
	return cryted, nil
}

func DesDecryption(cipherText []byte) ([]byte, error) {
	var iv = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	key, _ := hex.DecodeString(msgSessionKey)
	if len(key) == 16 {
		key = append(key, key[0:8]...)
	}

	block, err := des.NewTripleDESCipher(key)

	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)

	origData = PKCS5Unpadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
