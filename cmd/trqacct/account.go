package main

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Account struct {
	ID       string
	Njobs    uint64
	Walltime uint64
	Memory   uint64
	Cputime  uint64
}

func (a Account) String() string {
	return fmt.Sprintf("|%-12s|%12d|%20d|%20d|", a.ID, a.Njobs, a.Walltime, a.Memory)
}

func AccountTabletHeader() string {
	return AccountTabletSeparator() +
		fmt.Sprintf("\n|%-12s|%12s|%20s|%20s|\n", "ID", "jobs", "used walltime", "requested memory") +
		AccountTabletSeparator()
}

func AccountTabletSeparator() string {
	return fmt.Sprint("+" + strings.Repeat("=", 67) + "+")
}

// FindAccount finds the account registry with given `id`.
//
// It uses the `slices` feature available in Go >= 1.18
func FindAccount(accounts []Account, id string) int {
	return slices.IndexFunc(accounts, func(a Account) bool {
		return a.ID == id
	})
}
