package storage

import (
	"fmt"
	"os"

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

	NTID_FILENAME = "name_to_ind"

	MAXNLEN = 32
	MAXILEN = 2  // ids from 0..99
	MAXHLEN = 64 // for sha256

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
}

func NewFStorage(log *logger.Logger, dir string) (*FStorage, error) {
	return &FStorage{
		dir:  dir,
		buff: make([]byte, BUFFSIZE),
	}, nil
}

func (s *FStorage) AddUser(name, pass string) error {
	return nil
}

func (s *FStorage) GetPHashByID(id uint8) [32]byte {
	return [32]byte{}
}

func (s *FStorage) GetIDByName(name string) (uint8, error) {
	loc := GLOC + "GetIDByName()"

	file, err := os.Open(s.dir + NTID_FILENAME)
	if err != nil {
		return 0, ErrWhenOpenFile(loc)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return 0, ErrCannotGetFileInfo(loc)
	}
	bname := []byte(name)
	left, right := int64(0), fi.Size()/NTID_LINELEN-1
	for left <= right {
		mid := (left + right) / 2
		if _, err := file.ReadAt(s.buff[:NTID_LINELEN-1], mid*NTID_LINELEN); err != nil {
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
		if end && s.buff[len(bname)+1] == WSB {
			fmt.Println(string(s.buff))
			break
		}
	}

	return 1, nil
}
