---
title: Architecture
weight: 90
---

## How `orch8rion` works

Orch8rion leverages the standard Go toolchain's `-toolexec` flag to intercept
invocations to specific tools of the toolchain as part of the build:

1. Invocations of `go tool compile`
2. Invocations of `go tool link`

It uses a job server to ensure a given package is built exactly once, even if it
is a shared dependency between the instrumented application and some injected
packages. The job server also centralizes calls to {{<godoc
  import-path="golang.org/x/tools/go/packages"
  package="packages"
  name="Load"
>}} that resolve injected package's code objects, so that this somewhat
expensive process is done only once per package.

### Toolchain Versions

```mermaid
sequenceDiagram
  autonumber

  participant Toolchain as go toolchain
  participant Orch8rion as orch8rion toolexec
  participant JobServer as orch8rion job server
  participant Compiler as go tool compile
  participant Linker as go tool link

  Toolchain ->>+ Orch8rion: compile -V=full
  Orch8rion ->>+ Compiler: -V=full
  Compiler -->>- Orch8rion: version string
  Orch8rion ->>+ JobServer: build.version
  Note right of JobServer: Cache miss
  JobServer ->>+ Toolchain: packages.Load
  Toolchain -->>- JobServer: packages
  JobServer -->>- Orch8rion: version suffix
  Orch8rion -->>- Toolchain: full version string

  Toolchain ->>+ Orch8rion: link -V=full
  Orch8rion ->>+ Linker: -V=full
  Linker -->>- Orch8rion: version string
  Orch8rion ->>+ JobServer: build.version
  Note right of JobServer: Cache hit
  JobServer -->>- Orch8rion: version suffix
  Orch8rion -->>- Toolchain: full version string
```

The standard Go toolchain invokes all tools involved in a given build with the
`-V=full` argument (①, ⑨), so it can use all tool's versions as build cache
invalidation inputs. Orch8rion intercepts those calls, and appends information
about itself to the results (④, ⑫). The version information added by
orch8rion changes:

<div class="hextra-code-block hx-relative hx-mt-6 first:hx-mt-0 hx-group/code"><div><div class="highlight"><pre tabindex="0" class="chroma"><code class="language-console" data-lang="console"><span class="k">compile version go1.23.6</span><span class="o">:</span><span class="nb">orch8rion@v1.1.0-rc.1</span><span class="o">;</span><span class="s2">&lt;base64-encoded-hash&gt;</span></code></pre></div></div></div>

- the version of orch8rion being used, as different versions may apply
  integrations differently
- a base64-encoded hash composed using:
  - the specific configuration being used, as different integrations configured
    result in different instrumented code
  - the details about all packages that may be injected by the configured
    integrations, as the Go toolchain is unaware of these dependencies, yet
    they affect the nature of the build output
    * All relevant modules are listed using {{<godoc
        import-path="golang.org/x/tools/go/packages"
        package="packages"
        name="Load"
      >}} (⑤), and the result is cached


This results in more cache invalidations than is strictly necessary, however
the Go toolchain does not currently offer a more granual way to influence build
identifiers used for caching.

### Compilation

```mermaid
sequenceDiagram
  autonumber

  participant Toolchain as go toolchain
  participant Orch8rion as orch8rion toolexec
  participant JobServer as orch8rion job server
  participant Compiler as go tool compile

  loop For each package
    Toolchain ->>+ Orch8rion: compile ${args...}
    note over Toolchain,JobServer: The job server ensures a given package is compiled exactly once
    alt first build of package
      Orch8rion ->>+ JobServer: build.Start
      JobServer ->>- Orch8rion: token
      Orch8rion ->> Orch8rion: instrument .go files
      Orch8rion ->>+ JobServer: packages.Resolve
      note right of JobServer: injected packages
      JobServer ->>+ Toolchain: packages.Load
      Toolchain -->>- JobServer: packages
      JobServer -->>- Orch8rion: archives
      opt When package is "main"
        Orch8rion ->> Orch8rion: write link-deps.go
      end
      Orch8rion ->> Orch8rion: update -importcfg file
      note over Orch8rion,Compiler: Invoke the actual compiler tool
      Orch8rion -->>+ Compiler: ${args...}
      Compiler ->>- Orch8rion: exit code
      Orch8rion ->> Orch8rion: add link.deps to -output file
      Orch8rion ->>+ JobServer: build.Finish
      JobServer -->>- Orch8rion: ack
    else subsequent build of package (idempotent)
      Orch8rion ->>+ JobServer: build.Start
      JobServer ->>- Orch8rion: idempotent
      Orch8rion ->> Orch8rion: Copy build artifacts
    end
    Orch8rion -->>- Toolchain: exit code
  end
```

The standard Go toolchain makes one invocation to `go tool compile` (①) for
each package being built (unless that particular package is already present in
the `GOCACHE`).

Orch8rion begins by registering the package build with the job server (②),
which will determine whether the build is new and should proceed (③); or if it
has already been done and should be re-used from cache (⑰).

When doing the first build of a package (④), orch8rion will:

- parse all `.go` source files using {{<godoc import-path="go/parser">}}
- type-check the {{<godoc import-path="go/ast" name="File">}}
   * this requires reading type information from dependencies using the archives
     listed in the file specified by the `-importcfg` flag
- processing the {{<godoc import-path="go/ast" name="File">}} with the
  configured integrations (they are decorated by {{<godoc
  import-path="github.com/dave/dst" package="dave/dst">}})
  * Modified copies of the files are written in the Go toolchain's working
    directory; and they include `//line` pragmas to retain the original file's
    line information
  * New compile-time dependencies may be introduced at this stage: integrations
    may inject new packages that are not part of the original build's closure,
    and the `-importcfg` file must provide an archive file for each imported
    package. Those dependencies are resolved using {{<godoc
      import-path="golang.org/x/tools/go/packages"
      package="packages"
      name="Load"
    >}} (⑥)
   * New link-time dependencies may be introduced at this stage (via
     `//go:linkname` pragmas), which must be recorded together with the
     package's build artifacts
* When building a `main` package, a new source file is created (⑨) that
  contains `import` statement for all link-time dependencies that were
  previously recorded and which are not present in the `-importcfg` file
  * This is necessary to ensure those package's `func init()` functions are
    correctly registered, and so that the Go toolchain presents those packages'
    archives to the linker
* The `go tool compile` command is executed (⑪), using modified and synthetic
  `.go` source files and the modified `-importcfg` file
* A `link.deps` file is added to the compiler-produced `.a` archive (⑬),
  listing all link-time dependencies implied by a dependency on this package.
  This is performed using `go tool pack`

Finally, the outcome of the build is registered with the job server (⑭),
unblocking concurrent attempts at building the same package.

### Link

```mermaid
sequenceDiagram
  autonumber

  participant Toolchain as go toolchain
  participant Orch8rion as orch8rion toolexec
  participant JobServer as orch8rion job server
  participant Linker as go tool link

  loop For each executable
    Toolchain ->>+ Orch8rion: link ${args...}
    loop For each -importcfg entry
      Orch8rion ->> Orch8rion: read link.deps object
      Orch8rion ->>+ JobServer: packages.Resolve
      note right of JobServer: un-satisfied link-time dependencies
      JobServer ->>+ Toolchain: packages.Load
      Toolchain -->>- JobServer: packages
      JobServer -->>- Orch8rion: archives
    end
    Orch8rion ->> Orch8rion: update -importcfg file
    note over Orch8rion,Linker: Invoke the actual linker tool
    Orch8rion -->>+ Linker: ${args...}
    Linker ->>- Orch8rion: exit code
    Orch8rion -->>- Toolchain: exit code
  end
```

The standard Go toolchain invokes `go tool link` (①) once for each executable
binary being produced. When using `go run` or `go build`, this is a single
invocation; however `go test` will invoke the linker once for each test package.

Orch8rion intercepts the linker commands to update the `-importcfg` file so
that it correctly lists all link-time dependencies introduced by instrumentation
of all linked packages (②). It uses {{<godoc
  import-path="golang.org/x/tools/go/packages"
  package="packages"
  name="Load"
>}} to locate the relevant archive files (④), and writes an updated
`-importcfg` file (⑦) with all necessary additions performed.

Finally, it invokes the `go tool link` with updated arguments (⑧).

## Code Injection

Orch8rion drives code injection using a process similar to classical
Aspect-oriented Programming (AoP) (see [Aspects][contrib-aspects]). These
combine a _Join Point_ (where code needs to be modified) with one or more
_Advice_ (what modifications need to be made).

In order to reduce the cost of evaluation ({{<godoc
import-path="gopkg.in/DataDog/dd-trace-go.v1">}}) ships more than 100 different
aspects), we apply heuristics to determine what aspects have a chance of
applying to any given package and source file. The heuristics are based on the
observable dependency closure of the package being built (there is no need to
consider instrumentation targeting the `net/http` package if that package is not
imported) as well as the content of source files (an aspect that looks for the
`//dd:span` directive will never match in a source file that does not contain
any occurrence of this string).

The injector performs a depth-first traversal of the entire Abstract Syntax
Trees (ASTs), evaluates every applicable join point on each node; and applies
the configured advice where join points match.

[contrib-aspects]: ../../contributing/aspects/

## The job server

Due to the design of the Go toolchain's `-toolexec` feature, orch8rion works
by wrapping a large number of short-lived processes, which makes it difficult to
share state between individual processes.

Some of the work performed during instrumentation can however be expensive, and
we can preserve resources by making sure that work is done exactly once,
regardless of how many times it is required.

Orch8rion addresses this by starting a _job server_, which uses the
[NATS][nats] protocol and stays up for the entire duration of the build. That
server is responsible for the following aspects:

- Computing the version information that is appended to the output of
  intercepted `-V=full` invocations;
- Resolving package archives for injected dependencies, both during the
  `compile` and `link` phases of the build &ndash; these may cause child builds
  to be created;
- Storing `compile` task results in order to avoid having to re-instrument and
  re-compile packages that are both in the build's original dependency closure
  and part of some injected package dependencies.

[nats]: https://nats.io/
