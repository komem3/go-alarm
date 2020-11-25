package sound_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/komem3/goalarm/internal/sound"
)

func TestAlarm(t *testing.T) {
	tests := []struct {
		name    string
		given   string
		wantErr string
	}{
		{"mp3", "./sample.mp3", "mp3: EOF"},
		{"wav", "./sample.wav", "wav: EOF"},
		{"not suport format", "./sample.mp4", "unsuported ext .mp4"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := emptyFile(tt.given); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tt.given)
			_, err := sound.NewAalarm(tt.given)
			if diff := cmp.Diff(err.Error(), tt.wantErr); diff != "" {
				t.Errorf("Alarm error: given(-), want(+)\n%s\n", diff)
			}
		})
	}
}

func emptyFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}
