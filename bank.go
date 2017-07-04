// Copyright 2017 Christoph Berger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bank

import (
	"encoding/gob"
	"os"

	"github.com/pkg/errors"
)

// Account is a bank account with a name, a balance, and a
// transaction history.
type Account struct {
	Name string
	Bal  int
	Hist []history
}

type history struct {
	Amt, Bal int
}

var accounts map[string]*Account

// NewAccount creates a new account with a name. Initial balance is 0.
// The new account is added to the bank's map of accounts.
func NewAccount(s string) *Account {
	if accounts == nil {
		accounts = make(map[string]*Account)
	}
	a := &Account{Name: s}
	accounts[s] = a
	return a
}

// GetAccount receives a name and returns the account of that name, if it exists.
// GetAccount panics if the bank has no accounts.
func GetAccount(name string) (*Account, error) {
	accnt, ok := accounts[name]
	if !ok {
		return nil, errors.New("account '" + name + "' does not exist")
	}
	return accnt, nil
}

// Name returns the name of account a.
func Name(a *Account) string {
	return a.Name
}

// Balance returns the current balance of account a.
func Balance(a *Account) int {
	return a.Bal
}

// Deposit adds amount m to account a's balance.
// The amount must be positive.
func Deposit(a *Account, m int) (int, error) {
	if m < 0 {
		return a.Bal, errors.Errorf("Deposit: amount must be positive, but is %d.", m)
	}
	a.Bal += m
	a.Hist = append(a.Hist, history{m, a.Bal})
	return a.Bal, nil
}

// Withdraw removes amount m from account a's balance.
// The amount must be positive.
func Withdraw(a *Account, m int) (int, error) {
	if m < 0 {
		return a.Bal, errors.Errorf("Withdraw: amount must be positive, but is %d.", m)
	}
	if m > a.Bal {
		return a.Bal, errors.Errorf("Withdraw: amount (%d) must be less than actual balance (%d).", m, a.Bal)
	}
	a.Bal -= m
	a.Hist = append(a.Hist, history{-m, a.Bal})
	return a.Bal, nil
}

// Transfer transfers amount m from account a to account b.
// The amount must be positive.
// The sending account must have at least as much money as the
// amount to be transferred.
func Transfer(a, b *Account, m int) (int, int, error) {
	switch {
	case m < 0:
		return a.Bal, b.Bal, errors.Errorf("Transfer: amount must be positive, but is %d.", m)
	case m > a.Bal:
		return 0, a.Bal, errors.Errorf("Withdraw: amount (%d) must be less than actual balance of sending account (%d).", m, a.Bal)
	}
	a.Bal -= m
	b.Bal += m
	a.Hist = append(a.Hist, history{-m, a.Bal})
	b.Hist = append(b.Hist, history{m, b.Bal})
	return a.Bal, b.Bal, nil
}

// History returns a closure that returns one account transaction at a time.
// On each call, the closure returns the amount of the transaction, the resulting balance,
// and a boolean that is true as long as there are more history elements to read.
// The closure returns the history items from oldest to newest.
func History(a Account) func() (int, int, bool) {
	i := 0
	more := true
	return func() (int, int, bool) {
		if i >= len(a.Hist)-1 {
			more = false
		}
		h := a.Hist[i]
		i++
		return h.Amt, h.Bal, more
	}
}

// Persist the accounts map on disk.
func Save() error {
	f, err := os.OpenFile("bank.data", os.O_WRONLY, 0666) // Note: octal #
	if err != nil {
		f, err = os.Create("bank.data")
		if err != nil {
			return errors.Wrap(err, "Save: Create failed")
		}
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	err = e.Encode(accounts)
	if err != nil {
		return errors.Wrap(err, "Save: Encode failed")
	}
	return nil
}

// Restore the accounts map from disk.
func Load() error {
	f, err := os.Open("bank.data")
	if err != nil {
		if os.IsNotExist(err) {
			// Expected. The file does not exist initially.
			return nil
		}
		return errors.Wrap(err, "Load: Open failed")
	}
	defer f.Close()
	d := gob.NewDecoder(f)
	err = d.Decode(&accounts)
	if err != nil {
		return errors.Wrap(err, "Load: Decode failed")
	}
	return nil
}
