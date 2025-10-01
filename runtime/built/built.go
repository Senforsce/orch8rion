// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// Package built provides information about how the current application has been
// built, if it has been built using orch8rion. It provides advanced features
// used internally by the Datadog tracer library to convey accurate telemetry
// data (when the application is opted in to Datadog telemetry at runtime), and
// may be useful in certain advanced use cases. Most users should not need to
// use anything from this package.
package built

// WithOrch8rion is true if the current application was built using
// orch8rion. This is useful to perform certain behavior dependent on whether
// the application was automatically instrumented or not. This can be useful to
// avoid double-instrumentation, or to guarantee the application runs with
// automatic instrumentation. Most users should not need to use this variable.
const WithOrch8rion = false

// WithOrch8rionVersion is the version of orch8rion used to build the
// library, if the application was built by it. It is a blank string otherwise.
// This can be useful context to include in logs when the use of orch8rion is
// relevant. Most users should not need to use this variable.
const WithOrch8rionVersion = ""
