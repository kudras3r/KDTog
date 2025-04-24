package storage

import (
	"fmt"
	"testing"

	"github.com/kudras3r/KDTog/pkg/logger"
)

func TestGetIDByName(t *testing.T) {
	l := logger.New("DEBUG")
	s, err := NewFStorage(l, "/home/kud/WorkSpace/go/src/KDTog/tst/")
	if err != nil {
		t.Fatalf("error when init FSt %v", err)
	}
	_, err = s.GetIDByName("bname1")
	if err != nil {
		t.Fatalf("error when getIDByName: %v", err)
	}

	names := []string{"aaame1", "abame1", "bname1", "bname1", "caame1", "cbame1", "notFound"}
	for _, n := range names {
		id, err := s.GetIDByName(n)
		if err != nil {
			if n == "notFound" {
				t.Logf("user %s not found as expected", n)
				continue
			}
			t.Fatalf("error when getIDByName: %v", err)
		}
		t.Logf("id for %s is %d", n, id)
	}
}

func TestGetPHashByID(t *testing.T) {
	l := logger.New("DEBUG")
	s, err := NewFStorage(l, "/home/kud/WorkSpace/go/src/KDTog/tst/")
	if err != nil {
		t.Fatalf("error when init FSt %v", err)
	}
	h := s.GetPHashByID(0)
	fmt.Printf("hash: %x\n", h)
}
