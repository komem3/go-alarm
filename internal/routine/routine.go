package routine

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/komem3/goalarm/internal/log"
	"github.com/komem3/goalarm/internal/sound"
	"github.com/komem3/goalarm/internal/timeserver"
)

type Routine []timeserver.Task

var newAlarm = sound.NewAalarm

func RunRoutine(r io.Reader, w io.Writer, routine Routine, file string, loop bool) error {
	alarm, err := newAlarm(file)
	if err != nil {
		return err
	}
	jw := json.NewEncoder(w)
	sort.Slice(routine, func(i, j int) bool {
		return routine[i].Index < routine[j].Index
	})
	for l := true; l; l = loop {
		for i, task := range routine {
			task.Index = i + 1
			result, err := runTask(r, jw, task)
			if err != nil {
				return err
			}
			if result.Status == timeserver.StopStatus {
				return nil
			}
			if len(routine)-1 == i && !loop {
				alarm.PlayWait()
			} else {
				alarm.Play()
			}
		}
	}
	return nil
}

func RunAlarm(r io.Reader, w io.Writer, d time.Duration, file string, loop bool) error {
	alarm, err := newAlarm(file)
	if err != nil {
		return err
	}
	jw := json.NewEncoder(w)
	for l := true; l; l = loop {
		result, err := runTask(r, jw,
			timeserver.Task{
				Index: 0,
				Range: d,
				Name:  "alarm",
			})
		if err != nil {
			return err
		}
		if result.Status == timeserver.StopStatus {
			return nil
		}
		if loop {
			alarm.Play()
		} else {
			alarm.PlayWait()
		}
	}
	return nil
}

func runTask(
	r io.Reader,
	jw *json.Encoder,
	task timeserver.Task,
) (result timeserver.Result, err error) {
	log.Printf("run task %s: %s\n", task.Name, task.Range)
	tserver := timeserver.NewTimeServer(task)
	tserver.StartTimer()
	tserver.HandlerFunc(func(r timeserver.Result) {
		err := jw.Encode(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	})

	result = tserver.Listen(r)
	if result.Error != nil {
		return result, fmt.Errorf("server error : %w", result.Error)
	}

	if result.Status == timeserver.StopStatus {
		return result, nil
	}

	return result, nil
}
