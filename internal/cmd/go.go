// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package cmd

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/senforsce/orch8rion/internal/binpath"
	"github.com/senforsce/orch8rion/internal/goproxy"
	"github.com/senforsce/orch8rion/internal/pin"
	"github.com/urfave/cli/v2"
)

var (
	Go = &cli.Command{
		Name:            "go",
		Usage:           "Executes standard go commands with automatic instrumentation enabled",
		UsageText:       "orch8rion go [go command arguments...]",
		Args:            true,
		SkipFlagParsing: true,
		Action: func(clictx *cli.Context) (err error) {
			span, ctx := tracer.StartSpanFromContext(clictx.Context, "go",
				tracer.ResourceName(strings.Join(clictx.Args().Slice(), " ")),
			)
			defer func() { span.Finish(tracer.WithError(err)) }()

			if err := pin.AutoPinOrch8rion(ctx, clictx.App.Writer, clictx.App.ErrWriter); err != nil {
				return cli.Exit(err, -1)
			}

			if err := goproxy.Run(ctx, clictx.Args().Slice(), goproxy.WithToolexec(binpath.Orch8rion, "toolexec")); err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					return cli.Exit(err, exitErr.ExitCode())
				}
				return cli.Exit(err, -1)
			}
			return nil
		},
	}
)
