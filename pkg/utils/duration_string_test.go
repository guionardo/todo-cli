package utils

import (
	"testing"
	"time"
)

func TestDurationString(t *testing.T) {

	tests := []struct {
		name string
		arg  time.Duration
		want string
	}{
		{"1 year", time.Hour * 24 * 365, "1 year"},
		{"1 month", time.Hour * 24 * 30, "1 month"},
		{"1 hour", time.Hour + time.Minute*15, "1 hour"},
		{"15 minutes", time.Minute * 20, "15 minutes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationString(tt.arg); got != tt.want {
				t.Errorf("DurationString() = %v, want %v", got, tt.want)
			}
		})
	}
}
