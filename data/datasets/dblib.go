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

package datasets

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
	"github.com/verizonlabs/northstar/data/datasets/model"
	"github.com/verizonlabs/northstar/data/util"
)

var (
	columns = "id, accountid, datasourceid, name, description, tables, createdon, updatedon"
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

type DatasetsService struct{}

func (s *DatasetsService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("datasets")
	g.POST(":accountId", addDataset)
	g.GET(":accountId", getDatasets)
	g.GET(":accountId/by-id/:datasetId", getDatasetById)
	g.GET(":accountId/by-name/:name", getDatasetByName)
	g.PUT(":accountId/:datasetId", updateDataset)
	g.DELETE(":accountId/:datasetId", deleteDataset)
}

func addDataset(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var dataset = new(model.DatasetData)
	if err := c.Bind(dataset); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertDataset.Incr()
		return
	}

	err := dataset.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrInsertDataset.Incr()
		return
	}

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertDataset.Incr()
		return
	}

	id := uuid.NewV4().String()
	err = session.Query(`INSERT INTO `+Datasets+`(`+columns+`) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
		id, accountId, dataset.DatasourceId, dataset.Name, dataset.Description, dataset.Tables, time.Now().In(time.UTC), nil).Exec()
	if err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway, management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError, management.GetInternalError(em))
		}
		ErrInsertDataset.Incr()
		return
	}

	InsertDataset.Incr()
	c.JSON(http.StatusCreated, id)
}

func getDatasetByName(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get dataset by name due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetDatasetByName.Incr()
		return
	}

	name := c.Params.ByName("name")
	mlog.Info("Retrieving dataset %s for account %s", name, accountId)

	dataset, err := getDatasetByNameQuery(accountId, name)
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
		ErrGetDatasetByName.Incr()
		return
	}

	mlog.Info("Dataset %s retrieved for account %s", name, accountId)
	GetDatasetByName.Incr()
	c.JSON(http.StatusOK, dataset)
}

func getDatasetByNameQuery(accountId string, name string) (*model.DatasetData, error) {
	var dataset model.DatasetData

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, Datasets).
		Value("id", &dataset.Id).
		Value("accountid", &dataset.AccountId).
		Value("datasourceid", &dataset.DatasourceId).
		Value("name", &dataset.Name).
		Value("description", &dataset.Description).
		Value("tables", &dataset.Tables).
		Value("createdon", &dataset.CreatedOn).
		Value("updatedon", &dataset.UpdatedOn).
		Where("accountid", accountId).
		Where("name", name).
		Scan(session); err != nil {
		return nil, err
	}

	return &dataset, nil
}

func getDatasetById(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get dataset by Id due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetDatasetById.Incr()
		return
	}

	datasetId := c.Params.ByName("datasetId")
	mlog.Info("Retrieving dataset %s for account %s", datasetId, accountId)

	dataset, err := getDatasetByIdQuery(accountId, datasetId)
	if err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway, management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError, management.GetInternalError(em))
		}
		ErrGetDatasetById.Incr()
		return
	}

	mlog.Info("Dataset %s retrieved for account %s", dataset.Id, accountId)
	GetDatasetById.Incr()
	c.JSON(http.StatusOK, dataset)
}

func getDatasetByIdQuery(accountId string, datasetId string) (*model.DatasetData, error) {
	var dataset model.DatasetData

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, Datasets).
		Value("id", &dataset.Id).
		Value("accountid", &dataset.AccountId).
		Value("datasourceid", &dataset.DatasourceId).
		Value("name", &dataset.Name).
		Value("description", &dataset.Description).
		Value("tables", &dataset.Tables).
		Value("createdon", &dataset.CreatedOn).
		Value("updatedon", &dataset.UpdatedOn).
		Where("accountid", accountId).
		Where("id", datasetId).
		Scan(session); err != nil {
		return nil, err
	}

	return &dataset, nil
}

func getDatasets(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	datasets, err := getDatasetsQuery(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrGetDatasets.Incr()
		return
	}
	GetDatasets.Incr()
	c.JSON(http.StatusOK, datasets)
}

func getDatasetsQuery(accountId string) ([]model.DatasetData, error) {
	mlog.Info("Retrieving datasets for account %s", accountId)

	results := make([]model.DatasetData, 0, 10)
	entry := new(model.DatasetData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	iter := session.Query(`SELECT `+columns+` FROM `+Datasets+` WHERE accountid=?`, accountId).Iter()
	for iter.Scan(&entry.Id,
		&entry.AccountId,
		&entry.DatasourceId,
		&entry.Name,
		&entry.Description,
		&entry.Tables,
		&entry.CreatedOn,
		&entry.UpdatedOn) {
		results = append(results, *entry)
		entry = new(model.DatasetData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func updateDataset(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to update dataset due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateDataset.Incr()
		return
	}

	dpeId := c.Params.ByName("datasetId")

	var update = new(model.DatasetData)
	if err := c.Bind(update); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateDataset.Incr()
		return
	}

	err := updateDatasetQuery(accountId, dpeId, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError(err.Error()))
		ErrUpdateDataset.Incr()
		return
	}

	UpdateDataset.Incr()
	c.String(http.StatusOK, "")
}

func updateDatasetQuery(accountId string, datasetId string, update *model.DatasetData) error {
	queryBuilder := database.Update(Keyspace, Datasets).
		Param("updatedon", time.Now().In(time.UTC)).
		Where("accountid", accountId).
		Where("id", datasetId)

	if update.Name != "" {
		queryBuilder = queryBuilder.Param("name", update.Name)
	}

	if update.Description != "" {
		queryBuilder = queryBuilder.Param("description", update.Description)
	}

	if update.DatasourceId != "" {
		queryBuilder = queryBuilder.Param("datasourceid", update.DatasourceId)
	}

	session, err := getSession()
	if err != nil {
		return err
	}

	_, err = queryBuilder.Exec(session)
	return err
}

func deleteDataset(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete dataset due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelDataset.Incr()
		return
	}

	datasetId := c.Params.ByName("datasetId")

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelDataset.Incr()
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, Datasets).
		Where("accountId", accountId).
		Where("id", datasetId).
		Exec(session); err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError, management.GetInternalError(em))
		}
		ErrDelDataset.Incr()
		return
	}

	if !success {
		ErrDelDataset.Incr()
		mlog.Error("Data set %s not found in account %s", datasetId, accountId)
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Data set not found."))
		return
	}

	DelDataset.Incr()
	c.String(http.StatusOK, "")
}
