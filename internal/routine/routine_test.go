package routine_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/komem3/goalarm/internal/routine"
	"github.com/komem3/goalarm/internal/testutil"
	"github.com/komem3/goalarm/internal/timeserver"
)

func TestRunRoutine(t *testing.T) {
	type given struct {
		r   routine.Routine
		cmd string
	}
	tests := []struct {
		name    string
		given   given
		wantErr error
	}{
		{
			"stop",
			given{
				r: routine.Routine{
					{
						Index: 2,
						Range: time.Second * 10,
						Name:  "second",
					},
					{
						Index: 1,
						Range: 0,
						Name:  "first",
					},
				},
				cmd: "stop\n",
			},
			nil,
		},
		{
			"finish",
			given{
				r: routine.Routine{
					{
						Index: 2,
						Range: 0,
						Name:  "second",
					},
					{
						Index: 1,
						Range: 0,
						Name:  "first",
					},
				},
				cmd: "get\n",
			},
			nil,
		},
		{
			"bad command error",
			given{
				r: routine.Routine{
					{
						Index: 1,
						Range: time.Second * 100,
						Name:  "first",
					},
					{
						Index: 2,
						Range: time.Second * 10,
						Name:  "second",
					},
				},
				cmd: "unknown\n",
			},
			timeserver.ErrUnknownCommand,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := routine.RunRoutine(
				testutil.MockIn(tt.given.cmd),
				ioutil.Discard,
				tt.given.r,
				"dummy",
				false,
			)
			if diff := cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("routine.RunRoutine error: given(-), want(+)\n%s\n", diff)
			}
		})
	}
}

func TestRunAlarm(t *testing.T) {
	type given struct {
		command  string
		duration time.Duration
	}
	tests := []struct {
		name    string
		given   given
		wantErr error
	}{
		{"stop", given{"stop\n", time.Second}, nil},
		{"finish", given{"get\n", 0}, nil},
		{"bad command error", given{"unknown\n", time.Second}, timeserver.ErrUnknownCommand},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := routine.RunAlarm(testutil.MockIn(tt.given.command), ioutil.Discard, tt.given.duration, "dummy", false)
			if diff := cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("routine.RunAlarm error: given(-), want(+)\n%s\n", diff)
			}
		})
	}
}
