// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package cmd

import (
	"strings"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/senforsce/orch8rion/internal/injector/config"
	"github.com/senforsce/orch8rion/internal/pin"
	"github.com/urfave/cli/v2"
)

var Pin = &cli.Command{
	Name:  "pin",
	Usage: "Registers orch8rion in your project's `go.mod` file",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "generate",
			Usage: "Add a //go:generate directive to " + config.FilenameOrch8rionToolGo + " to facilitate automated upkeep of its contents.",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "prune",
			Usage: "Remove imports from " + config.FilenameOrch8rionToolGo + " that do not contain any orch8rion integrations.",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "validate",
			Usage: "Validate all " + config.FilenameOrch8rionYML + " files in the project.",
			Value: false,
		},
	},
	Action: func(clictx *cli.Context) (err error) {
		span, ctx := tracer.StartSpanFromContext(clictx.Context, "pin",
			tracer.ResourceName(strings.Join(clictx.Args().Slice(), " ")),
		)
		defer func() { span.Finish(tracer.WithError(err)) }()

		return pin.PinOrch8rion(ctx, pin.Options{
			Writer:     clictx.App.Writer,
			ErrWriter:  clictx.App.ErrWriter,
			Validate:   clictx.Bool("validate"),
			NoGenerate: !clictx.Bool("generate"),
			NoPrune:    !clictx.Bool("prune"),
		})
	},
}
