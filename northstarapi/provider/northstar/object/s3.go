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

package object

import (
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
	objectClient "github.com/verizonlabs/northstar/object/client"
	objectModel "github.com/verizonlabs/northstar/object/model"
)

// S3ObjectProvider defines the type used to support operations on NorthStar S3 Objects.
type S3ObjectProvider struct {
	objectClient *objectClient.ObjectClient
}

//NewS3ObjectProvider creates a new instance of the NorthstarS3BucketProvider
func NewS3ObjectProvider() (*S3ObjectProvider, error) {
	mlog.Info("NewNorthstarS3ObjectProvider")

	client, err := objectClient.NewObjectClient()
	if err != nil {
		return nil, err
	}

	//Create the provider
	return &S3ObjectProvider{
		objectClient: client,
	}, nil
}

//GetBuckets returns the list of buckets belonging to the user.
func (provider *S3ObjectProvider) GetBuckets(user *model.User) ([]model.Bucket, *management.Error) {
	mlog.Info("GetBuckets")

	externalBuckets, err := provider.objectClient.ListBuckets(user.AccountId)
	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("list buckets returned error: %s", err.Error()))
	}

	var buckets []model.Bucket
	for _, bucket := range externalBuckets {
		buckets = append(buckets, *fromExternalBucket(&bucket))
	}
	return buckets, nil
}

//ListObjects returns the list of object in the bucket
func (provider *S3ObjectProvider) ListObjects(user *model.User, bucket string, prefix string) ([]model.Object, *management.Error) {
	mlog.Info("ListObjects")

	externalObjects, err := provider.objectClient.ListFiles(user.AccountId, bucket)
	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("list objects returned error: %s", err.Error()))
	}

	var objects []model.Object
	for _, object := range externalObjects {
		objects = append(objects, *fromExternalObject(&object))
	}

	return objects, nil

}

//GetObject returns the specified object
func (provider *S3ObjectProvider) GetObject(user *model.User, bucket string, object string) (*model.Data, *management.Error) {
	mlog.Info("GetObject")

	data, err := provider.objectClient.DownloadFile(user.AccountId, bucket, object)
	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("get object returned error: %s", err.Error()))
	}

	return &model.Data{
		ContentType: data.ContentType,
		Payload:     data.Payload,
	}, nil
}

//fromExternalBucket translates from an object bucket to a northstarapi bucket
func fromExternalBucket(externalBucket *objectModel.Bucket) *model.Bucket {
	mlog.Info("fromExternalBucket")

	bucket := &model.Bucket{
		Name:         externalBucket.Name,
		CreationDate: externalBucket.CreationDate,
	}
	return bucket
}

//toExternalBucket translates from a northstarapi bucket to an object bucket
func toExternalBucket(bucket *model.Bucket) *objectModel.Bucket {
	mlog.Info("toExternalBucket")

	externalBucket := &objectModel.Bucket{
		Name:         bucket.Name,
		CreationDate: bucket.CreationDate,
	}

	return externalBucket
}

//fromExternalObject translates from an object object to a northstarapi object
func fromExternalObject(externalObject *objectModel.Object) *model.Object {
	mlog.Info("fromExternalObject")

	object := &model.Object{
		Key:          externalObject.Key,
		LastModified: externalObject.LastModified,
		Size:         externalObject.Size,
		Etag:         externalObject.Etag,
	}

	return object
}
