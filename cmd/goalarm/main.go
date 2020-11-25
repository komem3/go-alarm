package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/komem3/goalarm/internal/log"
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
	loop     bool
	verbose  bool
	describe bool
)

func init() {
	flag.StringVar(&file, "file", "", "Path of sound file.")
	flag.Int64Var(&sec, "sec", 0, "Wait second.")
	flag.Int64Var(&min, "min", 0, "Wait minute.")
	flag.Int64Var(&hour, "hour", 0, "Wait hour.")
	flag.StringVar(&tim, "time", "", "Call time.(15:00:01)")
	flag.StringVar(&routine, "routine", "", `Alarm routine. Format is json array. [{"range":20,"name":"working"},{"range":5,"name":"break"}]`)
	flag.BoolVar(&loop, "loop", false, "Loop Alarm.")
	flag.BoolVar(&describe, "describe", false, "Describe command or status.")
	flag.BoolVar(&verbose, "v", false, "Ouput verbose.")
	flag.Parse()
}

func main() {
	log.SetVerbose(verbose)

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
		log.Printf("input routine: %s\n", routine)
		var rj []taskJson
		err := json.Unmarshal([]byte(routine), &rj)
		if err != nil {
			fmt.Printf("%+v\n", routine) // output for debug

			fatalf("parse routine: %v\n", err)
		}
		if err := rtn.RunRoutine(os.Stdin, os.Stdout, convertTask(rj), file, loop); err != nil {
			fatalf(err.Error())
		}
		return
	}

	var (
		duration time.Duration
		err      error
	)
	if tim != "" {
		log.Printf("input time: %s\n", tim)
		duration, err = timeParse(tim)
		if err != nil {
			fatalf("parse time arg: %v\n", err)
		}
	} else {
		log.Printf("input hour(%d), min(%d), sec(%d)\n", hour, min, sec)
		duration = time.Hour*time.Duration(hour) + time.Minute*time.Duration(min) + time.Second*time.Duration(sec)
	}

	if err := rtn.RunAlarm(os.Stdin, os.Stdout, duration, file, loop); err != nil {
		fatalf("%v\n", err)
	}

}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	flag.Usage()
	os.Exit(1)
}
