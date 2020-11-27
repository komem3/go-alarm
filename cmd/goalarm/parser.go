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

type timeParser struct {
	hour int
	min  int
	sec  int
	err  error
}

func (t *timeParser) setHour(h string) *timeParser {
	if t.err != nil {
		return t
	}
	t.hour, t.err = strconv.Atoi(h)
	return t
}

func (t *timeParser) setMin(m string) *timeParser {
	if t.err != nil {
		return t
	}
	t.min, t.err = strconv.Atoi(m)
	return t
}

func (t *timeParser) setSec(s string) *timeParser {
	if t.err != nil {
		return t
	}
	t.sec, t.err = strconv.Atoi(s)
	return t
}

func (t *timeParser) time(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), t.hour, t.min, t.sec, 0, time.Local)
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

func timeParse(tstr string, now time.Time) (d time.Duration, err error) {
	times := strings.Split(tstr, ":")
	parser := &timeParser{
		hour: now.Hour(),
		min:  now.Minute(),
		sec:  now.Second(),
	}
	switch len(times) {
	case 3:
		parser.setHour(times[0]).setMin(times[1]).setSec(times[2])
	case 2:
		parser.setHour(times[0]).setMin(times[1])
	case 1:
		parser.setHour(times[0])
	}

	if parser.err != nil {
		return 0, parser.err
	}

	return parser.time(now).Sub(now), parser.err
}
