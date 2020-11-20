package timeserver

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

type Task struct {
	Index int
	Range time.Duration
	Name  string
}

type timeServer struct {
	running bool
	ticker  *time.Timer
	start   time.Time
	task    Task
	now     func() time.Time
	handler Handler
	reader  io.Reader
}

type Handler interface {
	Serve(r Result)
}

type handlerFunc func(r Result)

func (h handlerFunc) Serve(r Result) {
	h(r)
}

func NewTimeServer(task Task) *timeServer {
	tserver := &timeServer{
		task:    task,
		now:     time.Now,
		running: false,
	}
	return tserver
}

func (t *timeServer) StartTimer() {
	t.start = t.now()
	t.ticker = time.NewTimer(t.task.Range)
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
				Task:   t.task,
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
			left = t.task.Range - t.now().Sub(t.start)
		}
		leftSec := fmt.Sprintf("%s", left.Round(time.Second))

		switch Command(line[:len(line)-1]) {
		case GetCommand:
			result = Result{
				Left:   leftSec,
				Status: status,
				Task:   t.task,
			}
		case StartCommand:
			status = RunningStatus
			t.start = t.now()
			t.ticker.Stop()
			t.ticker.Reset(left)
			result = Result{
				Left:   leftSec,
				Status: status,
				Task:   t.task,
			}
		case PauseCommand:
			status = PauseStatus
			t.ticker.Stop()
			pauseLeft = left
			result = Result{
				Left:   leftSec,
				Status: status,
				Task:   t.task,
			}
		case StopCommand:
			t.ticker.Stop()
			result = Result{
				Left:   leftSec,
				Status: StopStatus,
				Task:   t.task,
			}
			t.running = false
		case RestartCommand:
			status = RunningStatus
			t.start = t.now()
			t.ticker.Stop()
			t.ticker.Reset(t.task.Range)
			result = Result{
				Left:   fmt.Sprintf("%s", t.task.Range.Round(time.Second)),
				Status: status,
				Task:   t.task,
			}
		default:
			err = fmt.Errorf("'%s' is %w", line[:len(line)-1], ErrUnknownCommand)
			result = Result{
				Status: ErrorStatus,
				Error:  err,
				Task:   t.task,
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
		Task:   t.task,
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
