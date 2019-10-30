// Copyright (c) 2015-2016 The btcsuite developers
// Copyright (c) 2019 The Soteria DAG developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/soteria-dag/soterd/soterutil"
	"github.com/soteria-dag/soterwallet/waddrmgr"
	"github.com/soteria-dag/soterwallet/walletdb"
	_ "github.com/soteria-dag/soterwallet/walletdb/bdb"
	"github.com/soteria-dag/soterwallet/wtxmgr"
)

const defaultNet = "mainnet"

var datadir = soterutil.AppDataDir("soterwallet", false)

// Flags.
var opts = struct {
	Force  bool   `short:"f" description:"Force removal without prompt"`
	DbPath string `long:"db" description:"Path to wallet database"`
}{
	Force:  false,
	DbPath: filepath.Join(datadir, defaultNet, "wallet.db"),
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
}

var (
	// Namespace keys.
	waddrmgrNamespace = []byte("waddrmgr")
	wtxmgrNamespace   = []byte("wtxmgr")
)

func yes(s string) bool {
	switch s {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}

func no(s string) bool {
	switch s {
	case "n", "N", "no", "No":
		return true
	default:
		return false
	}
}

func main() {
	os.Exit(mainInt())
}

func mainInt() int {
	fmt.Println("Database path:", opts.DbPath)
	_, err := os.Stat(opts.DbPath)
	if os.IsNotExist(err) {
		fmt.Println("Database file does not exist")
		return 1
	}

	for !opts.Force {
		fmt.Print("Drop all soterwallet transaction history? [y/N] ")

		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		if !scanner.Scan() {
			// Exit on EOF.
			return 0
		}
		err := scanner.Err()
		if err != nil {
			fmt.Println()
			fmt.Println(err)
			return 1
		}
		resp := scanner.Text()
		if yes(resp) {
			break
		}
		if no(resp) || resp == "" {
			return 0
		}

		fmt.Println("Enter yes or no.")
	}

	db, err := walletdb.Open("bdb", opts.DbPath)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return 1
	}
	defer db.Close()

	fmt.Println("Dropping soterwallet transaction history")

	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
		err := tx.DeleteTopLevelBucket(wtxmgrNamespace)
		if err != nil && err != walletdb.ErrBucketNotFound {
			return err
		}
		ns, err := tx.CreateTopLevelBucket(wtxmgrNamespace)
		if err != nil {
			return err
		}
		err = wtxmgr.Create(ns)
		if err != nil {
			return err
		}

		ns = tx.ReadWriteBucket(waddrmgrNamespace)
		birthdayBlock, err := waddrmgr.FetchBirthdayBlock(ns)
		if err != nil {
			fmt.Println("Wallet does not have a birthday block " +
				"set, falling back to rescan from genesis")

			startBlock, err := waddrmgr.FetchStartBlock(ns)
			if err != nil {
				return err
			}
			return waddrmgr.PutSyncedTo(ns, startBlock)
		}

		return waddrmgr.PutSyncedTo(ns, &birthdayBlock)
	})
	if err != nil {
		fmt.Println("Failed to drop and re-create namespace:", err)
		return 1
	}

	return 0
}