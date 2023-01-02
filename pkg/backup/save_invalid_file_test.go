package backup

import (
	"os"
	"path"
	"testing"
)

func TestMoveFileToBackup(t *testing.T) {
	tmpDir := t.TempDir()
	realFile := path.Join(tmpDir, "realFile.txt")
	os.WriteFile(realFile, []byte("Hello"), 0644)
	dirFile := path.Join(tmpDir, "dirFile")
	os.Mkdir(dirFile, 0755)

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"Real file", realFile, false},
		{"Dir file", dirFile, true},
		{"Non existent file", path.Join(tmpDir, "nonExistentFile.txt"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewFileName, err := MoveFileToBackup(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveFileToBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if _, err := os.Stat(gotNewFileName); os.IsNotExist(err) {
					t.Errorf("MoveFileToBackup() new file %v does not exists", gotNewFileName)
				}
			}
		})
	}
}
