package main

import "context"

type appConfig struct {
	co *ColorOutput
}

type appConfigContextKeyType struct{}

func appConfigFromContext(ctx context.Context) *appConfig {
	return ctx.Value(appConfigContextKey()).(*appConfig)
}

func appConfigColorOutput(ctx context.Context) *ColorOutput {
	return ctx.Value(appConfigContextKey()).(*appConfig).co
}

func appConfigContextKey() appConfigContextKeyType {
	return appConfigContextKeyType{}
}
