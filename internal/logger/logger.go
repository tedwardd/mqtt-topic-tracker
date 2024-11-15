package logger

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fatih/color"
)

type logColors struct {
	statusColor      int
	healthyColor     int
	recommendColor   int
	warningColor     int
	failureColor     int
	errorColor       int
	setupColor       int
	setupPromptColor int
	setupAdviceColor int
}

type LogCache struct {
	out     io.Writer
	err     io.Writer
	ErrLog  *bytes.Buffer
	command string
	message string
	seconds string
	colors  logColors
	verbose bool
	debug   bool
}

func NewLogCache(w io.Writer, e io.Writer) *LogCache {
	lc := LogCache{out: w, err: e}
	lc.ErrLog = new(bytes.Buffer)

	return &lc
}

func (lc *LogCache) SetCommandName(command string) {
	lc.command = command
}

func (lc *LogCache) EnableVerbose() {
	lc.verbose = true
}

func (lc *LogCache) Verbose() bool {
	return lc.verbose
}

func (lc *LogCache) EnableDebug() {
	lc.debug = true
}

// Printlnf prints a plain string (no color) and appends to log
func (lc *LogCache) Printlnf(format string, a ...interface{}) {
	lc.message += fmt.Sprintf(format+"\n", a...)
	fmt.Fprintf(lc.out, format+"\n", a...)
}

// Printf prints a plain string (no color) with no new line (mainly used instead of fmt.Printf())
func (lc *LogCache) Printf(format string, a ...interface{}) {
	lc.message += fmt.Sprintf(format, a...)
	fmt.Fprintf(lc.out, format, a...)
}

// Verbosef prints messages when `--verbose` is specified.  not logged
func (lc *LogCache) Verbosef(format string, a ...interface{}) {
	if lc.verbose {
		fmt.Fprintf(lc.out, format+"\n", a...)
	}
}

// Debugf prints messages when `--debug` is specififed.  not logged
func (lc *LogCache) Debugf(format string, a ...interface{}) {
	if lc.debug {
		fmt.Fprintf(lc.out, "🪲 "+format+"\n", a...)
	}
}

func (lc *LogCache) Errorf(format string, a ...interface{}) error {
	lc.message += fmt.Sprintf("Error: "+format+"\n", a...)
	return fmt.Errorf(color.New(color.Attribute(lc.colors.errorColor)).Sprintf(format, a...))
}
