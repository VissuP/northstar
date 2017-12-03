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

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/object/model"
	"github.com/verizonlabs/northstar/object/s3"
	"github.com/verizonlabs/northstar/object/stats"
	"github.com/verizonlabs/northstar/object/util"
)

type Controller struct {
	StorageProvider s3.StorageProvider
}

func NewController(sProvider s3.StorageProvider) (controller *Controller) {
	mlog.Debug("Creating controller")
	controller = &Controller{StorageProvider: sProvider}
	return controller
}

func (controller *Controller) UploadFile(c *gin.Context) {
	mlog.Debug("Uploadfile starts")
	stats.UploadFileReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	bucketName := c.Params.ByName("bucketName")

	if accountId == "" {
		stats.ErrUploadFileMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	if bucketName == "" {
		stats.ErrUploadFileMissingBucketNameCount.Incr()
		errM := fmt.Sprint(util.BucketNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	var data = new(model.UploadData)
	c.Bind(data)

	err := data.Validate()
	if err != nil {
		stats.ErrUploadValidateCount.Incr()
		mlog.Error("Failed to validate upload data: %v", err)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(err.Error()))
		return
	}

	mlog.Debug("UploadFile(): fileName: %v", data.FileName)
	bucketName = getBucketName(accountId, bucketName)
	mErr := controller.StorageProvider.Upload(bucketName, data)
	if mErr != nil {
		stats.ErrUploadWriteFailCount.Incr()
		mlog.Error("Uploadfile():  failed due to %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	stats.ErrUploadWriteDataCount.Incr()
	c.JSON(http.StatusOK, "Data received")
	mlog.Debug("Uploadfile ends")
}

func (controller *Controller) DownloadFile(c *gin.Context) {
	mlog.Debug("DownloadFile starts")
	stats.DownloadFileReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	bucketName := c.Params.ByName("bucketName")
	fileName := c.Params.ByName("fileName")

	mlog.Debug("DownloadFile(): : %v", fileName)

	if accountId == "" {
		stats.ErrDownloadMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	if bucketName == "" {
		stats.ErrDownloadFileMissingBucketNameCount.Incr()
		errM := fmt.Sprint(util.BucketNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	if fileName == "" {
		stats.ErrDownloadMissingFileNameCount.Incr()
		errM := fmt.Sprint(util.FileNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	bucketName = getBucketName(accountId, bucketName)
	data, mErr := controller.StorageProvider.Download(bucketName, fileName)
	if mErr != nil {
		stats.ErrDownloadReadFailCount.Incr()
		mlog.Error("Downloadfile():  failed due to %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	byteArr, err := json.Marshal(data)
	if err != nil {
		mlog.Error("Downloadfile():  error marshaling data %v", err)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	c.Data(http.StatusOK, "application/json", byteArr)
	mlog.Debug("DownloadFile ends")
}

func (controller *Controller) DeleteFile(c *gin.Context) {
	mlog.Debug("DeleteFile starts")
	stats.DeleteFileReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	bucketName := c.Params.ByName("bucketName")
	fileName := c.Params.ByName("fileName")

	if accountId == "" {
		stats.ErrDeleteFileMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	if bucketName == "" {
		stats.ErrDeleteFilesMissingBucketNameCount.Incr()
		errM := fmt.Sprint(util.BucketNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	if fileName == "" {
		stats.ErrDeleteFileMissingFileNameCount.Incr()
		errM := fmt.Sprint(util.FileNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	fileName = strings.TrimLeft(fileName, "/")
	bucketName = getBucketName(accountId, bucketName)
	mlog.Debug("DeleteFile(): : %s/%s", bucketName, fileName)
	mErr := controller.StorageProvider.Delete(bucketName, fileName)
	if mErr != nil {
		stats.ErrDeleteFileFailCount.Incr()
		mlog.Error("Downloadfile():  failed due to %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	c.JSON(http.StatusOK, "Data deleted")
	mlog.Debug("DeleteFile ends")
}

func (controller *Controller) ListFiles(c *gin.Context) {
	mlog.Debug("ListFiles starts")
	stats.ListFilesReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	bucketName := c.Params.ByName("bucketName")

	if accountId == "" {
		stats.ErrDeleteFileMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	if bucketName == "" {
		stats.ErrDeleteFilesMissingBucketNameCount.Incr()
		errM := fmt.Sprint(util.BucketNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	bucketName = getBucketName(accountId, bucketName)
	files, _, mErr := controller.StorageProvider.List(bucketName, "")
	mlog.Debug("Data from status file is data=%v", files)
	if mErr != nil {
		stats.ErrListFilesCount.Incr()
		mlog.Error("Downloadfile(): failed due to %v", mErr.Error())
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	c.JSON(http.StatusOK, files)
	mlog.Debug("ListFiles ends")
}

func (controller *Controller) CreateBucket(c *gin.Context) {
	mlog.Debug("CreateBucket starts")
	stats.CreateBucketReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	var bucket = new(model.Bucket)
	c.Bind(bucket)

	if accountId == "" {
		stats.ErrCreateBucketMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	if bucket.Name == "" {
		stats.ErrCreateBucketMissingBucketNameCount.Incr()
		errM := fmt.Sprint(util.BucketNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	bucketName := getBucketName(accountId, bucket.Name)
	mErr := controller.StorageProvider.CreateBucket(bucketName)
	if mErr != nil {
		stats.ErrCreateBucketCount.Incr()
		mlog.Error("CreateBucket(): failed due to %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	mlog.Debug("CreateBucket ends")
	c.JSON(http.StatusOK, "Data received")
}

func (controller *Controller) ListBuckets(c *gin.Context) {
	mlog.Debug("ListBuckets starts")
	stats.ListBucketsReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		stats.ErrListBucketsMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	buckets, mErr := controller.StorageProvider.ListBuckets(accountId)
	if mErr != nil {
		stats.ErrListBucketsCount.Incr()
		mlog.Error("ListBuckets(): failed due to %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	output := []model.Bucket{}
	for _, bucket := range buckets {
		name, err := getRealBucketName(bucket.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
			return
		}
		bucket.Name = name
		output = append(output, bucket)
	}

	mlog.Debug("ListBuckets ends")
	c.JSON(http.StatusOK, output)
}

func (controller *Controller) DeleteBucket(c *gin.Context) {
	mlog.Debug("DeleteBucket starts")
	stats.DeleteBucketReqCount.Incr()

	accountId := c.Params.ByName("accountId")
	bucketName := c.Params.ByName("bucketName")

	if accountId == "" {
		stats.ErrListFilesMissingAccountIdCount.Incr()
		mlog.Error(util.AccountIdMissing)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}

	if bucketName == "" {
		stats.ErrDeleteBucketMissingBucketNameCount.Incr()
		errM := fmt.Sprint(util.BucketNameMissing)
		mlog.Error(errM)
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(errM))
		return
	}

	bucketName = getBucketName(accountId, bucketName)
	mErr := controller.StorageProvider.DeleteBucket(bucketName)
	if mErr != nil {
		stats.ErrDeleteBucketsCount.Incr()
		mlog.Error("DeleteBucket(): failed due to %d, %v", mErr.HttpStatus, mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		return
	}

	mlog.Debug("DeleteBucket ends")
	c.JSON(http.StatusOK, "Data received")
}

func getBucketName(accountId, bucketName string) string {
	return fmt.Sprintf("%s_%s", accountId, bucketName)
}

func getRealBucketName(bucketName string) (string, error) {
	name := strings.SplitN(bucketName, "_", 2)
	if len(name) < 2 {
		err := fmt.Sprint("Failed to split name, wrong size")
		return "", errors.New(err)
	}

	return name[1], nil
}
