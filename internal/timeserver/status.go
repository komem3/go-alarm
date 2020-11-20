package timeserver

type Status string

const (
	RunningStatus Status = "running"
	PauseStatus   Status = "pause"
	StopStatus    Status = "stop"
	FinishStatus  Status = "finish"
	ErrorStatus   Status = "error"
)

type StatusDescribe struct {
	Status   Status
	Describe string
}

func AllStatuses() []StatusDescribe {
	return []StatusDescribe{
		{RunningStatus, "Running timer."},
		{PauseStatus, "Pause timer."},
		{StopStatus, "Stopped timer."},
		{FinishStatus, "Finish timer."},
		{ErrorStatus, "Error has occurred."},
	}
}
