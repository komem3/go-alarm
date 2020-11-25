package routine

import "github.com/komem3/goalarm/internal/testutil"

func SetMock() {
	newAlarm = testutil.NewMockAlarm
}
