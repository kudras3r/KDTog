package storage

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/kudras3r/KDTog/pkg/logger"
)

func TestGetIDByName(t *testing.T) {
	tempDir := t.TempDir() + "/"

	nameToIDFile := filepath.Join(tempDir, "name_to_id")
	err := os.WriteFile(nameToIDFile, []byte("aaame1                          :0 \n"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	l := logger.New("DEBUG")
	s, err := NewFStorage(l, tempDir)
	if err != nil {
		t.Fatalf("error when init FSt: %v", err)
	}

	names := []string{"aaame1", "notFound"}
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
	tempDir := t.TempDir() + "/"

	tmp := sha256.Sum256([]byte("123"))
	data := make([]byte, 36)
	data[0] = byte('0')
	data[1] = byte(' ')
	data[2] = byte(':')
	i := 3
	for _, b := range tmp {
		data[i] = b
		i++
	}
	data[i] = byte('\n')

	idToHashFile := filepath.Join(tempDir, "id_to_hash")
	err := os.WriteFile(idToHashFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	l := logger.New("DEBUG")
	s, err := NewFStorage(l, tempDir)
	if err != nil {
		t.Fatalf("error when init FSt: %v", err)
	}

	h := s.GetPHashByID(0)

	if h != tmp {
		t.Fatalf("expected hash %s, got %x", tmp, h)
	}
	t.Logf("hash for id 0 is %x", h)
}

func TestAddUser(t *testing.T) {
	tempDir := t.TempDir() + "/"

	l := logger.New("DEBUG")
	s, err := NewFStorage(l, tempDir)
	if err != nil {
		t.Fatalf("error when init FSt: %v", err)
	}

	err = s.AddUser("test", "test")
	if err != nil {
		t.Fatalf("error when add user: %v", err)
	}

	nti, err := s.GetIDByName("test")
	if err != nil {
		t.Fatalf("error when getIDByName: %v", err)
	}
	if nti != 0 {
		t.Fatalf("expected id 0, got %d", nti)
	}

	expectedHash := sha256.Sum256([]byte("test"))
	actualHash := s.GetPHashByID(0)
	if actualHash != expectedHash {
		t.Fatalf("expected hash %x, got %x", expectedHash, actualHash)
	}

	t.Logf("user 'test' added successfully with id %d and hash %x", nti, actualHash)
}
