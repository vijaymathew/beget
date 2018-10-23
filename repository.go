// Copyright 2018, 2019 Vijay Mathew Pandyalakal<vijay.the.lisper@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
