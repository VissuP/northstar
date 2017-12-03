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

package northstar

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	northstarApiClient "github.com/verizonlabs/northstar/northstarapi/client"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/client/auth"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/utils"
	"strconv"
)

const UserCount = 2

// Defines the type used as the ThingSpace base test.
type NorthstarApiBaseTest struct {
	Users              []User
	ClientToken        *auth.Token
	NorthstarApiClient *northstarApiClient.Client
}

type User struct {
	AccountId string
	Token     *auth.Token
	Email     string
	Password  string
}

// Method used to create a new ThingSpace base test.
func NewNorthstarApiBaseTest() (*NorthstarApiBaseTest, error) {
	northstarApiClient, err := northstarApiClient.NewClient(config.Configuration.NorthstarProtocol, config.Configuration.NorthstarApiHostPort)
	if err != nil {
		return nil, err
	}

	return &NorthstarApiBaseTest{
		NorthstarApiClient: northstarApiClient,
	}, err
}

// Creates resources needed to execute ThingSpace tests.
func (northstarApiBaseTest *NorthstarApiBaseTest) Initialize(logs *utils.Logger) (err error) {
	logs.LogInfo("Initialize Base")

	// Get the token associated with this test application.
	var getClientTokenStep utils.Step = "Get client token."
	var mErr *management.Error
	northstarApiBaseTest.ClientToken, mErr = auth.GetClientToken()
	if mErr != nil {
		return logs.LogStep(getClientTokenStep, false, mErr.Description)
	}
	logs.LogStep(getClientTokenStep, true, "Success.")

	// Initialize the global account manager.
	var getAccountsStep utils.Step = "Get accounts."
	accountManager = GetAccountManager()
	accounts, err := accountManager.CreateAccounts(false, UserCount)
	if err != nil {
		return logs.LogStep(getAccountsStep, false, "Error, failed to get accounts with error: %s", err.Error())
	}
	logs.LogStep(getAccountsStep, true, "Success. Accounts: %+v", accounts)

	var getUserTokensStep utils.Step = "Get user tokens."
	for _, account := range accounts {
		user := User{
			Email:    account.Email,
			Password: account.Password,
		}
		user.Token, mErr = auth.GetUserToken(account.Email, account.Password)
		if mErr != nil {
			return logs.LogStep(getUserTokensStep, false, "Failed to get user token with error: %s", mErr.Description)
		}

		//LATER: get account here
		
		user.AccountId = "nullAccountID"  //LATER: fix me
		northstarApiBaseTest.Users = append(northstarApiBaseTest.Users, user)
	}
	logs.LogStep(getUserTokensStep, true, "Success.")

	logs.LogInfo("Base test initialized with config: %+v", northstarApiBaseTest, northstarApiBaseTest.Users[0].Token, northstarApiBaseTest.Users[1].Token)
	return nil
}

// Cleanup removes resources allocated for a northstarapi test.
func (northstarApiBaseTest *NorthstarApiBaseTest) Cleanup(logs *utils.Logger) error {
	logs.LogInfo("Cleanup")

	var deleteClientTokenStep utils.Step = "Delete client token."
	logs.LogStep(deleteClientTokenStep, true, "Success.")

	return nil
}

// GetEventCallbackURL is a helper method used to generate callback url.
func (northstarApiBaseTest *NorthstarApiBaseTest) GetEventCallbackURL(uuid string) string {
	return fmt.Sprintf("http://%s/sim/v1/callback/%s", config.Configuration.ServiceHostPort, uuid)
}

//GetUUID returns a UUID
func (northstarApiBaseTest *NorthstarApiBaseTest) GetUUID() (string, error) {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		mlog.Error("Couldn't generate a UUID")
		return "", err
	}
	return uuid.String(), nil
}

func (northstarApiBaseTest *NorthstarApiBaseTest) ExecuteCell(logs *utils.Logger, callbacks *utils.ThreadSafeMap, code string) (*json.RawMessage, error) {
	logs.LogInfo("Execute Cell")

	uuid, err := northstarApiBaseTest.GetUUID()
	if err != nil {
		logs.LogInfo("Failed to get callback ID with error: %s", err)
		return nil, err
	}

	callbackURL := northstarApiBaseTest.GetEventCallbackURL(uuid)

	writeChannel := make(chan json.RawMessage)
	callbacks.Set(uuid, writeChannel)
	defer callbacks.Delete(uuid)
	logs.LogInfo("Registered callback URL: %s", callbackURL)

	cellUuid, err := northstarApiBaseTest.GetUUID()
	if err != nil {
		logs.LogError("Failed to get cell ID with error: %s", err)
		return nil, err
	}
	cell := northstarApiModel.Cell{
		Id:   cellUuid,
		Name: "Test Cell",
		Input: northstarApiModel.Input{
			Type:       northstarApiModel.CodeCellType,
			Language:   "lua",
			EntryPoint: "main",
			Body:       code,
		},
	}

	//Base64 encode the code section for execution
	cell.Input.Body = base64.StdEncoding.EncodeToString([]byte(cell.Input.Body))

	logs.LogInfo("Executing: %+v", cell)
	serviceErr := northstarApiBaseTest.NorthstarApiClient.ExecuteCell(northstarApiBaseTest.Users[0].Token.AccessToken, callbackURL, &cell)
	if serviceErr != nil {
		logs.LogError("Error executing cell. Error was: %s", serviceErr.Description)
		return nil, fmt.Errorf(serviceErr.Description)
	}

	select {
	case response := <-writeChannel:
		return &response, nil
	case <-time.After(config.Configuration.ExecutionResponseTimeout):
		logs.LogError("Error, no response received for cell execution. Timing out.")
		return nil, fmt.Errorf("Error, no response received for cell execution. Timing out.")
	}
	return nil, fmt.Errorf("Fell through execution switch. Shouldn't have been able to do that.")
}

//DeleteNotebook deletes the specified notebook
func (northstarApiBaseTest *NorthstarApiBaseTest) DeleteNotebook(logs *utils.Logger, token string, notebookID string) error {
	logs.LogInfo("Deleting notebook: %s", notebookID)
	mErr := northstarApiBaseTest.NorthstarApiClient.DeleteNotebook(token, notebookID)
	if mErr != nil {
		logs.LogError("Error, failed to delete notebook with error: %s", mErr.Description)
		return fmt.Errorf(mErr.Description)
	}

	logs.LogInfo("Notebook %v deleted", notebookID)
	return nil
}

//DeleteTemplate deletes the specified template
func (northstarApiBaseTest *NorthstarApiBaseTest) DeleteTemplate(logs *utils.Logger, token string, templateID string) error {
	logs.LogInfo("Deleting template: %s", templateID)
	if err := northstarApiBaseTest.NorthstarApiClient.DeleteTemplate(token, templateID); err != nil {
		return logs.LogError("Error, could not delete public cell template. Error was: %s", err)
	}
	return nil
}

//DeleteTransformation deletes the specified transformation
func (northstarApiBaseTest *NorthstarApiBaseTest) DeleteTransformation(logs *utils.Logger, token string, transformationID string) error {
	logs.LogInfo("Deleting transformation: %s", transformationID)

	serviceErr := northstarApiBaseTest.NorthstarApiClient.DeleteTransformation(token, transformationID)
	if serviceErr != nil {
		return logs.LogError("Error, failed to delete transformation: %s", serviceErr.Description)
	}
	return nil
}

//DeleteSchedule deletes the specified transformation
func (northstarApiBaseTest *NorthstarApiBaseTest) DeleteSchedule(logs *utils.Logger, token string, scheduleID string) error {
	logs.LogInfo("Deleting schedule: %s", scheduleID)

	serviceErr := northstarApiBaseTest.NorthstarApiClient.DeleteSchedule(token, scheduleID)
	if serviceErr != nil {
		return logs.LogError("Error, failed to delete schedule: %s", serviceErr.Description)
	}
	return nil
}

type tableRow struct {
	Row      []interface{}
	DataType string
}

//TableToMap will take the table output by GetRow and parse the row and column values from it.
func (northstarApiBaseTest *NorthstarApiBaseTest) TableToMap(logs *utils.Logger, results *northstarApiModel.CellResults, types map[string]string) (map[string]tableRow, error) {
	logs.LogInfo("parseTable")

	if results == nil || results.Type != "application/vnd.vz.table" {
		return nil, logs.LogError("Error, incorrect results type. Can't parse table. Results: %+v", results)
	}

	content, ok := results.Content.(map[string]interface{})
	if !ok {
		return nil, logs.LogError("Cannot parse content")
	}
	rows, ok := content["rows"].([]interface{})
	if !ok {
		return nil, logs.LogError("Cannot parse rows from content")
	}

	columns, ok := content["columns"].([]interface{})
	if !ok {
		return nil, logs.LogError("Cannot parse columns from content")
	}

	dataTypes, _ := content["types"].([]interface{})

	output := make(map[string]tableRow)
	for columnIndex, columnValue := range columns {
		columnName, ok := columnValue.(string)
		if !ok {
			return nil, logs.LogError("Error, cannot parse column name")
		}
		columnType, _ := types[columnName]

		//Get the data type off of the table.
		dataType := ""
		if dataTypes != nil {
			dataTypeIntf := dataTypes[columnIndex]
			if ok {
				dataType, _ = dataTypeIntf.(string)
			}
		}

		row := make([]interface{}, 0)
		for _, rowValue := range rows {
			rowData, ok := rowValue.([]interface{})
			if !ok {
				return nil, logs.LogError("Cannot parse row value")
			}

			value := rowData[columnIndex]
			if columnType == "date" {
				parsedTime, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v", value))
				if err != nil {
					return nil, fmt.Errorf("Error, could not parse date %s with error: %s", value, err.Error())
				}
				value = parsedTime.Format("2006-01-02")
			} else if columnType == "time" {
				parsedTime, err := parseLuaTime(logs, value)
				if err != nil {
					return nil, err
				}
				value = parsedTime.Format("15:04:05")
			}

			row = append(row, value)
		}
		output[columnName] = tableRow{
			Row:      row,
			DataType: dataType,
		}
	}

	logs.LogInfo("Output: %+v", output)

	return output, nil
}

func parseLuaTime(logs *utils.Logger, value interface{}) (*time.Time, error) {
	valueString := fmt.Sprintf("%v", value)

	logs.LogInfo("Time as string: %s", valueString)
	timeInt, err := strconv.ParseInt(valueString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse time %s", valueString)
	}
	parsedTime := time.Unix(0, timeInt)
	return &parsedTime, nil
}
