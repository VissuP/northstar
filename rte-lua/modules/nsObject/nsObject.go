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

package nsObject

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/object/client"
	"github.com/verizonlabs/northstar/object/model"
	"strconv"
	"time"
)

type NsObjectModule struct {
	Client    *client.ObjectClient
	AccountId string
}

func NewNsObjectModule(accountid string) (*NsObjectModule, error) {
	cli, err := client.NewObjectClient()
	if err != nil {
		return nil, err
	}
	return &NsObjectModule{Client: cli, AccountId: accountid}, nil
}

func (nsObject *NsObjectModule) Loader(L *lua.LState) int {
	api := map[string]lua.LGFunction{
		"createBucket": nsObject.createBucket,
		"deleteBucket": nsObject.deleteBucket,
		"listBuckets":  nsObject.listBuckets,
		"uploadFile":   nsObject.uploadFile,
		"downloadFile": nsObject.downloadFile,
		"deleteFile":   nsObject.deleteFile,
		"listFiles":    nsObject.listFiles,
	}
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

func (nsObject *NsObjectModule) createBucket(L *lua.LState) int {
	bucketName := L.CheckString(1)

	bucket := &model.Bucket{Name: bucketName, CreationDate: time.Now()}
	if _, mErr := nsObject.Client.CreateBucket(nsObject.AccountId, bucket); mErr != nil {
		errM := fmt.Sprintf("Failed to create bucket: %s", bucket.Name)
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "createBucket")
	}

	CreateBucket.Incr()
	return 0
}

func (nsObject *NsObjectModule) deleteBucket(L *lua.LState) int {
	bucketName := L.CheckString(1)

	if mErr := nsObject.Client.DeleteBucket(nsObject.AccountId, bucketName); mErr != nil {
		errM := fmt.Sprintf("Failed to delete bucket: %s", bucketName)
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "deleteBucket")
	}

	DeleteBucket.Incr()
	return 0
}

func (nsObject *NsObjectModule) listBuckets(L *lua.LState) int {
	buckets, mErr := nsObject.Client.ListBuckets(nsObject.AccountId)
	if mErr != nil {
		errM := fmt.Sprintf("Failed to list buckets")
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "listBuckets")
	}

	arr := L.CreateTable(len(buckets), 0)
	for _, bucket := range buckets {
		tbl := L.CreateTable(0, 1)
		tbl.RawSetH(lua.LString("name"), lua.LString(bucket.Name))
		tbl.RawSetH(lua.LString("date"), lua.LString(bucket.CreationDate.String()))
		arr.Append(tbl)
	}

	L.Push(arr)
	ListBuckets.Incr()
	return 1
}

func (nsObject *NsObjectModule) uploadFile(L *lua.LState) int {
	bucketName := L.CheckString(1)
	fileName := L.CheckString(2)
	input := L.CheckAny(3)
	contentType := L.CheckString(4)

	var data []byte
	switch v := input.(type) {
	case lua.LString:
		data = []byte(string(v))
	case *lua.LTable:
		size := v.MaxN()
		if size == 0 {
			return nsObject.error(L, "unexpected table", nil, "uploadFile")
		}

		for i := 1; i <= size; i++ {
			value := v.RawGetInt(i)
			switch val := value.(type) {
			case lua.LNumber:
				data = append(data, byte(val))
			default:
				return nsObject.error(L, "unexpected value in array, byte expected", nil, "uploadFile")
			}
		}
	default:
		return nsObject.error(L, "unexpected value, string or byte array expected", nil, "uploadFile")
	}

	uploadData := &model.UploadData{FileName: fileName, Payload: data, ContentType: contentType}
	_, mErr := nsObject.Client.UploadFile(nsObject.AccountId, bucketName, uploadData)
	if mErr != nil {
		errM := fmt.Sprintf("Failed to upload file %s", fileName)
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "uploadFile")
	}

	UploadFile.Incr()
	return 0
}

func (nsObject *NsObjectModule) downloadFile(L *lua.LState) int {
	bucketName := L.CheckString(1)
	fileName := L.CheckString(2)

	data, mErr := nsObject.Client.DownloadFile(nsObject.AccountId, bucketName, fileName)
	if mErr != nil {
		errM := fmt.Sprintf("Failed to download file %s", fileName)
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "downloadFile")
	}

	payload := L.CreateTable(len(data.Payload), 0)
	for _, d := range data.Payload {
		payload.Append(lua.LNumber(d))
	}

	output := L.CreateTable(0, 2)
	output.RawSetH(lua.LString("Payload"), payload)
	output.RawSetH(lua.LString("ContentType"), lua.LString(data.ContentType))

	L.Push(output)
	DownloadFile.Incr()
	return 1
}

func (nsObject *NsObjectModule) deleteFile(L *lua.LState) int {
	bucketName := L.CheckString(1)
	fileName := L.CheckString(2)

	mErr := nsObject.Client.DeleteFile(nsObject.AccountId, bucketName, fileName)
	if mErr != nil {
		errM := fmt.Sprintf("Failed to delete file %s", fileName)
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "deleteFile")
	}

	DeleteFile.Incr()
	return 0
}

func (nsObject *NsObjectModule) listFiles(L *lua.LState) int {
	bucketName := L.CheckString(1)

	objects, mErr := nsObject.Client.ListFiles(nsObject.AccountId, bucketName)
	if mErr != nil {
		errM := fmt.Sprintf("Failed to list files")
		mlog.Error(mErr.Error())
		return nsObject.error(L, errM, nil, "listFiles")
	}

	arr := L.CreateTable(len(objects), 0)
	for _, object := range objects {
		tbl := L.CreateTable(0, 1)
		tbl.RawSetH(lua.LString("key"), lua.LString(object.Key))
		tbl.RawSetH(lua.LString("last_modified"), lua.LString(object.LastModified.String()))
		tbl.RawSetH(lua.LString("size"), lua.LString(strconv.FormatInt(object.Size, 10)))
		tbl.RawSetH(lua.LString("etag"), lua.LString(object.Etag))
		tbl.RawSetH(lua.LString("storage_class"), lua.LString(object.StorageClass))
		arr.Append(tbl)
	}

	L.Push(arr)
	ListFiles.Incr()
	return 1
}
