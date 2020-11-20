package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	rtn "github.com/komem3/goalarm/internal/routine"
	"github.com/komem3/goalarm/internal/timeserver"
)

var (
	file     string
	sec      int64
	min      int64
	hour     int64
	tim      string
	routine  string
	describe bool
)

func init() {
	flag.StringVar(&file, "file", "", "Path of sound file.")
	flag.Int64Var(&sec, "sec", 0, "Wait second.")
	flag.Int64Var(&min, "min", 0, "Wait minute.")
	flag.Int64Var(&hour, "hour", 0, "Wait hour.")
	flag.StringVar(&tim, "time", "", "Call time.(15:00:01)")
	flag.StringVar(&routine, "routine", "", `Alarm routine. Format is json array. [{"range":20,"name":"working"},{"range":5,"name":"break"}]`)
	flag.BoolVar(&describe, "describe", false, "Describe command or status.")
	flag.Parse()
}

func main() {
	if describe {
		target := flag.Arg(0)
		jw := json.NewEncoder(os.Stdin)
		jw.SetIndent("", "  ")
		switch target {
		case "command":
			err := jw.Encode(timeserver.AllCommands())
			if err != nil {
				fatalf(err.Error())
			}
		case "status":
			err := jw.Encode(timeserver.AllStatuses())
			if err != nil {
				fatalf(err.Error())
			}
		default:
			fatalf("%s has no describe. support command or status.\n", target)
		}
		return
	}

	if file == "" || sec == 0 && min == 0 && hour == 0 && tim == "" && routine == "" {
		fatalf("")
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fatalf("%s does not exist\n", file)
	}

	if routine != "" {
		var rj []taskJson
		err := json.Unmarshal([]byte(routine), &rj)
		if err != nil {
			fmt.Printf("%+v\n", routine) // output for debug

			fatalf("parse routine: %v\n", err)
		}
		if err := rtn.RunRoutine(os.Stdin, os.Stdout, convertTask(rj), file); err != nil {
			fatalf(err.Error())
		}
		return
	}

	var (
		duration time.Duration
		err      error
	)
	if tim != "" {
		duration, err = timeParse(tim)
		if err != nil {
			fatalf("parse time arg: %v\n", err)
		}
	} else {
		duration = time.Hour*time.Duration(hour) + time.Minute*time.Duration(min) + time.Second*time.Duration(sec)
	}

	if err := rtn.RunAlarm(os.Stdin, os.Stdout, duration, file); err != nil {
		fatalf(err.Error())
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	flag.Usage()
	os.Exit(1)
}
