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
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/portal/model"
)

// ListPortfolios returns all portfolios associated with authenticated user. Note that this method only
// returns name and creationDate of the portfolio.
func (provider *NorthStarPortalProvider) ListPortfolios(token string) ([]model.Portfolio, *management.Error) {
	mlog.Debug("ListPortfolios")

	// Get the northstarapi (external) portfolios for user (i.e. from token).
	externalPortfolios, mErr := provider.northstarApiClient.ListBuckets(token)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("List portfolios returned error: %v", mErr))
	}

	// Translate to portal portfolios.
	var portfolios []model.Portfolio

	for _, externalPortfolio := range externalPortfolios {
		portfolio := model.Portfolio{
			Name:         externalPortfolio.Name,
			CreationDate: externalPortfolio.CreationDate,
		}

		portfolios = append(portfolios, portfolio)
	}

	return portfolios, nil
}

// ListFiles returns all files associated with authenticated user for a particular portfolio.
func (provider *NorthStarPortalProvider) ListFiles(token string, portfolio string, prefix string, count int, marker string) ([]model.File, *management.Error) {
	mlog.Debug("ListFiles")

	// Get the northstarapi (external) files for user (i.e. from token).
	externalFiles, mErr := provider.northstarApiClient.ListObjects(token, portfolio, prefix, count, marker)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("List portfolios returned error: %v", mErr))
	}

	// Translate to portal files.
	var files []model.File

	for _, externalFile := range externalFiles {
		file := model.File{
			Name:         externalFile.Key,
			LastModified: externalFile.LastModified,
			Size:         externalFile.Size,
			Etag:         externalFile.Etag,
		}

		files = append(files, file)
	}

	return files, nil
}

// GetFile returns the specified file under the selected portfolio
func (provider *NorthStarPortalProvider) GetFile(token string, portfolio string, file string) (*model.Data, *management.Error) {
	mlog.Debug("GetFile")

	// Get the northstarapi api (external) file under a portfolio.
	externalData, mErr := provider.northstarApiClient.GetObject(token, portfolio, file)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get file with file name %s returned error: %v", file, mErr))
	}

	return fromExternalData(externalData), nil
}

// Helper method used to translate nortstar api model to portal model.
func fromExternalData(externalData *northstarApiModel.Data) *model.Data {
	mlog.Debug("fromExternalData")

	// Create portal model data.
	data := &model.Data{
		Payload:     externalData.Payload,
		ContentType: externalData.ContentType,
	}

	return data
}
