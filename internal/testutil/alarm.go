package testutil

import (
	"github.com/komem3/goalarm/internal/sound"
)

type MockAlarm struct{}

var _ sound.Player = (*MockAlarm)(nil)

func NewMockAlarm(_ string) (sound.Player, error) {
	return &MockAlarm{}, nil
}

func (m *MockAlarm) Play()     {}
func (m *MockAlarm) PlayWait() {}
