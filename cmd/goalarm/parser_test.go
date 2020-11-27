package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTimeParse(t *testing.T) {
	type want struct {
		d   time.Duration
		err error
	}
	tests := []struct {
		name  string
		given string
		want  want
	}{
		{
			"hour only",
			"16",
			want{time.Hour, nil},
		},
		{
			"hour + minute",
			"15:01",
			want{time.Minute, nil},
		},
		{
			"hour + minute + sec",
			"15:00:01",
			want{time.Second, nil},
		},
		{
			"bad hour",
			"bad",
			want{0, strconv.ErrSyntax},
		},
	}
	now := time.Date(2010, 1, 1, 15, 0, 0, 0, time.Local)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d, err := timeParse(tt.given, now)
			if diff := cmp.Diff(d, tt.want.d); diff != "" {
				t.Errorf("timeParse duration, given(+), want(-)\n%s\n", diff)
			}
			if diff := cmp.Diff(err, tt.want.err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("timeParse error, given(+), want(-)\n%s\n", diff)
			}
		})
	}
}
