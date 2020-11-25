package routine_test

import (
	"os"
	"testing"

	"github.com/komem3/goalarm/internal/routine"
)

func TestMain(m *testing.M) {
	routine.SetMock()
	os.Exit(m.Run())
}
