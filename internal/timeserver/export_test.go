package timeserver

import "time"

func (t *timeServer) SetNow(tim time.Time) {
	t.now = func() time.Time { return tim }
}
