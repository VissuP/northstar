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

package utils

import (
	"fmt"
	"sync"
)

//ThreadSafeMap implements a multithreading-safe map
type ThreadSafeMap struct {
	sync.RWMutex
	threadMap map[string]interface{}
}

//NewThreadSafeMap creates a new ThreadSafeMap
func NewThreadSafeMap() *ThreadSafeMap {
	return &ThreadSafeMap{
		threadMap: make(map[string]interface{}),
	}
}

//Get retrieves the map entry matching the provided key.
func (tsMap *ThreadSafeMap) Get(key string) (interface{}, error) {
	tsMap.RLock()
	value, found := tsMap.threadMap[key]
	tsMap.RUnlock()

	if !found {
		return nil, fmt.Errorf("error, value for key %s not found", key)
	}

	return value, nil
}

//Set sets the value corresponding to the provided key.
func (tsMap *ThreadSafeMap) Set(key string, value interface{}) {
	tsMap.Lock()
	tsMap.threadMap[key] = value
	tsMap.Unlock()
}

//Delete removes the corresponding key from the map.
func (tsMap *ThreadSafeMap) Delete(key string) {
	tsMap.Lock()
	delete(tsMap.threadMap, key)
	tsMap.Unlock()
}

//Size returns the number of elements in the map.
func (tsMap *ThreadSafeMap) Size() int {
	return len(tsMap.threadMap)
}
