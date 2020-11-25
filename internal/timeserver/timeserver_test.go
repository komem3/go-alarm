package timeserver_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/komem3/goalarm/internal/testutil"
	"github.com/komem3/goalarm/internal/timeserver"
)

type (
	given struct {
		task        timeserver.Task
		commandTime time.Time
		command     string
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "get",
			},
			commandTime: shortTime(1, 0, 5),
			command:     string(timeserver.GetCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
					Task: timeserver.Task{
						Index: 1,
						Range: time.Second * 10,
						Name:  "get",
					},
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Minute * 11,
				Name:  "pause",
			},
			commandTime: shortTime(1, 1, 0),
			command:     string(timeserver.PauseCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:   "10m0s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Hour * 11,
				Name:  "stop",
			},
			commandTime: shortTime(1, 1, 1),
			command:     string(timeserver.StopCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.StopStatus,
					Left:   "10h58m59s",
				},
			},
		},
	},
	{
		"restart command",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 5,
				Name:  "pause",
			},
			commandTime: shortTime(0, 0, 1),
			command:     string(timeserver.RestartCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	{
		"start command",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 5,
				Name:  "start",
			},
			commandTime: shortTime(0, 0, 1),
			command:     string(timeserver.StartCommand),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "unknown",
			},
			commandTime: shortTime(1, 0, 5),
			command:     "unknown",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "get get",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.GetCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "pause get",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.GetCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "restart get",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.GetCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "get stop",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.StopCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.StopStatus,
					Left:   "5s",
					Error:  nil,
				},
			},
		},
	},
	{
		"pause stop",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "pause stop",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.StopCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.StopStatus,
					Left:   "5s",
					Error:  nil,
				},
			},
		},
	},
	{
		"restart stop",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "restart stop",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.StopCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.StopStatus,
					Left:   "10s",
					Error:  nil,
				},
			},
		},
	},
	// pause
	{
		"get pause",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "get pause",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.PauseCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "pause pause",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.PauseCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
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
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "restart pause",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.PauseCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.PauseStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	// restart
	{
		"get restart",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "get restart",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.RestartCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
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
		"pause restart",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "pause restart",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.RestartCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
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
		"restart restart",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "restart pause",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.RestartCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.ErrorStatus,
					Error:  io.EOF,
				},
			},
		},
	},
	// start
	{
		"get start",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "get start",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.GetCommand,
				timeserver.StartCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
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
		"pause start",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "pause start",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.PauseCommand,
				timeserver.StartCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.PauseStatus,
					Left:   "5s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "5s",
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
		"restart start",
		given{
			task: timeserver.Task{
				Index: 1,
				Range: time.Second * 10,
				Name:  "restart start",
			},
			commandTime: shortTime(1, 0, 5),
			command: fmt.Sprintf("%s\n%s",
				timeserver.RestartCommand,
				timeserver.StartCommand,
			),
		},
		want{
			results: []timeserver.Result{
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
					Error:  nil,
				},
				{
					Status: timeserver.RunningStatus,
					Left:   "10s",
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
			task: timeserver.Task{
				Index: 1,
				Range: 0,
				Name:  "alarm finish",
			},
			commandTime: shortTime(1, 0, 0),
			command:     string(timeserver.GetCommand),
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
					tserver := timeserver.NewTimeServer(tt.given.task)
					var results []timeserver.Result
					{
						tserver.SetNow(now)
						tserver.HandlerFunc(func(r timeserver.Result) {
							results = append(results, r)
						})
						tserver.StartTimer()
						tserver.SetNow(tt.given.commandTime)
					}

					lastResult := tserver.Listen(testutil.MockIn(fmt.Sprintf("%s\n", tt.given.command)))
					if diff := cmp.Diff(len(results), len(tt.want.results)); diff != "" {
						t.Fatalf("result length: given(-), want(+)\n%s\n", diff)
					}

					for i, r := range results {
						tt.want.results[i].Task = tt.given.task
						if diff := cmp.Diff(r.Error, tt.want.results[i].Error, cmpopts.EquateErrors()); diff != "" {
							t.Errorf("result.Error: given(-), want(+)\n%s\n", diff)
						}

						if i == len(results)-1 {
							if diff := cmp.Diff(lastResult.Error, tt.want.results[i].Error, cmpopts.EquateErrors()); diff != "" {
								t.Errorf("latestResult.Error: given(-), want(+)\n%s\n", diff)
							}
							if diff := cmp.Diff(lastResult, tt.want.results[i], cmpopts.IgnoreFields(timeserver.Result{}, "Error")); diff != "" {
								t.Errorf("result: given(-), want(+)\n%s\n", diff)
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
