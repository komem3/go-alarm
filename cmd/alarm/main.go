package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/komem3/go-alarm"
)

func main() {
	file := flag.String("file", "bell.mp3", "Path of sound file.")
	sec := flag.Int64("sec", 0, "Wait second.")
	min := flag.Int64("min", 0, "Wait minute.")
	hour := flag.Int64("hour", 0, "Wait hour.")
	tim := flag.String("time", "", "Call time.(15:00:01)")
	flag.Parse()

	if *sec == 0 && *min == 0 && *hour == 0 && *tim == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *tim != "" {
		alarmD, err := timeParse(*tim)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse time arg: %v\n", err)
			flag.Usage()
			os.Exit(1)
		}
		alarm.Timer(alarmD)
	} else {
		alarm.Timer(time.Hour*time.Duration(*hour) + time.Minute*time.Duration(*min) + time.Second*time.Duration(*sec))
	}
	if err := alarm.Alarm(*file); err != nil {
		fmt.Fprintf(os.Stderr, "alarm : %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
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
