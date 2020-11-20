package main

import (
	"strconv"
	"strings"
	"time"
)

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
