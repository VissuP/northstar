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
	"encoding/json"
	"fmt"
	"time"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	notebooks "github.com/verizonlabs/northstar/data/notebooks/client"
	notebooksModel "github.com/verizonlabs/northstar/data/notebooks/model"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

// Defines the type used to support operations on NorthStar notebooks.
type NorthStarNotebooksProvider struct {
	notebooksClient *notebooks.NotebooksClient
}

// Returns a new NorthStar notebook provider.
func NewNorthStarNotebooksProvider() (*NorthStarNotebooksProvider, error) {
	mlog.Info("NewNorthStarNotebooksProvider")

	client, err := notebooks.NewNotebooksClient()
	if err != nil {
		return nil, err
	}

	// Create the provider.
	return &NorthStarNotebooksProvider{notebooksClient: client}, nil
}

// Creates a new notebook with associated access information.
func (provider *NorthStarNotebooksProvider) Create(user *model.User,
	notebook *model.Notebook) (*model.Notebook, *management.Error) {
	mlog.Debug("Create: notebook:%+v", notebook)

	// Create external notebook.
	externalNotebook, err := provider.toExternalNotebook(notebook)

	if err != nil {
		mlog.Error("To external notebook returned error: %v", err)
		return nil, model.ErrorToExternalNotebook
	}

	// Store in database.
	createdNotebook, mErr := provider.notebooksClient.CreateNotebook(externalNotebook)

	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("Create notebook returned error: %v", mErr))
	}

	// Create access information for the notebook owner using authenticated user.
	access := &notebooksModel.Access{
		AccountId:  user.AccountId,
		UserId:     user.Id,
		Permission: notebooksModel.OwnerPermissions,
		NotebookId: createdNotebook.Id,
	}

	// Store in database.
	if _, mErr := provider.notebooksClient.CreateAccess(access); mErr != nil {
		// In case of error. Attempt to delete the notebook.
		provider.notebooksClient.DeleteNotebook(externalNotebook.Id)
		return nil,
			management.GetExternalError(fmt.Sprintf("Create notebook access returned error: %v", mErr))
	}

	// Update the notebook we return to the user.
	notebook.Id = createdNotebook.Id
	notebook.CreatedOn = createdNotebook.CreatedOn.Format(time.RFC3339Nano)
	notebook.LastUpdated = utils.GetTimeUUIDInISO8601(createdNotebook.Version)
	notebook.Permissions = model.OwnerPermissions

	return notebook, nil
}

// Updates an existing notebook.
func (provider *NorthStarNotebooksProvider) Update(user *model.User,
	notebook *model.Notebook) *management.Error {
	mlog.Debug("Update")

	access, mErr := provider.getAccess(user, notebook.Id)

	if mErr != nil {
		return mErr
	}

	// Update requires ownership or read,write permissions.
	switch access.Permission {
	case notebooksModel.OwnerPermissions:
	case notebooksModel.ReadWritePermissions:
	case notebooksModel.ReadWriteExecutePermissions:
	default:
		// For every other permissions, return error.
		return model.ErrorNoNotebookPermissions
	}

	// Create external notebook.
	externalNotebook, err := provider.toExternalNotebook(notebook)

	if err != nil {
		mlog.Error("To external notebook returned error: %v", err)
		return model.ErrorToExternalNotebook
	}

	// Update database.
	if _, mErr := provider.notebooksClient.UpdateNotebook(externalNotebook.Id,
		externalNotebook); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Update notebook returned error: %v", mErr))
	}

	return nil
}

// Get returns notebook for the specified id.
func (provider *NorthStarNotebooksProvider) Get(user *model.User,
	notebookID string) (*model.Notebook, *management.Error) {
	mlog.Debug("Get")

	access, mErr := provider.getAccess(user, notebookID)

	if mErr != nil {
		return nil, mErr
	}

	// Note that this should return only one entry.
	if access.Permission == notebooksModel.NoPermissions {
		return nil, model.ErrorNoNotebookPermissions
	}

	// Note that by default. If permissions are found, user can read.
	externalNotebook, mErr := provider.notebooksClient.GetNotebook(notebookID)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get notebook returned error: %v", mErr))
	}

	// Get notebook from external representation.
	notebook, err := provider.fromExternalNotebook(externalNotebook, access.Permission)

	if err != nil {
		mlog.Error("From external returned error: %v", mErr)
		return nil, model.ErrorFromExternalNotebook
	}

	return notebook, nil
}

// List returns list of notebooks.
func (provider *NorthStarNotebooksProvider) List(user *model.User) ([]model.Notebook, *management.Error) {
	mlog.Debug("List")

	// Query access information for the user.
	query := &notebooksModel.Query{
		AccountId: user.AccountId,
		UserId:    user.Id,
	}

	accesses, mErr := provider.notebooksClient.QueryAccess(query)

	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("Query notebook access returned error: %v", mErr))
	}

	var notebooks []model.Notebook

	// Get notebooks for the current user.
	for _, access := range accesses {
		externalNotebook, mErr := provider.notebooksClient.GetNotebook(access.NotebookId)
		if mErr != nil {
			// This should never happen. But if the two tables get out of sync we delete
			// access information and allow the user to continue.
			if mErr.Id == management.ERR_NOT_FOUND {
				mlog.Info("Warning, the notebook with id %s was not found. Access information will be deleted.",
					access.NotebookId)
				provider.notebooksClient.DeleteAccess(access.Id)
				continue
			}

			mlog.Error("Get notebook with id %s returned error: %v", access.NotebookId, mErr)
			continue
		}

		// Get notebook from external representation.
		notebook, err := provider.fromExternalNotebook(externalNotebook, access.Permission)

		if err != nil {
			mlog.Error("From external returned error: %v", mErr)
			continue
		}

		notebooks = append(notebooks, *notebook)
	}

	return notebooks, nil
}

// Deletes notebook with specified id.
func (provider *NorthStarNotebooksProvider) Delete(user *model.User,
	notebookId string) *management.Error {
	mlog.Debug("Delete")

	access, mErr := provider.getAccess(user, notebookId)
	if mErr != nil {
		return mErr
	}
	mlog.Debug("Access received: %v", access)

	// Only allow the owner to delete the notebook.
	if access.Permission != notebooksModel.OwnerPermissions {
		return model.ErrorNoNotebookPermissions
	}

	// Before deleting the notebook, remove access data.
	if mErr := provider.notebooksClient.DeleteAccess(access.Id); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete access returned error: %v", mErr))
	}

	mlog.Debug("Access %v removed", access.Id)

	// Delete the notebook.
	if mErr := provider.notebooksClient.DeleteNotebook(notebookId); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete notebook returned error: %v", mErr))
	}

	mlog.Debug("Notebook %v deleted", notebookId)
	return nil
}

// GetUsers returns the users associated with the specified notebook id.
func (provider *NorthStarNotebooksProvider) GetUsers(user *model.User,
	notebookId string) ([]model.User, *management.Error) {
	mlog.Debug("GetUsers")

	// Query access information for the notebook.
	query := &notebooksModel.Query{
		NotebookId: notebookId,
	}

	accesses, mErr := provider.notebooksClient.QueryAccess(query)

	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("Query notebook access returned error: %v", mErr))
	}

	// Note that we allow only the owner of the notebook to see the permissions.
	permissions := notebooksModel.NoPermissions

	for _, access := range accesses {
		if access.UserId == user.Id {
			permissions = access.Permission
			break
		}
	}

	if permissions != notebooksModel.OwnerPermissions {
		return nil, model.ErrorNoNotebookUsrPermission
	}

	// Otherwise, create collection of users.
	users := []model.User{}

	for _, access := range accesses {
		user := model.User{
			Id:          access.UserId,
			AccountId:   access.AccountId,
			Permissions: access.Permission,
		}
		users = append(users, user)
	}

	return users, nil
}

// GetExecutionInformation returns information needed to execute the notebook with id.
func (provider *NorthStarNotebooksProvider) GetExecutionInformation(user *model.User,
	notebookId string) (*model.User, *management.Error) {
	mlog.Debug("GetExecutionInformation")

	// Query access information for the notebook.
	query := &notebooksModel.Query{
		NotebookId: notebookId,
	}

	accesses, mErr := provider.notebooksClient.QueryAccess(query)

	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("Query notebook access returned error: %v", mErr))
	}

	// Note that we allow only users with execute permissions to access information.
	permissions := notebooksModel.NoPermissions

	for _, access := range accesses {
		if access.UserId == user.Id {
			permissions = access.Permission
			break
		}
	}

	// If current user is owner or can execute, return user information.
	switch permissions {
	case notebooksModel.OwnerPermissions:
	case notebooksModel.ReadWriteExecutePermissions:
	case notebooksModel.ReadExecute:
	default:
		return nil, model.ErrorNoNotebookExecPermission
	}

	// Search for owner.
	for _, access := range accesses {
		if access.Permission == notebooksModel.OwnerPermissions {
			user := &model.User{
				Id:          access.UserId,
				AccountId:   access.AccountId,
				Permissions: access.Permission,
			}

			return user, nil
		}
	}

	// Note that this is an internal error and should never happen.
	return nil, model.ErrorNotebookOwnerNotFound
}

// UpdateUsers updates the users associated with the specified notebook id.
func (provider *NorthStarNotebooksProvider) UpdateUsers(user *model.User,
	notebookId string,
	users []model.User) *management.Error {
	mlog.Debug("UpdateUsers")

	// Query access information for the notebook.
	query := &notebooksModel.Query{
		NotebookId: notebookId,
	}

	accesses, mErr := provider.notebooksClient.QueryAccess(query)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Query notebook access returned error: %v", mErr))
	}

	// Note that we allow only the owner of the notebook to update the permissions.
	permissions := notebooksModel.NoPermissions
	accessMap := map[string]notebooksModel.Access{}

	for _, access := range accesses {
		// Capture the permissions of the current (authenticated) user.
		if access.UserId == user.Id {
			permissions = access.Permission
		}

		// Add to the map of access. Note that we ignore the owner. E.g.,
		// owner cannot be updated.
		if access.Permission != notebooksModel.OwnerPermissions {
			accessMap[access.UserId] = access
		}
	}

	// If the current (authenticated) user is not the owner, do not allow to this method.
	if permissions != notebooksModel.OwnerPermissions {
		return model.ErrorNoNotebookUsrPermission
	}

	ownerId := user.Id
	newAccesses := []notebooksModel.Access{}
	updatedAccesses := []notebooksModel.Access{}

	// Update the maps used to carry the actual operations.
	for _, user := range users {

		// Make sure owner is ignore.
		if user.Id == ownerId {
			continue
		}

		// Make sure the request does not attempt to set a new owner.
		if user.Permissions == notebooksModel.OwnerPermissions {
			return management.GetForbiddenError("New owners cannot be assign to the notebook.")
		}

		// If the actual user already have access to the notebook.
		// Verify if updates is needed. Otherwise, it is a new
		// user.

		if access, found := accessMap[user.Id]; found == true {
			if user.Permissions != access.Permission {
				mlog.Debug("Updating permissions (%s) for user to %s.",
					access.Permission, user.Id, user.Permissions)
				access.Permission = user.Permissions
				updatedAccesses = append(updatedAccesses, access)
			}

			// Remove from the collection. Note that we assume left overs are to be deleted.
			delete(accessMap, user.Id)
		} else {
			mlog.Debug("Creating permissons for user %s.", user.Id)
			newAccess := notebooksModel.Access{
				AccountId:  user.AccountId,
				UserId:     user.Id,
				Permission: user.Permissions,
				NotebookId: notebookId,
			}
			newAccesses = append(newAccesses, newAccess)
		}
	}

	// Create new access entries.
	for _, access := range newAccesses {
		if _, mErr := provider.notebooksClient.CreateAccess(&access); mErr != nil {
			return management.GetExternalError(fmt.Sprintf("Create access returned error: %v", mErr))
		}
	}

	// Update existing entries.
	for _, access := range updatedAccesses {
		if _, mErr := provider.notebooksClient.UpdateAccess(access.Id, &access); mErr != nil {
			return management.GetExternalError(fmt.Sprintf("Update access returned error: %v", mErr))
		}
	}

	// Delete existing entries.
	for _, access := range accessMap {
		if mErr := provider.notebooksClient.DeleteAccess(access.Id); mErr != nil {
			return management.GetExternalError(fmt.Sprintf("Delete access returned error: %v", mErr))
		}
	}

	return nil
}

// Helper method used to get user access information for a notebook id.
func (provider *NorthStarNotebooksProvider) getAccess(user *model.User,
	notebookId string) (*notebooksModel.Access, *management.Error) {
	mlog.Debug("getAccess")

	// Query access information for the notebook.
	query := &notebooksModel.Query{
		AccountId:  user.AccountId,
		UserId:     user.Id,
		NotebookId: notebookId,
	}

	accesses, mErr := provider.notebooksClient.QueryAccess(query)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Query notebook access returned error: %v", mErr))
	}

	if len(accesses) > 0 {
		return &accesses[0], nil
	}

	// By default, user has no access.
	return &notebooksModel.Access{Permission: notebooksModel.NoPermissions}, nil
}

// Helper method used to translate portal api model to data service model.
func (provider *NorthStarNotebooksProvider) toExternalNotebook(notebook *model.Notebook) (*notebooksModel.Notebook, error) {
	mlog.Debug("toExternalNotebook: notebook:%+v", notebook)

	// Note that we store notebook content as JSON.
	data, err := json.Marshal(notebook)

	if err != nil {
		return nil, err
	}

	// Create the data service representation.
	externalNotebook := &notebooksModel.Notebook{
		Id:      notebook.Id,
		Version: notebook.Etag,
		Data:    string(data),
	}

	return externalNotebook, nil
}

// Helper method used to translate data service model to portal api model.
func (provider *NorthStarNotebooksProvider) fromExternalNotebook(externalNotebook *notebooksModel.Notebook,
	permissions string) (*model.Notebook, error) {
	mlog.Debug("fromExternalNotebook: externalNotebook:%+v", externalNotebook)

	// Unmarshal notebook representation.
	notebook := &model.Notebook{}

	if err := json.Unmarshal([]byte(externalNotebook.Data), notebook); err != nil {
		return nil, err
	}

	notebook.Id = externalNotebook.Id
	notebook.CreatedOn = externalNotebook.CreatedOn.Format(time.RFC3339Nano)
	notebook.LastUpdated = utils.GetTimeUUIDInISO8601(externalNotebook.Version)
	notebook.Etag = externalNotebook.Version
	notebook.Permissions = permissions

	return notebook, nil
}
