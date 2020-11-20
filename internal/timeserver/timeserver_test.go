package timeserver_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/komem3/goalarm/internal/timeserver"
)

type (
	given struct {
		waitDuration time.Duration
		commandTime  time.Time
		command      string
	}
	want struct {
		results []timeserver.Result
	}
	testcases []struct {
		name  string
		given given
		want  want
	}
)

var oneCommandTestCase = testcases{
	{
		"get command",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command:      string(timeserver.GetCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "5s",
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"pause command",
		given{
			waitDuration: time.Minute * 11,
			commandTime:  shortTime(1, 1, 0),
			command:      string(timeserver.PauseCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:    "10m0s",
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"stop command",
		given{
			waitDuration: time.Hour * 11,
			commandTime:  shortTime(1, 1, 1),
			command:      string(timeserver.StopCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.StopStatus,
					Left:    "10h58m59s",
				},
			},
		},
	},
	{
		"restart command",
		given{
			waitDuration: time.Second * 5,
			commandTime:  shortTime(0, 0, 1),
			command:      string(timeserver.RestartCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "5s",
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"unknown command",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command:      "unknown",
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.ErrorStatus,
					Error:  timeserver.ErrUnknownCommand,
				},
			},
		},
	},
}

var twoCommandTestCase = testcases{
	// get
	{
		"get get",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.GetCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"pause get",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.GetCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"restart get",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.GetCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "10s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:    "10s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	// stop
	{
		"get stop",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.StopCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.StopStatus,
					Left:    "5s",
					Error:  nil,
				},
			},
		},
	},
	{
		"pause stop",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.StopCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.StopStatus,
					Left:    "5s",
					Error:  nil,
				},
			},
		},
	},
	{
		"restart stop",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.StopCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "10s",
					Error:  nil,
				},
				{
					Status: timeserver.StopStatus,
					Left:    "10s",
					Error:  nil,
				},
			},
		},
	},
	// pause
	{
		"get pause",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.PauseCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"pause pause",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.PauseCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:    "5s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"restart pause",
		given{
			waitDuration: time.Second * 10,
			commandTime:  shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.PauseCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:    "10s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:    "10s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
}

var finishTestCase = testcases{
	{
		name: "alarm finish",
		given: given{
			waitDuration: 0,
			commandTime:  shortTime(1, 0, 0),
			command:      string(timeserver.GetCommand),
		},
		want: want{
			results: []timeserver.Result{
				{
					Status: timeserver.FinishStatus,
				},
			},
		},
	},
}

func TestTimeSever_Listen(t *testing.T) {
	cases := []struct {
		name  string
		tests testcases
	}{
		{
			"one command",
			oneCommandTestCase,
		},
		{
			"two command",
			twoCommandTestCase,
		},
		{
			"finish alarm",
			finishTestCase,
		},
	}
	now := shortTime(1, 0, 0)

	for _, c := range cases {
		tests := c.tests
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					tserver := timeserver.NewTimeServer(tt.given.waitDuration)
					var results []timeserver.Result
					{
						tserver.SetNow(now)
						tserver.HandlerFunc(func(r timeserver.Result) {
							results = append(results, r)
						})
						tserver.StartTimer()
						tserver.SetNow(tt.given.commandTime)
					}
					reader, err := mockIn(fmt.Sprintf("%s\n", tt.given.command))
					if err != nil {
						t.Fatal(err)
					}

					lastResult := tserver.Listen(reader)
					if diff := cmp.Diff(len(results), len(tt.want.results)); diff != "" {
						t.Fatalf("result length: given(-), want(+)\n%s\n", diff)
					}

					for i, r := range results {
						if diff := cmp.Diff(r.Error, tt.want.results[i].Error, cmpopts.EquateErrors()); diff != "" {
							t.Errorf("result.Error: given(-), want(+)\n%s\n", diff)
						}
						if diff := cmp.Diff(r.Status, tt.want.results[i].Status); diff != "" {
							t.Errorf("result.Status: given(-), want(+)\n%s\n", diff)
						}
						if diff := cmp.Diff(r.Left, tt.want.results[i].Left); diff != "" {
							t.Errorf("result.Sec: given(-), want(+)\n%s\n", diff)
						}
						if i == len(results)-1 {
							if diff := cmp.Diff(lastResult.Error, tt.want.results[i].Error, cmpopts.EquateErrors()); diff != "" {
								t.Errorf("latestResult.Error: given(-), want(+)\n%s\n", diff)
							}
							if diff := cmp.Diff(lastResult.Status, tt.want.results[i].Status); diff != "" {
								t.Errorf("latestResult.Status: given(-), want(+)\n%s\n", diff)
							}
							if diff := cmp.Diff(lastResult.Left, tt.want.results[i].Left); diff != "" {
								t.Errorf("latestResult.Sec: given(-), want(+)\n%s\n", diff)
							}
						}
					}
				})
			}
		})
	}
}

func shortTime(h, min, sec int) time.Time {
	return time.Date(2010, 1, 1, h, min, sec, 0, time.Local)
}

func mockIn(in string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	_, err := buf.WriteString(in)
	return buf, err
}
