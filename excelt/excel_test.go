package excelt

import (
	"fmt"
	"testing"
)

func TestGetFileList(t *testing.T) {
	var listpath = "."
	//listpath, _ = os.Getwd()
	_, _ = fmt.Scanf("%s", &listpath)
	GetFileList(listpath)
	ListFileFunc(listfile)
}
