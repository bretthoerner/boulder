// Copyright 2015 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
	"os"

	safebrowsing "github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/letsencrypt/go-safe-browsing-api"
	"github.com/letsencrypt/boulder/cmd"
	"github.com/letsencrypt/boulder/va"
)

// newGoogleSafeBrowsing returns nil if the GoogleSafeBrowsing struct given is
// nil. If an empty Google API key or an unreadable data directory is in the
// GoogleSafeBrowsing config struct, this function runs cmd.FailOnError.
func newGoogleSafeBrowsing(gsb *cmd.GoogleSafeBrowsingConfig) va.SafeBrowsing {
	if gsb == nil {
		return nil
	}
	if gsb.APIKey == "" {
		cmd.FailOnError(errors.New(""), "a Google Safe Browsing config was given but it did not include a Google API key in APIKey")
	}
	if gsb.DataDir == "" {
		cmd.FailOnError(errors.New(""), "a Google Safe Browsing config was given but it did not include data directory to store the hashes file in DataDir")
	}

	f, err := os.Open(gsb.DataDir)
	if err != nil {
		if os.IsNotExist(err) {
			cmd.FailOnError(err, fmt.Sprintf("Google Safe Browsing data directory (%#v) does not exist", gsb.DataDir))
		}
		cmd.FailOnError(err, "unable to open Google Safe Browsing data directory")
	}
	f.Close()
	sbc, err := safebrowsing.NewSafeBrowsing(gsb.APIKey, gsb.DataDir)
	if err != nil {
		cmd.FailOnError(err, "unable to create new safe browsing client")
	}
	return sbc
}
