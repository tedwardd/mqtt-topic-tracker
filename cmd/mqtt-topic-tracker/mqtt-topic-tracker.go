package main

import (
	"os"

	"github.com/fatih/color"
)

var hostname, _ = os.Hostname()

var channel = make(chan os.Signal, 1)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
)
