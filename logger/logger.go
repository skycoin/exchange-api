/*
package logger inplements simple logger with levels
by default, logger writes to STDOUT
if you set level, sets out for logger too
you may change logger out for each level by calling SetOut func
*/

package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Info same as fmt.Println
func Info(v ...interface{}) {
	if info != nil {
		fmt.Fprintln(info, format(All, v...))
	}
}

// Infof same as fmt.Printf
func Infof(f string, v ...interface{}) {
	if info != nil {
		fmt.Fprintln(info, formatf(All, f, v...))
	}
}

// Warning same as fmt.Println
func Warning(v ...interface{}) {
	if warn != nil {
		fmt.Fprintln(warn, format(Warnings, v...))
	}
}

// Warningf same as fmt.Println
func Warningf(f string, v ...interface{}) {
	if warn != nil {
		fmt.Fprintln(warn, formatf(Warnings, f, v...))
	}
}

// Error same as fmt.Println
func Error(v ...interface{}) {
	if err != nil {
		fmt.Fprintln(err, format(Errors, v...))
	}
}

// Errorf same as fmt.Printf
func Errorf(f string, v ...interface{}) {
	if err != nil {
		fmt.Fprintln(err, formatf(Errors, f, v...))
	}
}

// Fatal same as fmt.Println and panic after print
func Fatal(v ...interface{}) {
	if fatal != nil {
		fmt.Fprintln(fatal, format(Criticals, v...))
	}
	panic("logger.Fatalln")
}

// Fatalf as fmt.Printf and panic after print
func Fatalf(f string, v ...interface{}) {
	if fatal != nil {
		fmt.Fprintln(fatal, formatf(Criticals, f, v...))
	}
	panic("logger.Fatalf")
}

// SetOut sets outputs for given logger level
func SetOut(w io.Writer, lvl Level) {
	switch lvl {
	case All:
		info = w
	case Warnings:
		warn = w
	case Errors:
		err = w
	case Criticals:
		fatal = w
	}
}

// SetLogLevel sets output for lvl and all loglevels, there higher than given lvl
func SetLogLevel(lvl Level, w io.Writer) {
	info, warn, err, fatal = nil, nil, nil, nil
	for lvl > 0 {
		SetOut(w, lvl)
		lvl = lvl >> 1
	}
}

var (
	info  io.Writer
	warn  io.Writer
	err   io.Writer
	fatal io.Writer
)

// Level is a type for logger Level
type Level int

//Logger levels
const (
	Criticals Level = 1 << iota
	Errors
	Warnings
	All
)

func format(lvl Level, v ...interface{}) string {
	var strlvl string
	switch lvl {
	case All:
		strlvl = "Info "
	case Warnings:
		strlvl = "Warning "
	case Errors:
		strlvl = "Error "
	case Criticals:
		strlvl = "Fatal "
	}
	strlvl += time.Now().Format("15:04:05 01/02") + ": "
	return strlvl + fmt.Sprintln(v...)
}
func formatf(lvl Level, f string, v ...interface{}) string {
	var strlvl string
	switch lvl {
	case All:
		strlvl = "Info "
	case Warnings:
		strlvl = "Warning "
	case Errors:
		strlvl = "Error "
	case Criticals:
		strlvl = "Fatal "
	}
	strlvl += time.Now().Format("15:04:05 01/02") + ": "
	f = strlvl + f
	return fmt.Sprintf(f, v...)
}

func init() {
	info = os.Stdout
	warn = os.Stdout
	err = os.Stdout
	fatal = os.Stdout
}
