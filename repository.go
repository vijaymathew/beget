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
	"bytes"
	"io/ioutil"
	"path/filepath"
	"net/http"
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
	Put(key string, data []byte) (status bool, err error)
}

type NewRepo func (string) (Repository)

// A basic Repository that writes data to the local file system.
type fileRepository struct {
	targetDir string
}

// Intialize a fileRepository to point to a target directory.
// Returns the fileRepository.
// Note that the existence of targetDir is checked lazily, when the
// write happens.
func NewFileRepository(targetDir string) (Repository) {
	repository := fileRepository{targetDir: targetDir}
	return &repository
}

func (repository fileRepository) Put(key string, data []byte) (bool, error) {
	path := filepath.Join(repository.targetDir, key)
	if _, err := os.Stat(path); os.IsExist(err) {
		return false, nil
	}
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return false, fmt.Errorf("fileRepository.Put failed: %v", err)
	}
	return true, nil
}

// A HTTP repository - POST the fetched document to a generic endpoint that accepts JSON data.
type simpleHTTPRepository struct {
	url string
}

func NewSimpleHTTPRepository(url string) (Repository) {
	repository := simpleHTTPRepository{url: url}
	return &repository
}

func (repository simpleHTTPRepository) Put(key string, data []byte) (bool, error) {
	req, err := http.NewRequest("POST", repository.url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 204 {
		return true, nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return false, fmt.Errorf("simpleHTTPRepository.Put failed: HTTPStatus: %v, HTTPResponse: %v",
		resp.Status, body)
}
