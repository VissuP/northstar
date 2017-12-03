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

package s3

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"

	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/object/model"
)

// Define the parameters used to sign the request.
var ParametersToSign = map[string]bool{
	"acl":                          true,
	"location":                     true,
	"logging":                      true,
	"notification":                 true,
	"partNumber":                   true,
	"policy":                       true,
	"requestPayment":               true,
	"torrent":                      true,
	"uploadId":                     true,
	"uploads":                      true,
	"versionId":                    true,
	"versioning":                   true,
	"versions":                     true,
	"response-content-type":        true,
	"response-content-language":    true,
	"response-expires":             true,
	"response-cache-control":       true,
	"response-content-disposition": true,
	"response-content-encoding":    true,
	"website":                      true,
	"delete":                       true,
}

// Defines the storage provider interface
type S3StorageProvider struct {
	S3Storage  *s3.S3
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

// Creates a new S3 storage provider.
func NewS3StorageProvider(input *ProviderInput) (StorageProvider, error) {
	mlog.Info("NewS3StorageProvider")

	if input == nil {
		mlog.Alarm("Storage provider can not be created since input data is nil")
		return nil, management.GetInternalError(
			fmt.Sprintf("Error, failed to create storage provider due to empty input"))
	}

	// Modify the default configuration
	aws.DefaultConfig.Endpoint = input.Host
	aws.DefaultConfig.Region = "vzlabs"
	aws.DefaultConfig.Credentials = credentials.NewStaticCredentials(input.UserId, input.Secret, "")
	aws.DefaultConfig.DisableSSL = true
	aws.DefaultConfig.DisableParamValidation = false
	aws.DefaultConfig.S3ForcePathStyle = true

	// Debugging
	if input.DebugFlag {
		aws.DefaultConfig.LogLevel = uint(1)
		aws.DefaultConfig.LogHTTPBody = true
	}

	// Create the S3 client using default configuration.
	awsS3StorageProvider := &S3StorageProvider{
		S3Storage:  s3.New(nil),
		downloader: s3manager.NewDownloader(nil),
		uploader:   s3manager.NewUploader(nil),
	}

	// Setup the handler used to sign the request before a send. Note that
	// by default AWS GO SDK only supports S3 Sign Version 4 but EMC ViPr
	// uses Version 2.
	awsS3StorageProvider.S3Storage.Handlers.Sign.Clear()
	awsS3StorageProvider.S3Storage.Handlers.Sign.PushBack(aws.BuildContentLength)
	awsS3StorageProvider.S3Storage.Handlers.Sign.PushBack(awsS3StorageProvider.signVersion2)
	return awsS3StorageProvider, nil
}

// Stores item under the specified path.
func (S3StorageProvider *S3StorageProvider) Store(bucketName, path string, data string) *management.Error {
	mlog.Info("Ready to Store - path: %s data=%s", path, data)
	// Set the S3 item (full) name
	key := path

	// Create the S3 put object request
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		ACL:    aws.String("private"),
		Body:   strings.NewReader(data),
	}

	_, err := S3StorageProvider.S3Storage.PutObject(params)
	if err != nil {
		// If external error, alert external service error.
		mlog.Alarm("External Error, S3 Service Failed with error: %s", err.Error())
		return management.GetExternalError(
			fmt.Sprintf("Error, failed to put item in remote storage with error: %s", err.Error()))
	}

	mlog.Info("Successfully uploaded with path: %s data=%s", path, data)
	return nil
}

// Deletes the data for the specified file name.
func (S3StorageProvider *S3StorageProvider) Delete(bucketName, fileName string) *management.Error {
	mlog.Debug("Ready to Delete - fileName: %s", fileName)
	key := fileName

	params := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: []*s3.ObjectIdentifier{
				{
					Key: aws.String(key),
				},
			},
			Quiet: aws.Boolean(true),
		},
	}
	dobjs, err := S3StorageProvider.S3Storage.DeleteObjects(params)
	if err != nil {
		return management.GetInternalError(
			fmt.Sprintf("Error, failed to delete file %s with error: %s", fileName, err.Error()))
	}

	mlog.Debug("Deleted - %s", dobjs.String())
	return nil
}

// Helper method used to sign the S3 request using version 2.
func (S3StorageProvider *S3StorageProvider) signVersion2(request *aws.Request) {
	// If the request does not need to be signed ignore the signing of the
	// request if the AnonymousCredentials object is used.
	creds := aws.DefaultConfig.Credentials
	if request.Service.Config.Credentials == credentials.AnonymousCredentials {
		return
	}

	if err := request.Build(); err != nil {
		mlog.Error("Error, failed to build request with error: %s", err.Error())
		return
	}

	method := request.Operation.HTTPMethod
	if method == "" {
		method = "GET"
	}

	baseUrl := request.HTTPRequest.URL.String()
	if baseUrl == "" {
		// This should not happen. If it does, we should return error.
		baseUrl = aws.DefaultConfig.Endpoint
	}

	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		mlog.Error("Error, failed to parse url with error: %s", err.Error())
		return
	}

	signpath := parsedUrl.Path
	if strings.HasPrefix(signpath, "/") == false {
		signpath = "/" + signpath
	}

	signpath = (&url.URL{Path: signpath}).String()
	request.HTTPRequest.Header["Host"] = []string{parsedUrl.Host}
	request.HTTPRequest.Header["Date"] = []string{time.Now().In(time.UTC).Format(time.RFC1123)}

	value, _ := creds.Get()

	if value.SessionToken != "" {
		request.HTTPRequest.Header["X-Amz-Security-Token"] = []string{value.SessionToken}
	}

	var md5, ctype, date, xamz string
	var xamzDate bool
	var sarray []string
	for k, v := range request.HTTPRequest.Header {
		k = strings.ToLower(k)
		switch k {
		case "content-md5":
			md5 = v[0]
		case "content-type":
			ctype = v[0]
		case "date":
			if !xamzDate {
				date = v[0]
			}
		default:
			if strings.HasPrefix(k, "x-amz-") {
				vall := strings.Join(v, ",")
				sarray = append(sarray, k+":"+vall)
				if k == "x-amz-date" {
					xamzDate = true
					date = ""
				}
			}
		}
	}
	if len(sarray) > 0 {
		sort.StringSlice(sarray).Sort()
		xamz = strings.Join(sarray, "\n") + "\n"
	}

	rawQuery := strings.Split(parsedUrl.RawQuery, "&")
	sarray = sarray[0:0]
	for _, param := range rawQuery {
		query := strings.Split(param, "=")
		key := query[0]
		value := ""

		if len(query) == 2 {
			value = query[1]
		}

		if ParametersToSign[key] {
			if value == "" {
				sarray = append(sarray, key)
			} else {
				sarray = append(sarray, key+"="+value)
			}
		}
	}
	if len(sarray) > 0 {
		sort.StringSlice(sarray).Sort()
		signpath = signpath + "?" + strings.Join(sarray, "&")
	}

	payload := method + "\n" + md5 + "\n" + ctype + "\n" + date + "\n" + xamz + signpath
	hash := hmac.New(sha1.New, []byte(value.SecretAccessKey))
	hash.Write([]byte(payload))
	signature := make([]byte, base64.StdEncoding.EncodedLen(hash.Size()))
	base64.StdEncoding.Encode(signature, hash.Sum(nil))

	request.HTTPRequest.Header["Authorization"] = []string{"AWS " + value.AccessKeyID + ":" + string(signature)}
}

// Returns metadata under the specified path.
func (S3StorageProvider *S3StorageProvider) List(bucketName,
	path string) ([]model.Object, string, *management.Error) {
	mlog.Debug("List - path: %s", path)

	prefix := path
	// For now, we use no delimiter.
	delimiter := ""

	mlog.Debug("Listing content of S3 bucket %s - prefix: %s delimiter: %s ",
		bucketName, prefix, delimiter)
	params := &s3.ListObjectsInput{
		Bucket:    aws.String(bucketName),
		Delimiter: aws.String(delimiter),
		Prefix:    aws.String(prefix),
	}
	response, err := S3StorageProvider.S3Storage.ListObjects(params)
	if err != nil {
		// If external error, alert external service error.
		mlog.Alarm("External Error, S3 Service failed with error: %s.", err.Error())

		return nil, "", management.GetExternalError(
			fmt.Sprintf("Error, failed to list items in remote storage with error: %s", err.Error()))
	}

	// Generate the string items
	str := make([]model.Object, len(response.Contents), len(response.Contents))
	for index, item := range response.Contents {
		if item.Key != nil {
			str[index] = model.Object{Key: *item.Key,
				Etag:         *item.ETag,
				LastModified: *item.LastModified,
				Size:         *item.Size,
				StorageClass: *item.StorageClass}
			mlog.Debug("Data is %v, str[index]=%v", item, str[index])
		}
	}
	// Set the next marker
	next := ""
	if response.NextMarker != nil {
		next = *response.NextMarker
	}

	return str, next, nil
}

// Upload given file to S3
func (S3StorageProvider *S3StorageProvider) Upload(bucketName string,
	data *model.UploadData) *management.Error {
	mlog.Info("Uploading file %s", data.FileName)

	// Create the S3 put object request
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(data.FileName),
		ACL:    aws.String("private"),
		Body:   aws.ReadSeekCloser(bytes.NewReader(data.Payload)),
	}

	if data.ContentType != "" {
		params.ContentType = aws.String(data.ContentType)
	}

	response, err := S3StorageProvider.S3Storage.PutObject(params)
	if err != nil {
		// If external error, alert external service error.
		mlog.Alarm("External Error, S3 Service Failed with error: %s", err.Error())
		return management.GetExternalError(
			fmt.Sprintf("Error, failed to put item in remote storage with error: %s", err.Error()))
	}

	mlog.Info("Successfully uploaded object %s with response %v", data.FileName, response)
	return nil
}

// Download given file to S3
func (S3StorageProvider *S3StorageProvider) Download(bucketName,
	fileName string) (*model.DownloadData, *management.Error) {

	mlog.Debug("Get - fileName: %s", fileName)
	// Get the S3 object.
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}

	response, err := S3StorageProvider.S3Storage.GetObject(params)
	if err != nil {
		return nil, management.GetExternalError(
			fmt.Sprintf("Error, failed to get file %s with error: %s", fileName, err.Error()))
	}

	mlog.Info("Response is %v", response)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		mlog.Error("Get(): ioutil.ReadAll failed due to %v", err)
		return nil, management.GetInternalError(
			fmt.Sprintf("Error, failed to read file %s with error: %s", fileName, err.Error()))
	}

	mlog.Debug("Downloaded file %s successfully", fileName)
	return &model.DownloadData{ContentType: *response.ContentType, Payload: data}, nil
}

func (S3StorageProvider *S3StorageProvider) CreateBucket(bucketName string) *management.Error {
	mlog.Debug("CreateBucket - bucketName: %s", bucketName)

	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{},
	}

	response, err := S3StorageProvider.S3Storage.CreateBucket(params)
	if err != nil {
		return management.GetExternalError(
			fmt.Sprintf("Error, failed to create bucket %s with error: %s", bucketName, err.Error()))
	}

	mlog.Debug("Created bucket %s successfully: %v", bucketName, response)
	return nil
}

func (S3StorageProvider *S3StorageProvider) DeleteBucket(bucketName string) *management.Error {
	mlog.Debug("DeleteBucket - bucketName: %s", bucketName)

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}

	response, err := S3StorageProvider.S3Storage.DeleteBucket(params)
	if err != nil {
		return management.GetExternalError(
			fmt.Sprintf("Error, failed to delete bucket %s with error: %s", bucketName, err.Error()))
	}

	mlog.Debug("Deleted bucket %s successfully: %v", bucketName, response)
	return nil
}

func (S3StorageProvider *S3StorageProvider) ListBuckets(nameFilter string) ([]model.Bucket,
	*management.Error) {
	mlog.Debug("ListBuckets start")

	params := &s3.ListBucketsInput{}
	response, err := S3StorageProvider.S3Storage.ListBuckets(params)
	if err != nil {
		return nil, management.GetExternalError(
			fmt.Sprintf("Error, failed to list buckets with error: %s", err.Error()))
	}

	mlog.Debug("Buckets retrieved successfully: %v", response)

	var buckets = []model.Bucket{}
	for _, bucket := range response.Buckets {
		var newBucket *model.Bucket
		if nameFilter == "" {
			newBucket = &model.Bucket{Name: *bucket.Name, CreationDate: *bucket.CreationDate}
		} else {
			if strings.Contains(*bucket.Name, nameFilter) {
				newBucket = &model.Bucket{Name: *bucket.Name, CreationDate: *bucket.CreationDate}
			}
		}

		if newBucket != nil {
			buckets = append(buckets, *newBucket)
		}
	}

	mlog.Debug("List buckets successfull: %v", len(buckets))
	return buckets, nil
}
