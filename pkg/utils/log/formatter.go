/*
 * @Date: 2021.01.06 19:24
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 21:56
 */

package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoMsgColors - not apply colors to the msg
	NoMsgColors bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool
}

var defaultFormatter = &Formatter{
	HideKeys:         true,
	NoColors:         false,
	NoFieldsColors:   false,
	NoMsgColors:      true,
	ShowFullLevel:    false,
	NoUppercaseLevel: false,
	TrimMessages:     true,
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat

	if timestampFormat == "" {
		timestampFormat = "2006-01-02 15:04:05.000"
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(entry.Time.Local().Format(timestampFormat))

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if !f.NoColors {
		_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}

	b.WriteString("]")

	if !f.NoColors {
		b.WriteString("\x1b[0m")
	}

	b.WriteString(" ")

	if !f.NoFieldsColors {
		_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	if !f.NoMsgColors {
		_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	if f.TrimMessages {
		entry.Message = strings.TrimSpace(entry.Message)
	}

	// write message
	b.WriteString(entry.Message)

	if !f.NoMsgColors {
		b.WriteString("\x1b[0m")
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		_, _ = fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		_, _ = fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
