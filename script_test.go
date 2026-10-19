package ntfilepicker_test

import (
	"testing"

	"nimbleterminal.dev/ntfilepicker/internal/pickercfg"
	"nimbleterminal.dev/ntfilepicker/internal/pickertest"
)

func TestScripts(t *testing.T) {
	pickertest.RunDir(t, pickercfg.Config(), "testdata/scripts")
}
