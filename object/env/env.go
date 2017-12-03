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

package env

import (
	"os"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"strings"
)

var (
	WebPort, _             = config.GetString("OBJECT_PORT", "80")
	DebugFlag, _           = config.GetBool("ENABLE_DEBUG", false)
	BlobStorageHostAndPort = getStorageHostPort()
	BlobStorageUserId      = os.Getenv("OBJECT_BLOB_STORAGE_USER_ID")
	BlobStorageUserSecret  = os.Getenv("OBJECT_BLOB_STORAGE_USER_SECRET")
)

func getStorageHostPort() string {
	val := os.Getenv("BLOB_STORAGE_HOST_AND_PORT")
	mlog.Info("BLOB_STORAGE_HOST_AND_PORT is %s", val)
	if val == "" {
		mlog.Alarm("Env variable BLOB_STORAGE_HOST_AND_PORT is not set")
	}

	// BLOB_STORAGE_HOST_AND_PORT can have entries
	// host:port
	// :port (in this case we should use $HOST)
	r := strings.Split(val, ":")
	mlog.Info("Value after split is %s", r)

	if r[0] == "" {
		ip := os.Getenv("HOST")
		iphost := ip + val
		mlog.Info("BLOB_STORAGE_HOST_AND_PORT is %s", iphost)
		return iphost
	}
	return val
}
