package tempdir_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/goark/toolbox/tempdir"
)

func TestNew(t *testing.T) {
	td := tempdir.New("")
	if td.String() != os.TempDir() {
		t.Errorf("TempDir.BaseDir() is \"%v\", want \"%v\"", td, os.TempDir())
	}
	if td.Path() != "" {
		t.Errorf("TempDir.Path() is \"%v\", want \"\"", td.Path())
	}
}

func TestMakeDir(t *testing.T) {
	td := tempdir.New("")
	err := td.MakeDir()
	if err != nil {
		t.Errorf("TempDir.MakeDir() is \"%+v\", want nil", err)
	}
	fmt.Println("temporary directory:", td.Path())
	_ = td.CleanUp()
}
