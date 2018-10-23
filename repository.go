package beget

import (
	"os"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// The Repository interface represents a key-value document store.
// The documents fetched by the crawler will be dumped to a Repository.
// A Repository could be implemented on top of the local file system,
// a remote key-value store or a traditional RDBMS.

type Repository interface {
	// Stores the data under key in the target Repository.
	// Returns true if the write was successfull, false if
	// the key already exists or the write failed.
	// If the write failed, err will contain details of the
	// failure condition.
	Put(key string, data string) (status bool, err error)
}

// A basic Repository that writes data to the local file system.
type fileRepository struct {
	targetDir string
}

// Intialize a fileRepository to point to a target directory.
// Returns the fileRepository.
// Note that the existence of targetDir is checked lazily, when the
// write happens.
func FileRepository (targetDir string) (repository fileRepository) {
	repository.targetDir = targetDir
	return
}

func (repository fileRepository) Put(key string, data string) (bool, error) {
	path := filepath.Join(repository.targetDir, key)
	if _, err := os.Stat(path); os.IsExist(err) {
		return false, nil
	}
	err := ioutil.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return false, fmt.Errorf("fileRepository.Put failed: %v", err)
	}
	return true, nil
}
