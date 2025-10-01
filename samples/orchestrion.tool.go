// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// This file was created by `orch8rion pin`, and is used to ensure the
// `go.mod` file contains the necessary entries to ensure repeatable builds when
// using `orch8rion`. It is also used to set up which tracer integrations are
// enabled.

//go:build tools

//go:generate go run github.com/senforsce/orch8rion pin

package tools

// Imports in this file determine which tracer integrations are enabled in
// orch8rion. New integrations can be automatically discovered by running
// `orch8rion pin` again. You can also manually add new imports here to
// enable additional integrations. When doing so, you can run `orch8rion pin`
// to make sure manually added integrations are valid (i.e, the imported package
// includes a valid `orch8rion.yml` file).
import (
	// Ensures `orch8rion` is present in `go.mod` so that builds are repeatable.
	// Do not remove.
	_ "github.com/senforsce/orch8rion"

	// Provides integrations for essential `orch8rion` features. Most users
	// should not remove this integration.
	_ "github.com/senforsce/orch8rion/instrument" // integration
)
