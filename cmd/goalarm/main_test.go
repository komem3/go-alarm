package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/komem3/goalarm/internal/testutil"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		file    string
		wantErr string
	}{
		{
			name: "describe status",
			args: []string{"goalarm", "-describe", "status"},
		},
		{
			name: "describe command",
			args: []string{"goalarm", "-describe", "command"},
		},
		{
			name:    "unsport descibe",
			args:    []string{"goalarm", "-describe", "unsport"},
			wantErr: "unsport has no describe. support command or status.",
		},
		{
			name:    "file error",
			args:    []string{"goalarm", "-file", "empty", "-sec", "10"},
			wantErr: "stat empty: no such file or directory",
		},
		{
			name:    "unsport ext(alarm)",
			args:    []string{"goalarm", "-file", "empty.mp4", "-sec", "10"},
			file:    "empty.mp4",
			wantErr: "open .mp4: unsuported ext",
		},
		{
			name:    "unsport ext(time hh:mm:ss)",
			args:    []string{"goalarm", "-file", "empty.time", "-time", "15:00:00"},
			file:    "empty.time",
			wantErr: "open .time: unsuported ext",
		},
		{
			name:    "unsport ext(time hh:mm)",
			args:    []string{"goalarm", "-file", "empty.hhmm", "-time", "15:00"},
			file:    "empty.hhmm",
			wantErr: "open .hhmm: unsuported ext",
		},
		{
			name:    "unsport ext(time hh)",
			args:    []string{"goalarm", "-file", "empty.hh", "-time", "15"},
			file:    "empty.hh",
			wantErr: "open .hh: unsuported ext",
		},
		{
			name:    "bad time format",
			args:    []string{"goalarm", "-file", "badtime.mp3", "-time", "date::"},
			file:    "badtime.mp3",
			wantErr: `parse time arg: strconv.Atoi: parsing "date": invalid syntax`,
		},
		{
			name: "unsport ext(routine)",
			args: []string{"goalarm", "-file", "empty.png", "-routine",
				`[{"range":20,"name":"working"},{"range":5,"name":"break"}]`},
			file:    "empty.png",
			wantErr: "open .png: unsuported ext",
		},
		{
			name: "unmarshal routine",
			args: []string{"goalarm", "-file", "unmarshal.mp3", "-routine",
				`[{"range":20,"name":"working"},]`},
			file:    "unmarshal.mp3",
			wantErr: "parse routine: invalid character ']' looking for beginning of value",
		},
		{
			name:    "file empty",
			args:    []string{"goalarm", "-sec", "10"},
			wantErr: "insufficient arguments",
		},
		{
			name:    "sec emptry",
			args:    []string{"goalarm", "-file", "empty"},
			wantErr: "insufficient arguments",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.file != "" {
				if err := testutil.EmptyFile(tt.file); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tt.file)
			}
			parser := newParser()
			err := exec(parser, tt.args)
			if tt.wantErr == "" {
				if err != nil {
					t.Error(err)
				}
				return
			}
			if diff := cmp.Diff(err.Error(), tt.wantErr); diff != "" {
				t.Errorf("exec error given(-), want(+)\n%s\n", diff)
			}
		})
	}
}
