package main

import (
	"strconv"
	"strings"
	"time"

	rtn "github.com/komem3/goalarm/internal/routine"
	"github.com/komem3/goalarm/internal/timeserver"
)

type taskJson struct {
	Index int           `json:"index"`
	Range time.Duration `json:"range"`
	Name  string        `json:"name"`
}

func convertTask(tasks []taskJson) rtn.Routine {
	r := make(rtn.Routine, 0, len(tasks))
	for _, t := range tasks {
		r = append(r, timeserver.Task{
			Index: t.Index,
			Range: time.Minute * t.Range,
			Name:  t.Name,
		})
	}
	return r
}

func timeParse(tstr string) (d time.Duration, err error) {
	now := time.Now()
	times := strings.Split(tstr, ":")
	var alarmTime time.Time
	switch len(times) {
	case 3:
		hour, err := strconv.Atoi(times[0])
		if err != nil {
			return 0, err
		}
		min, err := strconv.Atoi(times[1])
		if err != nil {
			return 0, err
		}
		sec, err := strconv.Atoi(times[2])
		if err != nil {
			return 0, err
		}
		alarmTime = time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, time.Local)
	case 2:
		hour, err := strconv.Atoi(times[0])
		if err != nil {
			return 0, err
		}
		min, err := strconv.Atoi(times[1])
		if err != nil {
			return 0, err
		}
		alarmTime = time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)
	case 1:
		hour, err := strconv.Atoi(times[0])
		if err != nil {
			return 0, err
		}
		alarmTime = time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.Local)
	}

	d = time.Until(alarmTime)
	return
}
