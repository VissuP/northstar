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
	r "crypto/rand"
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"math/rand"
	"github.com/verizonlabs/northstar/pkg/management"
	"strings"
	"time"
)

const (
	strlen   = 16
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numerals = "0123456789"
)

func GetSecret() (string, error) {

	key := make([]byte, 32)
	_, err := r.Read(key)

	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(key), nil
}

func GetQrCode() string {

	charSet := []string{Alphabet, Numerals}
	chars := strings.Join(charSet, "")

	result := make([]byte, strlen)
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

// Helper method used to render http response from given management error object.
func RenderServiceError(context *gin.Context, serviceError *management.Error) {
	// Note that headers need to be set before calling context.JSON method
	for k, v := range serviceError.Header {
		for _, v1 := range v {
			context.Writer.Header().Add(k, v1)
		}
	}

	// Serialize rest of the response
	context.JSON(serviceError.HttpStatus, serviceError)
}
