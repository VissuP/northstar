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

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/object/controller"
	"github.com/verizonlabs/northstar/object/env"
	"github.com/verizonlabs/northstar/object/s3"
	"github.com/verizonlabs/northstar/object/util"
)

type Service struct {
	controller *controller.Controller
	engine     *gin.Engine
}

func NewService() (service *Service, err error) {
	mlog.Debug("NewService")

	input := &s3.ProviderInput{
		Host:      env.BlobStorageHostAndPort,
		UserId:    env.BlobStorageUserId,
		Secret:    env.BlobStorageUserSecret,
		DebugFlag: env.DebugFlag,
	}

	storageProvider, err := s3.NewS3StorageProvider(input)
	if err != nil {
		mlog.Error("Error, failed to create storage provider with error %s.\n", err.Error())
		return nil, err
	}
	controller := controller.NewController(storageProvider)

	engine := management.Engine()
	g := engine.Group(util.ObjectBasePath)

	// Buckets
	g.POST("/buckets/:accountId", controller.CreateBucket)
	g.GET("/buckets/:accountId", controller.ListBuckets)
	g.DELETE("/buckets/:accountId/:bucketName", controller.DeleteBucket)

	// File
	g.POST("/files/:accountId/:bucketName", controller.UploadFile)
	g.GET("/files/:accountId/:bucketName", controller.ListFiles)
	g.GET("/files/:accountId/:bucketName/*fileName", controller.DownloadFile)
	g.DELETE("/files/:accountId/:bucketName/*fileName", controller.DeleteFile)

	service = &Service{
		controller: controller,
		engine:     engine,
	}
	return service, nil
}

func (service *Service) Start() error {
	port := ":" + env.WebPort
	return management.Listen(port)
}
