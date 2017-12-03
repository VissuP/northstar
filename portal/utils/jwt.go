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
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"time"
)

func GetJWTFromString(reqJwt string) (*jwt.Token, error) {
	mlog.Debug("JWT String: %s", reqJwt)
 	token, parseErr := jwt.Parse(reqJwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("4xGKO2rHoLp4vbMW9z08Ql8sCEL1inv0"), nil
	})
	if parseErr != nil {
		mlog.Debug("Failed to parse JWT")
		mlog.Debug(parseErr.Error())
		return nil, fmt.Errorf("Failed to parse JWT")
	}
	expiration := token.Claims.(jwt.MapClaims)["exp"]
	if expiration == nil {
		return nil, fmt.Errorf("Token Expiry is Nil")
	}
	if expiration.(float64) < float64(time.Now().Unix()) {
		mlog.Debug("JWT validation failed - token is expired")
		return nil, fmt.Errorf("JWT validation failed - token is expired")
	}
	return token, nil
}

func ValidateTSCoreJWT(jwtObj *jwt.Token) bool{
	if jwtObj.Claims.(jwt.MapClaims)["email"] ==nil || jwtObj.Claims.(jwt.MapClaims)["token"]==nil{
		mlog.Debug("Validate TSCore Jwt Failed...")
		return false
	}
	return true
}

// build JWT containing token and email, and expiring when ts core token expires
func CreateJWTFromToken(token string, email string, cookieDurationSeconds int) (string, error) {
	mlog.Debug("Creating JWT...")
	replaceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token":      token,
		"email":      email,
		"exp":        time.Now().Unix() + int64(cookieDurationSeconds),
	})
	tokenString, err := replaceToken.SignedString([]byte("4xGKO2rHoLp4vbMW9z08Ql8sCEL1inv0"))
	if err != nil {
		return "", err
	} else {
		mlog.Debug(tokenString)
		return tokenString, nil
	}
}
