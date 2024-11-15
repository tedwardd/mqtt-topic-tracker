package main

import (
	"io"

	"github.com/tedwardd/mqtt-topic-tracker/internal/logger"
)

type commandBasics struct {
	LogCache *logger.LogCache
	verbose  bool
	debug    bool
}

func (basics *commandBasics) setupLogCache(out io.Writer, err io.Writer) {
	basics.LogCache = logger.NewLogCache(out, err)
}
