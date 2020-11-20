package timeserver_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/komem3/goalarm/internal/timeserver"
)

func TestResult_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		given timeserver.Result
		want  string
	}{
		{
			"normal response",
			timeserver.Result{
				Status: timeserver.RunningStatus,
				Left:   "10m4s",
				Task: timeserver.Task{
					Index: 1,
					Range: time.Second,
					Name:  "normal",
				},
			},
			`{"status":"running","left":"10m4s","error":"","task":{"index":1,"range":"1s","name":"normal"}}`,
		},
		{
			"error case",
			timeserver.Result{
				Status: timeserver.ErrorStatus,
				Error:  fmt.Errorf("err is %w", timeserver.ErrUnknownCommand),
				Task: timeserver.Task{
					Index: 2,
					Range: time.Hour,
					Name:  "second",
				},
			},
			`{"status":"error","left":"","error":"err is not support command","task":{"index":2,"range":"1h0m0s","name":"second"}}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b, err := json.Marshal(tt.given)
			if err != nil {
				t.Errorf("marshal error %v", err)
			}
			if diff := cmp.Diff(string(b), tt.want); diff != "" {
				t.Errorf("marshal output: given(-), want(+)\n%s\n", diff)
			}
		})
	}
}
