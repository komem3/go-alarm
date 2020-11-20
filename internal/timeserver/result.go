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

var _ json.Marshaler = Result{}

func (r Result) MarshalJSON() ([]byte, error) {
	w := new(bytes.Buffer)
	jw := json.NewEncoder(w)
	_, err := w.WriteRune('{')
	if err != nil {
		return nil, err
	}

	{
		_, err = w.WriteString("\"status\":")
		if err != nil {
			return nil, err
		}
		err = jw.Encode(r.Status)
		if err != nil {
			return nil, err
		}
	}
	{
		_, err = w.WriteString(",\"left\":")
		if err != nil {
			return nil, err
		}
		err = jw.Encode(r.Left)
		if err != nil {
			return nil, err
		}
	}
	{
		_, err = w.WriteString(",\"error\":")
		if err != nil {
			return nil, err
		}
		if r.Error == nil {
			err = jw.Encode("")
			if err != nil {
				return nil, err
			}
		} else {
			err = jw.Encode(r.Error.Error())
			if err != nil {
				return nil, err
			}
		}
	}
	{
		_, err = w.WriteString(",\"task\":{\"index\":")
		if err != nil {
			return nil, err
		}
		{
			err = jw.Encode(r.Task.Index)
			if err != nil {
				return nil, err
			}
			_, err = w.WriteString(
				fmt.Sprintf(",\"range\":\"%s\"", r.Task.Range.Round(time.Second)))
			if err != nil {
				return nil, err
			}
			_, err = w.WriteString(
				fmt.Sprintf(",\"name\":\"%s\"}", r.Task.Name))
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = w.WriteRune('}')
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
