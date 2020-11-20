package testutil

import (
	"bytes"
	"io"
)

func MockIn(in string) io.Reader {
	buf := new(bytes.Buffer)
	_, err := buf.WriteString(in)
	if err != nil {
		panic(err)
	}
	return buf
}
