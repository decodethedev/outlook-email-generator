package logrus

import "github.com/fatih/color"

var InfoAllowed bool = true

func Error(print string) {
	color.Red(print)
}

func Errorf(print string, args ...interface{}) {
	color.Red(print, args...)
}

func Info(print string) {
	if InfoAllowed {
		color.Cyan(print)
	}
}

func Infof(print string, args ...interface{}) {
	if InfoAllowed {
		color.Cyan(print, args...)
	}
}

func Print(print string) {
	color.Yellow(print)
}

func Printf(print string, args ...interface{}) {
	color.Yellow(print, args...)
}

func Success(print string) {
	color.Green(print)
}

func Successf(print string, args ...interface{}) {
	color.Green(print, args...)
}
