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

package config

import (
	"os"
	"github.com/verizonlabs/northstar/pkg/config"
	rteCfg "github.com/verizonlabs/northstar/pkg/rte/config"
	"github.com/verizonlabs/northstar/pkg/rte/topics"
)

var (
	RTEBotcAppName        = os.Getenv("MARATHON_APP_ID")
	RTEServiceName, _     = config.GetString("RTE_SERVICE_NAME", rteCfg.RTE_SERVICE_NAME)
	RTELuaCtrlTopic, _    = config.GetString("RTE_LUA_CTRL_TOPIC", topics.RTE_LUA_CTRL_TOPIC)
	RTELuaMarathonJson, _ = config.GetString("RTE_LUA_MARATHON_JSON", "")
)
