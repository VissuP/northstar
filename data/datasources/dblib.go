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

package datasources

import (
	"net/http"

	"fmt"

	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/datasources/model"
	"github.com/verizonlabs/northstar/data/util"
)

var (
	columns = "id, accountid, name, description, protocol, host, port, options, createdon, updatedon"
	sess    *gocql.Session
	lock    sync.Mutex
)

// Helper method used to get/create database session.
func getSession() (*gocql.Session, error) {
	var err error

	if sess == nil || sess.Closed() {
		lock.Lock()
		defer lock.Unlock()

		if sess == nil || sess.Closed() {
			sess, err = util.NewDB(Keyspace).GetSessionWithError()
		}
	}

	return sess, err
}

type DatasourcesService struct{}

func (s *DatasourcesService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("datasources")
	g.POST(":accountId", addDatasource)
	g.GET(":accountId", getDatasources)
	g.GET(":accountId/:datasourceId", getDatasource)
	g.PUT(":accountId/:datasourceId", updateDatasource)
	g.DELETE(":accountId/:datasourceId", deleteDatasource)
}

func addDatasource(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var datasource = new(model.DatasourceData)
	if err := c.Bind(datasource); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertDatasource.Incr()
		return
	}

	err := datasource.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrInsertDatasource.Incr()
		return
	}

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertDatasource.Incr()
		return
	}

	id := uuid.NewV4().String()
	err = session.Query(`INSERT INTO `+Datasources+`(`+columns+`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, accountId, datasource.Name, datasource.Description, datasource.Protocol, datasource.Host,
		datasource.Port, datasource.Options, time.Now().In(time.UTC), nil).Exec()
	if err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(em))
		}
		ErrInsertDatasource.Incr()
		return
	}

	InsertDatasource.Incr()
	c.String(http.StatusCreated, id)
}

func getDatasource(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get datasource due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetDatasource.Incr()
		return
	}

	datasourceId := c.Params.ByName("datasourceId")
	mlog.Info("Retrieving dataource %s for account %s", datasourceId, accountId)

	dataource, err := getDatasourceQuery(accountId, datasourceId)
	if err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(em))
		}
		ErrGetDatasource.Incr()
		return
	}

	mlog.Info("Datasource %s retrieved for account %s", dataource.Id, accountId)
	GetDatasource.Incr()
	c.JSON(http.StatusOK, dataource)
}

func getDatasourceQuery(accountId string, datasourceId string) (*model.DatasourceData, error) {
	var dataource model.DatasourceData

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, Datasources).
		Value("id", &dataource.Id).
		Value("accountid", &dataource.AccountId).
		Value("name", &dataource.Name).
		Value("description", &dataource.Description).
		Value("options", &dataource.Options).
		Value("createdon", &dataource.CreatedOn).
		Value("updatedon", &dataource.UpdatedOn).
		Where("accountid", accountId).
		Where("id", datasourceId).
		Scan(session); err != nil {
		return nil, err
	}

	return &dataource, nil
}

func getDatasources(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	datasources, err := getDatasourcesQuery(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrGetDatasources.Incr()
		return
	}
	GetDatasources.Incr()
	c.JSON(http.StatusOK, datasources)
}

func getDatasourcesQuery(accountId string) ([]model.DatasourceData, error) {
	mlog.Info("Retrieving datasources for account %s", accountId)

	results := make([]model.DatasourceData, 0, 10)
	entry := new(model.DatasourceData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	iter := session.Query(`SELECT `+columns+` FROM `+Datasources+` WHERE accountid=?`, accountId).Iter()
	for iter.Scan(&entry.Id,
		&entry.AccountId,
		&entry.Name,
		&entry.Description,
		&entry.Protocol,
		&entry.Host,
		&entry.Port,
		&entry.Options,
		&entry.CreatedOn,
		&entry.UpdatedOn) {
		results = append(results, *entry)
		entry = new(model.DatasourceData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func updateDatasource(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to update datasource due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateDatasource.Incr()
		return
	}

	dpeId := c.Params.ByName("datasourceId")

	var update = new(model.DatasourceData)
	if err := c.Bind(update); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateDatasource.Incr()
		return
	}

	err := updateDatasourceQuery(accountId, dpeId, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError(err.Error()))
		ErrUpdateDatasource.Incr()
		return
	}

	UpdateDatasource.Incr()
	c.String(http.StatusOK, "")
}

func updateDatasourceQuery(accountId string, datasourceId string, update *model.DatasourceData) error {
	queryBuilder := database.Update(Keyspace, Datasources).
		Param("updatedon", time.Now().In(time.UTC)).
		Where("accountid", accountId).
		Where("id", datasourceId)

	if update.Name != "" {
		queryBuilder = queryBuilder.Param("name", update.Name)
	}

	if update.Description != "" {
		queryBuilder = queryBuilder.Param("description", update.Description)
	}

	if update.Protocol != "" {
		queryBuilder = queryBuilder.Param("protocol", update.Protocol)
	}

	if update.Host != "" {
		queryBuilder = queryBuilder.Param("host", update.Host)
	}

	if update.Port > 0 {
		queryBuilder = queryBuilder.Param("port", update.Port)
	}

	if len(update.Options) > 0 {
		queryBuilder = queryBuilder.Param("options", update.Options)
	}

	session, err := getSession()
	if err != nil {
		return err
	}

	_, err = queryBuilder.Exec(session)
	return err
}

func deleteDatasource(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete datasource due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelDatasource.Incr()
		return
	}

	datasourceId := c.Params.ByName("datasourceId")

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelDatasource.Incr()
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, Datasources).
		Where("accountId", accountId).
		Where("id", datasourceId).
		Exec(session); err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetInternalError(em))
		}
		ErrDelDatasource.Incr()
		return
	}

	if !success {
		ErrDelDatasource.Incr()
		mlog.Error("Data source %s not found in account %s", datasourceId, accountId)
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Data source not found."))
		return
	}

	DelDatasource.Incr()
	c.String(http.StatusOK, "")
}
