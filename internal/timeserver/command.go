package timeserver

import "errors"

const (
	GetCommand     Command = "get"
	StartCommand   Command = "start"
	PauseCommand   Command = "pause"
	StopCommand    Command = "stop"
	RestartCommand Command = "restart"
)

var ErrUnknownCommand = errors.New("not support command")

type Command string

type CommandDescribe struct {
	Command  Command
	Describe string
}

func AllCommands() []CommandDescribe {
	return []CommandDescribe{
		{GetCommand, "Get status and left time."},
		{StartCommand, "Start timer when pause status."},
		{PauseCommand, "Pause timer."},
		{StopCommand, "Stop timer. This command stop process."},
		{RestartCommand, "Restart timer at the first."},
	}
}
