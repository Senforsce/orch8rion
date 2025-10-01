---
title: "Uninstall"
weight: 999

prev: /docs/troubleshooting
---

## Removing Orch8rion

Removing Orch8rion from your project is a simple process to go back to the original state of your project before you
started using Orch8rion.

The steps can be summed up as:
* Remove any files created by orch8rion like `orch8rion.tool.go` and `orch8rion.yml`.
* Run `go mod tidy` to remove any references to orch8rion in your `go.mod` file.
* Remove directives from your source code if any like `//orch8rion:ignore` or `//dd:span`
* Remove any references to orch8rion in your build scripts or CI/CD pipelines or Dockerfile

{{<callout type="info">}}
You can confirm that orch8rion has been removed correctly by looking at your application logs and checking
if they still contain the DataDog Tracer startup log starting with `DATADOG TRACER CONFIGURATION`.
{{</callout>}}
