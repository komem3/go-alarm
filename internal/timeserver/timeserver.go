package timeserver

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

type timeServer struct {
	running  bool
	ticker   *time.Timer
	start    time.Time
	duration time.Duration
	now      func() time.Time
	handler  Handler
	reader   io.Reader
}

type Handler interface {
	Serve(r Result)
}

type handlerFunc func(r Result)

func (h handlerFunc) Serve(r Result) {
	h(r)
}

func NewTimeServer(d time.Duration) *timeServer {
	tserver := &timeServer{
		duration: d,
		now:      time.Now,
		running:  false,
	}
	return tserver
}

func (t *timeServer) StartTimer() {
	t.start = t.now()
	t.ticker = time.NewTimer(t.duration)
}

func (t *timeServer) HandlerFunc(f func(r Result)) {
	t.handler = handlerFunc(f)
}

func (t *timeServer) Listen(in io.Reader) (result Result) {
	t.running = true
	t.reader = in
	defer t.ticker.Stop()
	return t.finishRace(t.readCommand, t.finish)
}

func (t *timeServer) readCommand() (result Result) {
	status := RunningStatus
	var pauseLeft time.Duration
	buf := bufio.NewReader(t.reader)

	for t.running {
		line, err := buf.ReadString('\n')
		if err != nil {
			result = Result{
				Status: ErrorStatus,
				Error:  err,
			}
			t.handler.Serve(result)
			t.running = false
		}

		// when finished while reading input
		if !t.running {
			break
		}

		var left time.Duration
		if status == PauseStatus {
			left = pauseLeft
		} else {
			left = t.duration - t.now().Sub(t.start)
		}
		leftSec := fmt.Sprintf("%s", left.Round(time.Second))

		switch Command(line[:len(line)-1]) {
		case GetCommand:
			result = Result{
				Left:   leftSec,
				Status: status,
			}
		case StartCommand:
			status = RunningStatus
			t.start = t.now()
			t.ticker.Stop()
			t.ticker.Reset(left)
			result = Result{
				Left:   leftSec,
				Status: status,
			}
		case PauseCommand:
			status = PauseStatus
			t.ticker.Stop()
			pauseLeft = left
			result = Result{
				Left:   leftSec,
				Status: status,
			}
		case StopCommand:
			t.ticker.Stop()
			result = Result{
				Left:   leftSec,
				Status: StopStatus,
			}
			t.running = false
		case RestartCommand:
			status = RunningStatus
			t.start = t.now()
			t.ticker.Stop()
			t.ticker.Reset(t.duration)
			result = Result{
				Left:   fmt.Sprintf("%s", t.duration.Round(time.Second)),
				Status: status,
			}
		default:
			err = fmt.Errorf("'%s' is %w", line[:len(line)-1], ErrUnknownCommand)
			result = Result{
				Status: ErrorStatus,
				Error:  err,
			}
			t.running = false
		}
		t.handler.Serve(result)
	}
	return result
}

func (t *timeServer) finish() Result {
	<-t.ticker.C
	t.running = false
	r := Result{
		Status: FinishStatus,
	}
	t.handler.Serve(r)
	return r
}

func (t *timeServer) finishRace(fs ...func() Result) Result {
	results := make(chan Result, len(fs))
	for _, f := range fs {
		go func(f func() Result) {
			results <- f()
		}(f)
	}
	return <-results
}
