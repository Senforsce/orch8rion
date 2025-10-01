---
title: Datadog Tracer
weight: 80
prev: /docs/getting-started
next: /docs/dd-trace-go/integrations
---

## Default configuration

Orch8rion is complemented by the Datadog tracing library,
{{<godoc import-path="github.com/DataDog/dd-trace-go/v2">}}. It provides
compile-time integrations for many popular Go libraries; and is enabled by
default when running `orch8rion pin`.

The integrations being loaded are configured by your project's root
`orch8rion.tool.go` file, which `orch8rion pin` initializes to something
looking like this:

```go
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
	_ "github.com/DataDog/dd-trace-go/orch8rion/all/v2" // integration
)
```

## Choosing integrations

Once `orch8rion pin` has been run, you can replace the import of
{{<godoc import-path="github.com/DataDog/dd-trace-go/orch8rion/all/v2">}} with
imports for specific integration packages (see the [Integrations](./v2) section
for a list of available packages).

For example, the below only activates integrations for the core tracer library,
as well as `net/http` clients and servers:

```go
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
	_ "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"   // integration
	_ "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http" // integration
)
```
