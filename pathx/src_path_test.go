package pathx

import (
	"os"
	"path"
	"testing"
)

func TestModuleSrcPath(t *testing.T) {
	curDir, _ := os.Getwd()
	tests := []struct {
		name             string
		wantGoModDirPath string
		wantErr          bool
	}{
		{
			name:             "should get module src path",
			wantGoModDirPath: path.Dir(curDir),
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGoModDirPath, err := ModuleSrcPath(curDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModuleSrcPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotGoModDirPath != tt.wantGoModDirPath {
				t.Errorf("ModuleSrcPath() gotGoModDirPath = %v, want %v", gotGoModDirPath, tt.wantGoModDirPath)
			}
		})
	}
}
