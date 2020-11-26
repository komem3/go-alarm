package sound

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	PlayWait()
}

var ErrUnsuportExt = fmt.Errorf("unsuported ext")

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
		err = fmt.Errorf("open %s: %w", filepath.Ext(path), ErrUnsuportExt)
	}
	if err != nil {
		return nil, err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	return &Alarm{
		buffer: buffer,
	}, nil
}

func (a *Alarm) Play() {
	log.Printf("async play sound\n")
	speaker.Play(a.buffer.Streamer(0, a.buffer.Len()))
}

func (a *Alarm) PlayWait() {
	log.Printf("wait play sound\n")
	done := make(chan struct{})
	speaker.Play(beep.Seq(a.buffer.Streamer(0, a.buffer.Len()), beep.Callback(func() {
		done <- struct{}{}
	})))
	<-done
}
