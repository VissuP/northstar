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
	"sync"
	"time"

	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/client/auth"
	"github.com/verizonlabs/northstar/nssim/client/user"
	"github.com/verizonlabs/northstar/nssim/config"
	"strings"
)

var (
	DefaultExpiration = 1 * time.Hour
	once              sync.Once
	accountManager    *AccountManager // singleton
)

//AccountPtr defines the type used to keep track of account reference.
type AccountPtr struct {
	Account *user.Account
	Token   auth.Token
}

//Stats defines type used to capture some account metrics.
type Stats struct {
	CreatedAccounts     int
	FailedCreateAccount int
}

// AccountManager defines the (private) type used to manager account references.
type AccountManager struct {
	sync.RWMutex
	expiration       time.Duration
	accounts         []*AccountPtr
	stats            Stats
}

//GetAccountManager returns the singleton instance of the account manager.
func GetAccountManager() *AccountManager {

	// Create the, one and only one, account manager instance.
	once.Do(func() {
		mlog.Debug("Creating new account manager")

		accountManager = &AccountManager{
			expiration:       DefaultExpiration,
			accounts:         make([]*AccountPtr, 0, 0),
		}
	})

	return accountManager
}

//CreateAccounts returns a new account.
func (manager *AccountManager) CreateAccounts(force bool, count int) ([]*user.Account, error) {
	mlog.Info("CreateAccount: force - %b", force)
	manager.Lock()
	defer manager.Unlock()

	accounts := []*user.Account{}
	for i := 0; i < count; i++ {
		email := config.Configuration.Credentials.User.Id + fmt.Sprintf("%d", i) + "@nssim.verizon.com"
		password := config.Configuration.Credentials.User.AccountPassword
		mlog.Info("Getting account. Username: %s, password: %s", email, password)

		if len(manager.accounts) > count {
			mlog.Info("Account already exists. Appending.")
			accounts = append(accounts, manager.accounts[i].Account)
			continue
		}

		// Create a new account.
		account, err := manager.createAccount(email, password)
		if err != nil {
			if !strings.Contains(err.Error(), "Loginname already used for id") {
				mlog.Info("Create account response: %s", err.Error())
				return nil, err
			}
			account = &user.Account{}
		}

		account.Email = email
		account.Password = password

		accountPtr := &AccountPtr{
			Account: account,
		}

		mlog.Info("Account (%s) for user %s created successfully.", account.Id, email)
		//Append to the list of accounts to return
		accounts = append(accounts, account)

		//Save the account for future use.
		manager.accounts = append(manager.accounts, accountPtr)
	}

	return accounts, nil
}

//LogStats logs status associated with account manager
func (manager *AccountManager) LogStats() {
	mlog.Info("Account Manager Stats: %+v", manager.stats)
}

//createAccount is a helper method used to create the account.
func (manager *AccountManager) createAccount(email, password string) (*user.Account, error) {
	// Get the token associated with this test application.
	_ , mErr := auth.GetClientToken()
	if mErr != nil {
		return nil, fmt.Errorf("Failed to get client token with error: %s", mErr.Description)
	}

        account := &user.Account{} //TEMP make this compile; stubbing out real account
 
	manager.stats.CreatedAccounts++
	return account, nil
}
