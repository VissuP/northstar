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

package notebooks

import (
	"sync"

	"github.com/gocql/gocql"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/util"
)

const (
	// Defines keyspace and table names.
	Keyspace      = "account"
	NotebookTable = "notebooks"
	AccessTable   = "access"

	// Defines the max limit to be used for queries.
	MaxQueryLimit = int(100)
)

var (
	lock    sync.Mutex
	session *gocql.Session
)

// Defines the notebook data service.
type NotebookService struct {
	db *util.DB
}

// Registers service routes.
func (service *NotebookService) AddRoutes() {
	// Register service endpoints.
	v1 := management.Engine().Group(util.DataBasePath)
	{
		// Register Access Endpoints
		v1.POST("/access", service.createAccess)
		v1.PUT("/access/:accessId", service.updateAccess)
		v1.POST("/access/actions/query", service.queryAccess)
		v1.DELETE("/access/:accessId", service.deleteAccess)

		// Register Notebook Endpoints
		v1.POST("/notebooks", service.createNotebook)
		v1.PUT("/notebooks/:notebookId", service.updateNotebook)
		v1.GET("/notebooks/:notebookId", service.getNotebook)
		v1.DELETE("/notebooks/:notebookId", service.deleteNotebook)
	}
}

// Helper method used to get/create database session.
func getSession() (*gocql.Session, error) {
	mlog.Debug("getSession")
	var err error

	if session == nil || session.Closed() {
		lock.Lock()
		defer lock.Unlock()

		if session == nil || session.Closed() {
			session, err = util.NewDB(Keyspace).GetSessionWithError()
		}
	}

	return session, err
}
