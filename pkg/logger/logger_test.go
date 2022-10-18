package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestSetLogger(t *testing.T) {
	t.Run("Logging", func(t *testing.T) {
		SetLogger(true)
		if !Log.DebugMode {
			t.Errorf("Expected debug mode to be true, got false")
		}
		SetLogger(false)
		if Log.DebugMode {
			t.Errorf("Expected debug mode to be false, got true")
		}
	})
}

func TestLogger(t *testing.T) {
	t.Run("Logging", func(t *testing.T) {
		SetLogger(true)
		Log.Output = &bytes.Buffer{}
		Infof("Info")
		Debugf("Debug")
		Warnf("Warn")
		output := Log.Output.(*bytes.Buffer).String()
		if !strings.Contains(output, "Info") {
			t.Errorf("Expected Info, got %s", output)
		}
		if !strings.Contains(output, "Debug") {
			t.Errorf("Expected Debug, got %s", output)
		}
		if !strings.Contains(output, "Warn") {
			t.Errorf("Expected Warn, got %s", output)
		}

	})
}
