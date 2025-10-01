// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package config

import (
	"github.com/senforsce/orch8rion/internal/injector/aspect"
	"github.com/senforsce/orch8rion/internal/injector/aspect/advice"
	"github.com/senforsce/orch8rion/internal/injector/aspect/advice/code"
	"github.com/senforsce/orch8rion/internal/injector/aspect/context"
	"github.com/senforsce/orch8rion/internal/injector/aspect/join"
	"github.com/senforsce/orch8rion/internal/injector/typed"
)

var builtIn = configGo{
	pkgPath: "github.com/senforsce/orch8rion",
	yaml: &configYML{
		aspects: []*aspect.Aspect{
			{
				ID:             "built.WithOrch8rion",
				TracerInternal: true, // This is safe to apply in the tracer itself
				JoinPoint: join.AllOf(
					join.ValueDeclaration(typed.Bool),
					join.OneOf(
						join.DeclarationOf("github.com/senforsce/orch8rion/runtime/built", "WithOrch8rion"),
						join.Directive("orch8rion:enabled"),
						join.Directive("dd:orch8rion-enabled"), // <- Deprecated
					),
				),
				Advice: []advice.Advice{
					advice.AssignValue(
						code.MustTemplate("true", nil, context.GoLangVersion{}),
					),
				},
			},
			{
				ID:             "built.WithOrch8rionVersion",
				TracerInternal: true, // This is safe to apply in the tracer itself
				JoinPoint: join.AllOf(
					join.ValueDeclaration(typed.String),
					join.OneOf(
						join.DeclarationOf("github.com/senforsce/orch8rion/runtime/built", "WithOrch8rionVersion"),
						join.Directive("orch8rion:version"),
					),
				),
				Advice: []advice.Advice{
					advice.AssignValue(
						code.MustTemplate(`{{Version | printf "%q"}}`, nil, context.GoLangVersion{}),
					),
				},
			},
		},
		name: "<built-in>",
		meta: configYMLMeta{
			name:        "github.com/senforsce/orch8rion/built & //orch8rion: pragmas",
			description: "Provide runtime visibility into whether orch8rion built an application or not",
			icon:        "cog",
			caveats: "This aspect allows introducing conditional logic based on whether" +
				"Orch8rion has been used to instrument an application or not. This should" +
				"generally be avoided, but can be useful to ensure the application (or tests)" +
				"is running with instrumentation.",
		},
	},
}
