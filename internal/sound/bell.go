package sound

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/komem3/goalarm/internal/log"
)

type Alarm struct {
	buffer *beep.Buffer
}

type Player interface {
	Play()
}

func NewAalarm(path string) (Player, error) {
	log.Printf("sound file is %s\n", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
	)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	default:
		err = fmt.Errorf("unsuported ext %s", filepath.Ext(path))
	}
	if err != nil {
		return nil, err
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	return &Alarm{
		buffer: buffer,
	}, nil
}

func (a *Alarm) Play() {
	speaker.Play(a.buffer.Streamer(0, a.buffer.Len()))
}
