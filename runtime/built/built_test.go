// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package built_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/senforsce/orch8rion/internal/injector/config"
	"github.com/senforsce/orch8rion/internal/version"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	tmp := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tmp, "main.go"), []byte(testProgram), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(tmp, config.FilenameOrch8rionToolGo), []byte(orch8rionToolGo), 0o644))

	cmd := exec.Command("go", "mod", "init", "dummy.test")
	cmd.Dir = tmp
	require.NoError(t, cmd.Run())

	_, thisFile, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(thisFile, "..", "..", "..")
	cmd = exec.Command("go", "mod", "edit", "-replace=github.com/senforsce/orch8rion="+rootDir)
	cmd.Dir = tmp
	require.NoError(t, cmd.Run())

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmp
	require.NoError(t, cmd.Run())

	var stdout bytes.Buffer
	cmd = exec.Command("go", "run", "github.com/senforsce/orch8rion", "go", "run", ".")
	cmd.Dir = tmp
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())

	require.Equal(t, version.Tag(), stdout.String())
}

const testProgram = `package main

import (
	"fmt"
	"os"

	"github.com/senforsce/orch8rion/runtime/built"
)

func main() {
	if !built.WithOrch8rion {
		os.Exit(42)
	}

	fmt.Print(built.WithOrch8rionVersion)
}
`

const orch8rionToolGo = `//go:build tools

package tools

import (
	_ "github.com/senforsce/orch8rion"
)
`
