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

package cluster

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/msgq"
	"github.com/verizonlabs/northstar/pkg/zklib"
	"strconv"
)

var (
	ConsumerPartitionsExcludeListBLockZkPath = "/msgq/%s/consumer/partitions/block"
	ProducerPartitionsExcludeListBLockZkPath = "/msgq/%s/producer/partitions/block"
	ConsumerPartitionsExcludeListWLockZkPath = "/msgq/%s/consumer/partitions/wlock"
	ProducerPartitionsExcludeListWLockZkPath = "/msgq/%s/producer/partitions/wlock"
)

const (
	KAFKA_MNGR_SERVICE_NAME = "kafkamngr"
)

type MsgQAdmin struct {
	Queue *msgq.MsgQ
}

func NewMsgQAdmin() (*MsgQAdmin, error) {
	msgQ, err := msgq.NewMsgQ(KAFKA_MNGR_SERVICE_NAME, nil, nil)
	if err != nil {
		mlog.Error("Error, failed to create a MsgQ object, %s.", err.Error())
		return nil, err
	}
	return &MsgQAdmin{Queue: msgQ.(*msgq.MsgQ)}, nil
}

func (m *MsgQAdmin) markSingleWLPart(choosen int32,
	serviceName string,
	topic string,
	prod bool) bool {
	excludeListZkPath := ""
	var err error
	lockPath := ""
	var lock *zk.Lock = nil
	if prod {
		lockPath = fmt.Sprintf(ProducerPartitionsExcludeListWLockZkPath, topic)
		excludeListZkPath = fmt.Sprintf(msgq.ProducerPartitionsExcludeListZkPath, topic)
	} else {
		lockPath = fmt.Sprintf(ConsumerPartitionsExcludeListWLockZkPath, topic)
		excludeListZkPath = fmt.Sprintf(msgq.ConsumerPartitionsExcludeListZkPath, serviceName, topic)
	}
	acls := zk.WorldACL(zk.PermAll)
	lock = zk.NewLock(m.Queue.Zk.Conn, lockPath, acls)
	err = lock.Lock()
	if err != nil {
		mlog.Error("Failed to acquire zookeeper lock %s for %s", lockPath, excludeListZkPath)
		return false
	}
	defer lock.Unlock()
	excludeListPartitions, err := m.Queue.Zk.GetChildrenWithCreate(excludeListZkPath)
	if err != nil {
		return false
	}

	if len(excludeListPartitions) != 0 {
		zNode := ""
		if choosen == int32(-1) {
			zNode = fmt.Sprintf("%s/%s", excludeListZkPath, excludeListPartitions[0])
		} else {
			zNode = fmt.Sprintf("%s/%d", excludeListZkPath, choosen)
		}
		//remove black list partitions based on # of partitions
		mlog.Debug("removing zNode=%s", zNode)
		err = m.Queue.Zk.DeleteAll(zNode)
		if err != nil {
			mlog.Error("Unable to whitelist partition=%s due to %v", zNode, err)
			return false
		}
		mlog.Debug("Whitelist partition=%s succeeded", zNode)
	}
	return true
}

func (m *MsgQAdmin) getPartitionTobeMarkedBL(excludeListZkPath, topic string) (int32, error) {
	//Wait to get an available partition
	mlog.Debug("getPartitionTobeMarkedBL")
	partno := int32(-1)
	var err error = nil
	for {
		partno, err = m.identifyPartitionTobeMarkedBL(excludeListZkPath, topic)
		if err == nil {
			break
		}
		return -1, err
	}
	mlog.Debug("getPartitionTobeMarkedBL=%d", partno)

	return partno, nil
}

func (m *MsgQAdmin) identifyPartitionTobeMarkedBL(excludeListZkPath, topic string) (int32, error) {

	mlog.Debug("identifyPartitionTobeMarkedBL excludeListZkPath=%s t=%s", excludeListZkPath, topic)
	kafkaPartitions, err := m.Queue.Client.Partitions(topic)
	if err != nil {
		return -1, err
	}
	choosen := int32(-1)
	excludeListPartitions, err := m.Queue.Zk.GetChildrenWithCreate(excludeListZkPath)
	if err == zklib.ErrNoNode {
		choosen = kafkaPartitions[0]
	} else if err != nil {
		mlog.Error("Failed to get partitions from Zookeeper %v", err)
		return -1, err
	} else {
		var found bool
		//Look for an element in kafkaPartitions slice
		//but not in zkPartitions slice
		for i := range kafkaPartitions {
			//Check if the kafka partition is found in any zk partition
			found = false
			//Check if its part of blacklist partitions
			for k := range excludeListPartitions {
				intVal, err := strconv.Atoi(excludeListPartitions[k])
				if err != nil {
					mlog.Error("Failed to convert blacklist partition value %s to int %v",
						excludeListPartitions[k], err)
					return -1, err
				}
				if kafkaPartitions[i] == int32(intVal) {
					mlog.Debug("Partition %d is blacklisted. Trying another partition", intVal)
					found = true
					break
				}
			}
			if !found {
				choosen = kafkaPartitions[i]
				break
			}
		}
	}
	if choosen < 0 {
		return -1, msgq.ErrAllPartitionUsed
	}
	return choosen, nil
}

func (m *MsgQAdmin) isSingleBLPart(choosen int32,
	serviceName string,
	topic string,
	prod bool) (bool, error) {
	// identify partition from kafkaPartitions and mark in excludeListPartitions
	// mark the partition black list
	mlog.Debug("isSingleBLPart choosen=%d topic=%s ispriducer=%t", choosen, topic, prod)
	excludeListZkPath := ""
	lockPath := ""
	var err error
	if prod {
		lockPath = fmt.Sprintf(ProducerPartitionsExcludeListBLockZkPath, topic)
		excludeListZkPath = fmt.Sprintf(msgq.ProducerPartitionsExcludeListZkPath, topic)
	} else {
		lockPath = fmt.Sprintf(ConsumerPartitionsExcludeListBLockZkPath, topic)
		excludeListZkPath = fmt.Sprintf(msgq.ConsumerPartitionsExcludeListZkPath, serviceName, topic)
	}

	var lock *zk.Lock = nil
	acls := zk.WorldACL(zk.PermAll)
	lock = zk.NewLock(m.Queue.Zk.Conn, lockPath, acls)
	err = lock.Lock()
	if err != nil {
		mlog.Error("Failed to acquire zookeeper lock %s for %s", lockPath, excludeListZkPath)
		return false, err
	}
	defer lock.Unlock()
	if choosen == -1 {
		choosen, err = m.getPartitionTobeMarkedBL(excludeListZkPath, topic)
		if err != nil {
			return false, err
		}
	}
	zNode := fmt.Sprintf("%s/%d", excludeListZkPath, choosen)
	found, err := m.Queue.Zk.Exists(zNode)
	if err != nil {
		return found, err
	}
	mlog.Debug("isSingleBLPart found=%t", found)
	return found, nil
}

func (m *MsgQAdmin) markSingleBLPart(choosen int32, serviceName string, topic string, prod bool) error {
	mlog.Debug("markSingleBLPart choosen=%d topic=%s ispriducer=%t", choosen, topic, prod)

	// identify partition from kafkaPartitions and mark in excludeListPartitions
	// mark the partition black list
	excludeListZkPath := ""
	lockPath := ""
	var lock *zk.Lock = nil
	var err error
	if prod {
		lockPath = fmt.Sprintf(ProducerPartitionsExcludeListBLockZkPath, topic)
		excludeListZkPath = fmt.Sprintf(msgq.ProducerPartitionsExcludeListZkPath, topic)
	} else {
		lockPath = fmt.Sprintf(ConsumerPartitionsExcludeListBLockZkPath, topic)
		excludeListZkPath = fmt.Sprintf(msgq.ConsumerPartitionsExcludeListZkPath, serviceName, topic)
	}
	acls := zk.WorldACL(zk.PermAll)
	lock = zk.NewLock(m.Queue.Zk.Conn, lockPath, acls)
	err = lock.Lock()
	if err != nil {
		mlog.Error("Failed to acquire zookeeper lock %s for %s", lockPath, excludeListZkPath)
		return err
	}
	defer lock.Unlock()
	if choosen == -1 {
		choosen, err = m.getPartitionTobeMarkedBL(excludeListZkPath, topic)
		if err != nil {
			return err
		}
	}
	zkpath := fmt.Sprintf("%s/%d", excludeListZkPath, choosen)
	exists, err := m.Queue.Zk.Exists(zkpath)
	if exists {
		mlog.Debug("Zookeeper node %s is already marked BL", zkpath, err)
		return nil

	}
	err = m.Queue.Zk.Create(zkpath, []byte(""), false)
	if err != nil {
		mlog.Error("Unable to create znode %s due to %v", zkpath, err)
		return err
	}
	return nil
}
