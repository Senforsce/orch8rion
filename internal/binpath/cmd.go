// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package binpath

import (
	"os"
	"path/filepath"
)

var Orch8rion string

func init() {
	var err error
	if Orch8rion, err = os.Executable(); err != nil {
		if Orch8rion, err = filepath.Abs(os.Args[0]); err != nil {
			Orch8rion = os.Args[0]
		}
	}
	Orch8rion = filepath.Clean(Orch8rion)
}
