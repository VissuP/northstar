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
	"time"

	"github.com/channelmeter/vault-gatekeeper-mesos/gatekeeper"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/verizonlabs/northstar/pkg/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

var (
	// config signature changed.  the default value 2 is not passsed anymore.
	CassandraDatacenter, _               = config.GetString("DATACENTER", "")
	CassandraProtoVersion, _             = config.GetInt("CASSANDRA_PROTO_VERSION", 3)
	CassandraCQLVersion, _               = config.GetString("CASSANDRA_CQL_VERSION", "")
	CassandraAuthEnabled, _              = config.GetBool("CASSANDRA_AUTH_ENABLED", false)
	VaultHostPort, _                     = config.GetString("VAULT_HOST_PORT", "")
	GatekeeperHostPort, _                = config.GetString("GATEKEEPER_HOST_PORT", "")
	CassandraUsername, CassandraPassword = GetCassandraAuthCredentials(GatekeeperHostPort, VaultHostPort)
)

func GetCassandraAuthCredentials(gatekeeperHostPort string, vaultHostPort string) (username string, password string) {
	if CassandraAuthEnabled == true {
		vaultAddress := "http://" + vaultHostPort
		gatekeeperAddress := "http://" + gatekeeperHostPort
		vaultToken := vaultAuth(gatekeeperAddress, vaultAddress)
		for {
			client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
			if err != nil {
				mlog.Alarm("Error initializing vault client: %s", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			client.SetAddress(vaultAddress)
			client.SetToken(vaultToken)
			c := client.Logical()
			secret, err := c.Read("secret/cassandra")
			if err != nil {
				return username, password
			}
			username := secret.Data["username"].(string)
			password := secret.Data["password"].(string)
			return username, password
		}
	} else {
		return username, password
	}

}

func vaultAuth(gatekeeperAddress string, vaultAddress string) string {
	mlog.Debug("Authenticating with Vault")
	mesosTaskId := fetchAndCheckEnv("MESOS_TASK_ID", true)
	for {
		mlog.Info("Getting token from gatekeeper for mesosTaskId %s", mesosTaskId)
		client, err := gatekeeper.NewClient(vaultAddress, gatekeeperAddress, nil)
		if err != nil {
			mlog.Alarm("Error initializing gatekeeper client: %s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		token, err := client.RequestVaultToken(mesosTaskId)
		if err == nil {
			return token
		} else {
			mlog.Alarm("Couldn't get vault token: %s", err.Error())
			time.Sleep(1 * time.Second)
		}
	}
}

func fetchAndCheckEnv(env string, mandatory bool) string {
	val := os.Getenv(env)
	if mandatory && len(val) <= 0 {
		mlog.Error("%s env variable is not set", env)
	}
	return val
}
