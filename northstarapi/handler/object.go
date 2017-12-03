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

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

//ListBuckets retrieves a list of buckets belonging to the user
func (controller *Controller) ListBuckets(context *gin.Context) {
	mlog.Info("ListBuckets")

	// Get user information.
	user, mErr := controller.getUser(context)
	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrListBuckets.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	buckets, err := controller.objectProvider.GetBuckets(user)
	if err != nil {
		mlog.Error("Failed to list buckets with error: %s", mErr.Description)
		utils.ErrListBuckets.Incr()
		controller.RenderServiceError(context, mErr)
	}

	context.JSON(http.StatusOK, buckets)
}

//GetObject returns the specified object
func (controller *Controller) GetObject(context *gin.Context) {
	mlog.Info("ListObjects")
	bucket := context.Params.ByName("bucket")
	path := context.Params.ByName("path")

	// Get user information.
	user, mErr := controller.getUser(context)
	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrGetObject.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	data, mErr := controller.objectProvider.GetObject(user, bucket, path)
	mlog.Info("DATA:%s contentType: %s", data.Payload, data.ContentType)
	if mErr != nil {
		mlog.Error("Failed to get object with error: %s", mErr.Description)
		controller.RenderServiceError(context, mErr)
		utils.ErrGetObject.Incr()
		return

	}
	context.Data(http.StatusOK, data.ContentType, data.Payload)
	utils.GetObject.Incr()
}

//ListObjects retrieves a list of files belonging to the user
func (controller *Controller) ListObjects(context *gin.Context) {
	mlog.Info("GetObject")

	mlog.Info("ListObjects")
	bucket := context.Params.ByName("bucket")
	path := context.Params.ByName("path")
	//TODO:Extend ListObjects to handle pagination.
	//TOOD: Extend listObjects to support s3 folders.

	// Get user information.
	user, mErr := controller.getUser(context)
	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrListObjects.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	objects, mErr := controller.objectProvider.ListObjects(user, bucket, path)
	if mErr != nil {
		mlog.Error("Failed to list objects with error: %s", mErr.Description)
		controller.RenderServiceError(context, mErr)
		utils.ErrListObjects.Incr()
		return
	}

	utils.ListObjects.Incr()
	context.JSON(http.StatusOK, objects)
}
