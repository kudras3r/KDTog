package storage

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/kudras3r/KDTog/pkg/logger"
)

/*
	At this point, we use simple file-based storage.
	There are two main dict-like files:
		1. name -> id
		2. id -> pHash
	We use a simple line-based format for the files:
		1. name:id
		2. id:pHash
	Max lens for name and id look at conts below.
	Max len for line in file (1) is MAXNLEN + len(':')=1 + MAXILEN + len('\n')=1,
	Where N - name, I - id
*/

const (
	WSB = ' '

	GLOC = "internal/storage/storage.go/"

	NTID_FILENAME = "name_to_id"
	ITH_FILENAME  = "id_to_hash"

	MAXNLEN = 32
	MAXILEN = 2  // ids from 0..99
	MAXHLEN = 32 // for sha256

	NTID_LINELEN = MAXNLEN + 1 + MAXILEN + 1
	ITH_LINELEN  = MAXILEN + 1 + MAXHLEN + 1

	BUFFSIZE = max(NTID_LINELEN, ITH_LINELEN)
)

type Storage interface {
	AddUser(name, pass string) error
	GetIDByName(name string) (uint8, error) // i know, really small :)
	GetPHashByID(id uint8) [32]byte
}

type FStorage struct {
	dir  string
	buff []byte

	log *logger.Logger

	mu sync.Mutex
}

func NewFStorage(log *logger.Logger, dir string) (*FStorage, error) {
	return &FStorage{
		dir:  dir,
		buff: make([]byte, BUFFSIZE),
		log:  log,
	}, nil
}

func (s *FStorage) AddUser(name, pass string) error {
	return nil
}

func (s *FStorage) GetPHashByID(id uint8) [32]byte {
	loc := GLOC + "GetPHashByID()"

	s.log.Infof("GetPHashByID(%d)", id)
	s.mu.Lock()
	defer s.mu.Unlock()

	var hash [32]byte
	file, err := os.Open(s.dir + ITH_FILENAME)
	if err != nil {
		s.log.Errorf("cannot open file %s: %v", loc, err)
		return hash
	}
	defer file.Close()

	if _, err := file.ReadAt(s.buff[:MAXHLEN], int64(id)*ITH_LINELEN+MAXILEN+1); err != nil {
		s.log.Errorf("cannot read file at %s: %v", loc, err)
		return hash
	}

	copy(hash[:], s.buff[:MAXHLEN])
	s.log.Infof("found hash for id %d", id)

	return hash
}

func (s *FStorage) GetIDByName(name string) (uint8, error) {
	loc := GLOC + "GetIDByName()"

	s.log.Infof("GetIDByName(%s)", name)
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.dir + NTID_FILENAME)
	if err != nil {
		s.log.Errorf("cannot open file %s: %v", loc, err)
		return 0, ErrWhenOpenFile(loc)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		s.log.Errorf("cannot get file info %s: %v", loc, err)
		return 0, ErrCannotGetFileInfo(loc)
	}

	bname := []byte(name)
	// bsearch
	left, right := int64(0), fi.Size()/NTID_LINELEN-1
	for left <= right {
		mid := (left + right) / 2
		if _, err := file.ReadAt(s.buff[:NTID_LINELEN-1], mid*NTID_LINELEN); err != nil {
			s.log.Errorf("cannot read file %s: %v", loc, err)
			return 0, ErrCannotReadName(loc)
		}
		end := false
		for i := 0; i < len(bname); i++ {
			if bname[i] != s.buff[i] {
				if s.buff[i] > bname[i] {
					right = mid - 1
				} else if s.buff[i] < bname[i] {
					left = mid + 1
				}
				break
			}
			end = true
		}
		if end {
			sid := string(s.buff[MAXNLEN+1 : MAXNLEN+3])
			uid, _ := strconv.ParseInt(strings.TrimSpace(sid), 10, 16) // we know that it is a valid number

			s.log.Infof("found %s:%s", name, sid)
			return uint8(uid), nil
		}
	}

	s.log.Warnf("name %s not found at: %s", name, loc)
	return 0, ErrNameNotFound(loc)
}
