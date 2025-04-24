package storage

import (
	"testing"

	"github.com/kudras3r/KDTog/pkg/logger"
)

func TestFStorage(t *testing.T) {
	l := logger.New("DEBUG")
	s, err := NewFStorage(l, "/home/kud/WorkSpace/go/src/KDTog/tst/")
	if err != nil {
		t.Fatalf("error when init FSt %v", err)
	}
	_, err = s.GetIDByName("bname1")
	if err != nil {
		t.Fatalf("error when getIDByName: %v", err)
	}
}
