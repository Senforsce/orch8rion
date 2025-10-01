---
title: Configuration
weight: 20
---

## Configuration

Orch8rion uses the file `orch8rion.tool.go` in the project root to import its configuration. This file is
used to specify which integrations are enabled, and which features are activated but also as a way to connect
orch8rion integrations to your `go.mod` file so they are appropriately versioned.

The file is a Go source file, and it must be valid Go code. If your project does not already contain one, you may run
`orch8rion pin` to create it.

{{<callout type="info">}}
Orch8rion is a vendor-agnostic tool. By default, `orch8rion pin` enables Datadog's tracer integrations by
importing `github.com/DataDog/dd-trace-go/orch8rion/all/v2` in `orch8rion.tool.go`, but other vendors (such as OpenTelemetry) may
provide alternate integrations that can be used instead.
{{</callout>}}

### Loading

Configuration is loaded from the `orch8rion.tool.go` located in the same directory as your application's `go.mod` file. Each import in this file
will be processed by orch8rion and will enable the corresponding integration. Configuration loading happens
recursively and will load all the integrations that are imported by `orch8rion.tool.go` in the imported package in a
tree-like structure (packages are de-duplicated so you don't have to worry about a package being transitively imported by multiple integrations).

```mermaid
flowchart TD
    root --> github.com/DataDog/dd-trace-go/orch8rion/all
    github.com/DataDog/dd-trace-go/orch8rion/all --> ddtrace/tracer
    github.com/DataDog/dd-trace-go/orch8rion/all --> contrib/net/http
    github.com/DataDog/dd-trace-go/orch8rion/all --> contrib/database/sql
    github.com/DataDog/dd-trace-go/orch8rion/all --> ...
```

Each package encountered in the configuration loading step is allowed to contain an `orch8rion.yml` file. These files
are the auto-instrumentation configuration backbone that modify your codebase. Please refer to the
[contributing guide][contributing] for more details on how to write these.

[contributing]: ../contributing/

### Finer grain instrumentation

The default `orch8rion.tool.go` imports all integrations provided by the `github.com/DataDog/dd-trace-go/orch8rion/all/v2`
package. But this can be cumbersome if you only want to use a subset of the integrations. You can expand the default
`orch8rion.tool.go` by replacing the import of `github.com/DataDog/dd-trace-go/v2` with the specific integrations you
want to use from the list available at one level deeper in the configuration loading tree [here][orch8rion-all].

[orch8rion-all]: https://github.com/DataDog/dd-trace-go/blob/main/orch8rion/all/orch8rion.tool.go

### Remove an integration

Sometimes auto-instrumentation simply does not fit your use case. Plenty of automatic instrumentation modules offer more
configuration option when using their SDK. If you plan on using an SDK integration you should first remove the
corresponding import from `orch8rion.tool.go` and then use the SDK's own configuration mechanism to enable it. This
may require you to opt for finer grain instrumentation like described in the previous section.

{{<callout type="warning">}}
Some auto instrumentation integrations (notably caller-side) may automatically be neutered by using the corresponding
manual instrumentation but this is not guaranteed. If you are using manual instrumentation, and you want to ensure that
2 similar spans are not created, you should remove the corresponding import from `orch8rion.tool.go`.
{{</callout>}}
