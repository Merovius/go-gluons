package log

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/prasannavl/go-gluons/ansicode"
)

type ColorStringer interface {
	ColorString() string
}

func DefaultTextFormatter(r *Record) string {
	var buf bytes.Buffer
	buf.WriteString(r.Meta.Time.Format(time.RFC3339))
	buf.WriteString("," + LogLevelString(r.Meta.Level) + ",")
	args := r.Args
	if r.Format == "" {
		buf.WriteString(strconv.Quote(fmt.Sprint(args...)))
	} else if len(args) > 0 {
		buf.WriteString(strconv.Quote(fmt.Sprintf(r.Format, args...)))
	} else {
		buf.WriteString(strconv.Quote(r.Format))
	}
	fields := GetFields(r.Meta.Logger)
	for _, x := range fields {
		buf.WriteString("," + strconv.Quote(x.Name) + "=" + strconv.Quote(fmt.Sprint(x.Value)))
	}
	buf.WriteString("\r\n")
	return buf.String()
}

var initTime = time.Now()

func DefaultTextFormatterForHuman(r *Record) string {
	var buf bytes.Buffer
	var timeFormat string
	t := r.Meta.Time
	if t.Sub(initTime).Hours() > 24 {
		timeFormat = "03:04:05 (Jan 02)"
	} else {
		timeFormat = "03:04:05"
	}
	buf.WriteString(t.Format(timeFormat))
	buf.WriteString("  " + PaddedString(LogLevelString(r.Meta.Level), 5) + "  ")
	args := r.Args
	if r.Format == "" {
		fmt.Fprint(&buf, args...)
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, r.Format, args...)
	} else {
		buf.WriteString(r.Format)
	}
	for _, x := range GetFields(r.Meta.Logger) {
		fmt.Fprintf(&buf, " %s=%v ", x.Name, x.Value)
	}
	buf.WriteString("\r\n")
	return buf.String()
}

func DefaultColorTextFormatterForHuman(r *Record) string {
	var buf bytes.Buffer
	t := r.Meta.Time
	var timeFormat string
	if t.Sub(initTime).Hours() > 24 {
		timeFormat = "03:04:05 (Jan 02)"
	} else {
		timeFormat = "03:04:05"
	}
	buf.WriteString(ansicode.BlackBright + t.Format(timeFormat) + ansicode.Reset)
	buf.WriteString("  " + logLevelColoredString(r.Meta.Level) + "  ")
	args := r.Args
	for i, a := range args {
		if colorable, ok := a.(ColorStringer); ok {
			args[i] = colorable.ColorString()
		}
	}
	if r.Format == "" {
		fmt.Fprint(&buf, args...)
	} else if len(args) > 0 {
		fmt.Fprintf(&buf, r.Format, args...)
	} else {
		buf.WriteString(r.Format)
	}
	for _, x := range GetFields(r.Meta.Logger) {
		value := x.Value
		if colorable, ok := value.(ColorStringer); ok {
			value = colorable.ColorString()
		}
		fmt.Fprintf(&buf, " %s=%v ", HashColoredText(x.Name), value)
	}
	buf.WriteString("\r\n")
	return buf.String()
}

var colorMap = []string{
	ansicode.BlackBright,
	ansicode.Cyan,
	ansicode.Green,
	ansicode.Magenta,
}

var colorMapLen = len(colorMap)

func HashColoredText(name string) string {
	const maxIterations = 10
	l := len(name)
	if l > 10 {
		l = 10
	}
	for i, x := range name {
		if i > maxIterations {
			break
		}
		l += int(x)
	}
	index := l % colorMapLen
	if index < 0 {
		index = 0
	}
	return colorMap[index] + name + ansicode.Reset
}

func logLevelColoredString(lvl Level) string {
	return LogLevelColoredMsg(lvl, PaddedString(LogLevelString(lvl), 5))
}

func LogLevelColoredMsg(lvl Level, msg string) string {
	switch lvl {
	case ErrorLevel:
		return ansicode.RedBright + msg + ansicode.Reset
	case WarnLevel:
		return ansicode.YellowBright + msg + ansicode.Reset
	case InfoLevel:
		return ansicode.Blue + msg + ansicode.Reset
	case DebugLevel:
		return ansicode.White + msg + ansicode.Reset
	case TraceLevel:
		return ansicode.BlackBright + msg + ansicode.Reset
	}
	return ansicode.BlackBright + msg + ansicode.Reset
}
