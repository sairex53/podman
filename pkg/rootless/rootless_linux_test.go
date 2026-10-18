//go:build linux && cgo

package rootless

import (
	"os"
	"testing"
)

func TestIsDefinitelyNotPauseProcess(t *testing.T) {
	tests := []struct {
		name        string
		isPause     bool
		checkErr    error
		processGone bool
		want        bool
	}{
		{name: "readable non-pause process", want: true},
		{name: "pause process", isPause: true},
		{name: "process no longer exists", checkErr: os.ErrNotExist, processGone: true, want: true},
		{name: "process is hidden but still alive", checkErr: os.ErrNotExist},
		{name: "environment is inaccessible", checkErr: os.ErrPermission},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDefinitelyNotPauseProcess(tt.isPause, tt.checkErr, tt.processGone); got != tt.want {
				t.Fatalf("isDefinitelyNotPauseProcess(%v, %v, %v) = %v, want %v", tt.isPause, tt.checkErr, tt.processGone, got, tt.want)
			}
		})
	}
}
