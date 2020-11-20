package routine

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/komem3/goalarm/internal/sound"
	"github.com/komem3/goalarm/internal/timeserver"
)

type Routine []timeserver.Task

func RunRoutine(r io.Reader, w io.Writer, routine Routine, file string) error {
	jw := json.NewEncoder(w)
	sort.Slice(routine, func(i, j int) bool {
		return routine[i].Index < routine[j].Index
	})
	for i, task := range routine {
		task.Index = i + 1
		result, err := runTask(r, jw, task, file)
		if err != nil {
			return err
		}
		if result.Status == timeserver.StopStatus {
			return nil
		}
	}
	return nil
}

func RunAlarm(r io.Reader, w io.Writer, d time.Duration, file string) error {
	_, err := runTask(r, json.NewEncoder(w),
		timeserver.Task{
			Index: 0,
			Range: d,
			Name:  "alarm",
		}, file)
	return err
}

func runTask(
	r io.Reader,
	jw *json.Encoder,
	task timeserver.Task,
	file string,
) (result timeserver.Result, err error) {
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

	if err := sound.Alarm(file); err != nil {
		return result, fmt.Errorf("alarm : %w", err)
	}

	return result, nil
}
