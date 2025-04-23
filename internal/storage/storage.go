package storage

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
	N and I is name and id.
*/

const (
	MAXNLEN  = 64
	MAXILEN  = 3
	NTIDLINE = MAXNLEN + 1 + MAXILEN + 1
)

type Storage interface {
	GetIDByName(name string) (uint8, error) // i know, really small :)
	GetPHashByID(id uint8) ([]byte, error)
}

type FStorage struct {
	path string
	buff []byte
}

func NewFStorage(path string, buffSize uint16) (*FStorage, error) {
	return &FStorage{
		path: path,
		buff: make([]byte, buffSize),
	}, nil
}

func (s *FStorage) GetIDByName(name string) (uint8, error) {
	return 1, nil
}
