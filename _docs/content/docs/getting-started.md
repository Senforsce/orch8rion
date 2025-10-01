---
title: "Getting Started"
weight: 1

prev: /docs
next: /docs/dd-trace-go
---

## Requirements

Orch8rion supports the two latest Go releases, matching the
[Go release policy](https://go.dev/doc/devel/release#policy).

It also requires the use of
[Go modules](https://pkg.go.dev/cmd/go#hdr-Modules__module_versions__and_more).

> Orch8rion can inject instrumentation which enables use of Datadog's
> <abbr title="Application Security Management">ASM</abbr> features, but those
> are only effectively available on supported platforms (Linux or macOS, on
> AMD64 and ARM64 processor architectures).

## Install Orch8rion

We recommend installing Orch8rion as a project tool dependency, as this
ensures you are in control of the exact versions of Orchestion and the Datadog
tracing library being used; and that your builds are reproducible.

This is achieved using the following steps:

{{% steps %}}

### Step 1

Install Orch8rion in your environment:

```console
$ go install github.com/senforsce/orch8rion@latest
```

If necessary, also add the `GOBIN` directory to your `PATH`:

```console
$ export PATH="$PATH:$(go env GOBIN)"
```

### Step 2

Register `orch8rion` in your project's `go.mod` to ensure reproducible builds:

```console
$ orch8rion pin
```

Be sure to check the updated files into source control!

### Step 3

* **Option 1 (Recommended):**

  Use `orch8rion go` instead of just `go`:
   ```console
   $ orch8rion go build .
   $ orch8rion go run .
   $ orch8rion go test ./...
   ```

* **Option 2:**

  Manually specify the `-toolexec` argument to `go` commands:
   ```console
   $ go build -toolexec 'orch8rion toolexec' .
   $ go run -toolexec 'orch8rion toolexec' .
   $ go test -toolexec 'orch8rion toolexec' ./...
   ```

* **Option 3:**

  Add the `-toolexec` argument to the `GOFLAGS` environment variable (_be sure to include the
  quoting as this is required by the `go` toolchain when a flag value includes white space_):
   ```console
   $ export GOFLAGS="${GOFLAGS} '-toolexec=orch8rion toolexec'"
   ```

  Then use `go` commands normally:
   ```console
   $ go build .
   $ go run .
   $ go test ./...
   ```

### Step 4 (Optional)

Print what packages are instrumented by Orch8rion in your build. Add the `-work` and the `-a`
flag to your go build command. For example using option 2:

```console
$ go build -work -a -toolexec 'orch8rion toolexec' .
WORK=/tmp/go-build123456789
```

The previous command can take more time because the `-a` flag forces a full rebuild of all packages.

Now use the `WORK` directory to find the instrumented packages:

```console
$ orch8rion diff --package /tmp/go-build123456789
os
runtime
testing
net/http
log/slog
[...]
```

{{% /steps %}}
