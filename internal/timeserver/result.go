package timeserver

import (
	"bytes"
	"encoding/json"
)

type Result struct {
	Status Status
	Left   string
	Error  error
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
	_, err = w.WriteRune('}')
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
