/*
Copyright (C) 2017 Verizon. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	cryprand "crypto/rand"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"fmt"
)

var encryptionKey = "K92j0jvGRar9mt8wQqV65lB28ELMHVT1"

var MacLength int = 32
func EncryptAndB64EncodeAccessToken(token string) (string, error) {

	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, MacLength+aes.BlockSize+len(token))
	iv := ciphertext[MacLength:MacLength+aes.BlockSize]
	if _, err := io.ReadFull(cryprand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(c, iv)
	cfb.XORKeyStream(ciphertext[MacLength+aes.BlockSize:], []byte(token))

	hashmac := ciphertext[:MacLength]
	copy(hashmac[:],ComputeHmac256(ciphertext[MacLength+aes.BlockSize:],[]byte(encryptionKey)))
	mlog.Debug("Encrypting: %s=>%x", []byte(token), ciphertext)
	sEnc := b64.StdEncoding.EncodeToString(ciphertext)
	return sEnc, nil

}

func B64DecodeAndDecryptAccessToken(encryptedToken string) (string, error) {

	sDec, _ := b64.StdEncoding.DecodeString(encryptedToken)
	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}
	tokenMac := sDec[:MacLength]
	expectedMac := ComputeHmac256(sDec[MacLength+aes.BlockSize:],[]byte(encryptionKey))
	if !hmac.Equal(tokenMac,expectedMac){
		return "", fmt.Errorf("Invalid token")
	}
	iv := sDec[MacLength:MacLength+aes.BlockSize]
	cfbdec := cipher.NewCFBDecrypter(c, iv)
	decryptedToken := make([]byte, len(sDec) - MacLength - aes.BlockSize)
	cfbdec.XORKeyStream(decryptedToken, []byte(sDec[MacLength+aes.BlockSize:]))
	mlog.Debug("Decrypting: %x=>%s", sDec, decryptedToken)
	return string(decryptedToken), nil

}

func ComputeHmac256(message []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return h.Sum(nil)
}
