package main

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type ColorOutput struct {
	stdout    io.Writer
	stderr    io.Writer
	infoColor *color.Color
	errColor  *color.Color
	warnColor *color.Color
}

func initColorOutput(stdout io.Writer, stderr io.Writer) *ColorOutput {
	return &ColorOutput{
		stdout:    stdout,
		stderr:    stderr,
		infoColor: green.Add(color.Bold),
		errColor:  red.Add(color.Bold),
		warnColor: yellow.Add(color.Bold),
	}
}

// Info
func (c ColorOutput) Info(msg string) {
	_, _ = c.infoColor.Fprintln(c.stdout, msg)
}

// Error
func (c ColorOutput) Error(msg string) {
	_, _ = c.errColor.Fprintln(c.stderr, msg)
}

// Errorf paints error messages red
func (c ColorOutput) Errorf(msg string, vars ...interface{}) error {
	return fmt.Errorf(c.errColor.Sprintf(msg, vars...))
}

// Warn
func (c ColorOutput) Warn(msg string) {
	_, _ = c.warnColor.Fprintln(c.stdout, msg)
}
