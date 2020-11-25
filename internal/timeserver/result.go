package timeserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Result struct {
	Status Status
	Left   string
	Error  error
	Task   Task
}

type jsonWriter struct {
	b   *bytes.Buffer
	j   *json.Encoder
	err error
}

var _ json.Marshaler = Result{}

func (r Result) MarshalJSON() ([]byte, error) {
	w := new(bytes.Buffer)
	jw := &jsonWriter{
		b: w,
		j: json.NewEncoder(w),
	}

	jw.writeRune('{')
	jw.writeString("\"status\":").encode(r.Status)
	jw.writeString(",\"left\":").encode(r.Left)

	jw.writeString(",\"error\":")
	if r.Error == nil {
		jw.encode("")
	} else {
		jw.encode(r.Error.Error())
	}

	jw.writeString(",\"task\":{\"index\":").encode(r.Task.Index)
	jw.writeFormat(",\"range\":\"%s\"", r.Task.Range.Round(time.Second))
	jw.writeFormat(",\"name\":\"%s\"}", r.Task.Name)

	jw.writeRune('}')
	return jw.b.Bytes(), jw.err
}

func (j *jsonWriter) encode(i interface{}) *jsonWriter {
	if j.err != nil {
		return j
	}
	j.err = j.j.Encode(i)
	return j
}

func (j *jsonWriter) writeString(str string) *jsonWriter {
	if j.err != nil {
		return j
	}
	_, j.err = j.b.WriteString(str)
	return j
}

func (j *jsonWriter) writeFormat(format string, args ...interface{}) *jsonWriter {
	if j.err != nil {
		return j
	}
	_, j.err = j.b.WriteString(fmt.Sprintf(format, args...))
	return j
}

func (j *jsonWriter) writeRune(r rune) *jsonWriter {
	if j.err != nil {
		return j
	}
	_, j.err = j.b.WriteRune(r)
	return j
}
