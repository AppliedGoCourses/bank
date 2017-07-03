// Copyright 2017 Christoph Berger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bank

import (
	"errors"
)

type Account struct {
	name string
	bal  int
	hist []history
}

type history struct {
	amt, bal int
}

// NewAccount creates a new account with a name. Initial balance is 0.
func NewAccount(s string) *Account {
	return &Account{name: s}
}

// Name returns the name of the account.
func Name(a Account) string {
	return a.name
}

// Balance returns the current balance of account a.
func Balance(a Account) int {
	return a.bal
}

// Deposit adds amount m to account a's balance.
func Deposit(a *Account, m int) (int, error) {
	if m < 0 {
		return a.bal, errors.New("You can only deposit positive amounts.")
	}
	a.bal += m
	a.hist = append(a.hist, history{m, a.bal})
	return a.bal, nil
}

// Withdraw removes amount m from account a's balance.
func Withdraw(a *Account, m int) (int, error) {
	if m < 0 {
		return a.bal, errors.New("You can only withdraw positive amounts.")
	}
	if m > a.bal {
		return a.bal, errors.New("You cannot take on debt.")
	}
	a.bal -= m
	a.hist = append(a.hist, history{-m, a.bal})
	return a.bal, nil
}

// Transfer transfers amount m from account a to account b.
func Transfer(a, b *Account, m int) (int, int, error) {
	switch {
	case m < 0:
		return a.bal, b.bal, errors.New("You can only transfer positive amounts.")
	case m > a.bal:
		return a.bal, b.bal, errors.New("You cannot take on debt.")
	}
	a.bal -= m
	b.bal += m
	a.hist = append(a.hist, history{-m, a.bal})
	b.hist = append(b.hist, history{m, b.bal})
	return a.bal, b.bal, nil
}

// History returns a closure that returns one account transaction at a time.
// On each call, the closure returns the amount of the transaction, the resulting balance,
// and a boolean that is true as long as there are more history elements to read.
// The closure returns the history items from oldest to newest.
func History(a Account) func() (int, int, bool) {
	i := 0
	more := true
	return func() (int, int, bool) {
		if i >= len(a.hist)-1 {
			more = false
		}
		h := a.hist[i]
		i++
		return h.amt, h.bal, more
	}
}
