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

type flagPaser struct {
	fset     *flag.FlagSet
	file     string
	sec      int64
	min      int64
	hour     int64
	time     string
	routine  string
	loop     bool
	describe string
	verbose  bool
}

func newParser() *flagPaser {
	e := &flagPaser{fset: flag.NewFlagSet("goalarm", flag.ExitOnError)}
	e.fset.StringVar(&e.file, "file", "", "Path of sound file.")
	e.fset.Int64Var(&e.sec, "sec", 0, "Wait second.")
	e.fset.Int64Var(&e.min, "min", 0, "Wait minute.")
	e.fset.Int64Var(&e.hour, "hour", 0, "Wait hour.")
	e.fset.StringVar(&e.time, "time", "", "Call time.(15:00:01)")
	e.fset.StringVar(&e.routine, "routine", "", `Alarm routine. Format is json array. [{"range":20,"name":"working"},{"range":5,"name":"break"}]`)
	e.fset.BoolVar(&e.loop, "loop", false, "Loop Alarm.")
	e.fset.StringVar(&e.describe, "describe", "", "Describe command or status.")
	e.fset.BoolVar(&e.verbose, "v", false, "Ouput verbose.")
	return e
}

func (e *flagPaser) parse(args []string) error {
	return e.fset.Parse(args)
}

func main() {
	parser := newParser()
	if err := exec(parser, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		parser.fset.Usage()
		os.Exit(1)
	}
}

func exec(parser *flagPaser, args []string) error {
	err := parser.parse(args[1:])
	if err != nil {
		return err
	}
	log.SetVerbose(parser.verbose)
	// describe mode
	if parser.describe != "" {
		jw := json.NewEncoder(os.Stdout)
		jw.SetIndent("", "  ")
		switch parser.describe {
		case "command":
			err := jw.Encode(timeserver.AllCommands())
			if err != nil {
				return err
			}
		case "status":
			err := jw.Encode(timeserver.AllStatuses())
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s has no describe. support command or status.", parser.describe)
		}
		return nil
	}

	if parser.file == "" || parser.sec == 0 && parser.min == 0 && parser.hour == 0 && parser.time == "" && parser.routine == "" {
		return fmt.Errorf("insufficient arguments")
	}

	if _, err := os.Stat(parser.file); os.IsNotExist(err) {
		return err
	}

	// routine mode
	if parser.routine != "" {
		log.Printf("input routine: %s\n", parser.routine)
		var rj []taskJson
		err := json.Unmarshal([]byte(parser.routine), &rj)
		if err != nil {
			return fmt.Errorf("parse routine: %w", err)
		}
		if err := rtn.RunRoutine(os.Stdin, os.Stdout, convertTask(rj), parser.file, parser.loop); err != nil {
			return err
		}
		return nil
	}

	// alarm mode
	var duration time.Duration
	if parser.time != "" {
		log.Printf("input time: %s\n", parser.time)
		duration, err = timeParse(parser.time, time.Now())
		if err != nil {
			return fmt.Errorf("parse time arg: %w", err)
		}
	} else {
		log.Printf("input hour(%d), min(%d), sec(%d)\n", parser.hour, parser.min, parser.sec)
		duration = time.Hour*time.Duration(parser.hour) + time.Minute*time.Duration(parser.min) + time.Second*time.Duration(parser.sec)
	}

	return rtn.RunAlarm(os.Stdin, os.Stdout, duration, parser.file, parser.loop)
}
