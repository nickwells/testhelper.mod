package testhelper_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestGoldenFilePathName(t *testing.T) {
	testCases := []struct {
		name         string
		gfc          testhelper.GoldenFileCfg
		fname        string
		expectedPath string
	}{
		{
			name: "has all parts",
			gfc: testhelper.GoldenFileCfg{
				DirNames: []string{"testdata", "subdir"},
				Pfx:      "PFX",
				Sfx:      "SFX",
			},
			fname: "file",
			expectedPath: "testdata" + string(filepath.Separator) +
				"subdir" + string(filepath.Separator) +
				"PFX.file.SFX",
		},
		{
			name: "no prefix",
			gfc: testhelper.GoldenFileCfg{
				DirNames: []string{"testdata", "subdir"},
				Sfx:      "SFX",
			},
			fname: "file",
			expectedPath: "testdata" + string(filepath.Separator) +
				"subdir" + string(filepath.Separator) +
				"file.SFX",
		},
		{
			name: "no suffix",
			gfc: testhelper.GoldenFileCfg{
				DirNames: []string{"testdata", "subdir"},
				Pfx:      "PFX",
			},
			fname: "file",
			expectedPath: "testdata" + string(filepath.Separator) +
				"subdir" + string(filepath.Separator) +
				"PFX.file",
		},
		{
			name: "no prefix or suffix",
			gfc: testhelper.GoldenFileCfg{
				DirNames: []string{"testdata", "subdir"},
			},
			fname: "file",
			expectedPath: "testdata" + string(filepath.Separator) +
				"subdir" + string(filepath.Separator) +
				"file",
		},
		{
			name: "bad filename (has embedded /)",
			gfc: testhelper.GoldenFileCfg{
				DirNames: []string{"testdata", "subdir"},
			},
			fname: ".." + string(filepath.Separator) + "file",
			expectedPath: "testdata" + string(filepath.Separator) +
				"subdir" + string(filepath.Separator) +
				"file",
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		path := tc.gfc.PathName(tc.fname)
		if path != tc.expectedPath {
			t.Log(tcID)
			t.Logf("\t: expected path: %s\n", tc.expectedPath)
			t.Logf("\t:           got: %s\n", path)
			t.Errorf("\t: Unexpected pathname\n")
		}
	}
}
