package logging

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

const jsonTimeFormat string = "2006-1-2|15:04:05.000"

//StdOutFormatter applies consistent formatting to every message
type StdOutFormatter struct {
	Name string
}

//Format satisfies logrus' Format interface while staying flexible
func (f *StdOutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	genericFormatter := logrus.TextFormatter{PadLevelText: true, TimestampFormat: "1-2|15:04:05.000", FullTimestamp: true}
	formatted, err := genericFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	label := fmt.Sprintf("%-10s ", f.Name)

	line := bytes.Join(
		[][]byte{[]byte(label), formatted},
		[]byte(" "),
	)

	return line, nil
}

//FileOutFormatter applies consistent formatting to every message
type FileOutFormatter struct {
	Name string
}

//Format satisfies logrus' Format interface while staying flexible
func (f *FileOutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	name, err := MarshalJsonNonHTML(f.Name)
	if err != nil {
		return nil, err
	}
	name = append(name, ',')

	level, err := MarshalJsonNonHTML(entry.Level)
	if err != nil {
		return nil, err
	}
	level = append(level, ',')

	time, err := MarshalJsonNonHTML(entry.Time.Format(jsonTimeFormat))
	if err != nil {
		return nil, err
	}
	time = append(time, ',')

	msg, err := MarshalJsonNonHTML(entry.Message)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf(`{ "module":%-12s "level":%-10s "time":%s "msg":%s`, name, level, time, msg))

	for k, v := range entry.Data {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		switch k {
		case "module", "level", "time", "msg":
			k = "fields." + k
		}

		buf.WriteString(fmt.Sprintf(`, %q:%s`, k, b))
	}
	buf.Write([]byte(" }\n"))
	return buf.Bytes(), nil
}

func MarshalJsonNonHTML(x interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)
	err := e.Encode(x)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	return b[:len(b)-1], nil // remove trailing newline
}
