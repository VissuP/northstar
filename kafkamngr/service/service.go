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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/kafkamngr/cluster"
	"github.com/verizonlabs/northstar/kafkamngr/model"
	"github.com/verizonlabs/northstar/kafkamngr/util"
)

type KafkaMngrService struct{}

var (
	kafkaCluster cluster.KafkaCluster
)

func (s *KafkaMngrService) SetKafkaCluster(cluster cluster.KafkaCluster) {
	kafkaCluster = cluster
}

func (s *KafkaMngrService) AddRoutes() {
	grp := management.Engine().Group(util.KafkaMngrBasePath)
	g := grp.Group("topics")
	g.POST("", createTopic)
	g.GET("", getTopicNames)
	g.PUT("/:serviceName/:topicName", updateTopic)
}

func createTopic(c *gin.Context) {
	var topic = new(model.Topic)
	c.Bind(topic)

	err := topic.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, ValidationFailed, err.Error()))
		ErrCreateTopic.Incr()
		return
	}

	exists, err := kafkaCluster.HasTopic(topic.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, HasTopicFailed, err.Error()))
		ErrCreateTopic.Incr()
		return
	}

	if exists {
		mErr := fmt.Errorf("Topic %s already exists", topic.Name)
		mlog.Error("%v", mErr)
		c.JSON(http.StatusConflict,
			management.NewError(http.StatusConflict, TopicExists, mErr.Error()))
		ErrCreateTopic.Incr()
		return
	}

	mlog.Info("Creating topic %s", topic.Name)
	err = kafkaCluster.CreateTopic(topic.Name, topic.Partitions, topic.Replication)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, CreateTopicFailed, err.Error()))
		ErrCreateTopic.Incr()
		return
	}

	mlog.Info("Topic %s created", topic.Name)
	CreateTopic.Incr()
	c.String(http.StatusCreated, "")
}

func getTopicNames(c *gin.Context) {
	mlog.Info("Retrieving topics")
	names, err := kafkaCluster.GetTopicNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, GetTopicsFailed, err.Error()))
		ErrGetTopicNames.Incr()
		return
	}

	mlog.Info("Topic names: %v", names)
	GetTopicNames.Incr()
	c.JSON(http.StatusOK, names)
}

func updateTopic(c *gin.Context) {
	serviceName := c.Params.ByName("serviceName")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("serviceName is empty"))
		ErrUpdateTopic.Incr()
		return
	}

	topicName := c.Params.ByName("topicName")
	if topicName == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("topicName is empty"))
		ErrUpdateTopic.Incr()
		return
	}

	var topic = new(model.Topic)
	c.Bind(topic)

	mlog.Info("Updating service name: %s, topic: %s, with data: %v", serviceName, topicName, topic)

	if topic.Partitions > 0 {
		err := kafkaCluster.UpdateNumberOfPartitions(serviceName,
			topicName,
			topic.Partitions)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				management.NewError(http.StatusInternalServerError, UpdateTopicFailed, err.Error()))
			ErrUpdateTopic.Incr()
			return
		}
	}

	mlog.Info("Topic %s with service name %s updated successfully", topicName, serviceName)
	UpdateTopic.Incr()
	c.String(http.StatusOK, "")
}
