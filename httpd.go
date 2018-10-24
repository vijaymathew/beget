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
	"fmt"
	"os"
	"log"
	"context"
	"os/signal"
	"encoding/json"
	"net/http"
)

type CrawlRequest struct {
	RepoType string             `json:"repository"`
	RepoConfig string           `json:"repositoryConfig"`
	Context HTTPRequestContext  `json:"context"`
	Resources map[string]string `json:"resources"`
}

type HTTPServerConfig struct {
	Interface string
	Port int
	TLSCertFile string
	TLSKeyFile string
}

var repoRegistry = map[string]NewRepo{"file": NewFileRepository}

func toCrawlRequestContext(req *CrawlRequest) (*CrawlRequestContext) {
	repo := repoRegistry[req.RepoType](req.RepoConfig)
	ctx := CrawlRequestContext{Repo: repo, Context: &req.Context}
	return &ctx
}

func startJobPopper(jobs chan CrawlRequest, abort chan struct{}, crawlCtx *CrawlContext) {
	for {
		select {
		case req := <-jobs:
			crawlCtx.GetResource(req.Resources, toCrawlRequestContext(&req))
		case <-abort:
			return
		}
	}
}

func StartHTTPServer(config HTTPServerConfig, crawlCtx *CrawlContext) {
	var srv http.Server	
	jobs := make(chan CrawlRequest)
	abort := make(chan struct{})
	go startJobPopper(jobs, abort, crawlCtx)
	
	mux := http.NewServeMux()

	crawlRequestHandler := func (w http.ResponseWriter, r *http.Request) {
		var req CrawlRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%v", err)
			return
		}
		_, ok := repoRegistry[req.RepoType]
		if !ok {
			w.WriteHeader(400)
			fmt.Fprintf(w, "bad repository type: %s", req.RepoType)
			return
		}
		jobs <- req
		w.WriteHeader(202)
	}
	
	mux.HandleFunc("/crawl", crawlRequestHandler)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		<-abort
	}()

	addr := fmt.Sprintf("%s:%d", config.Interface, config.Port)
	certf, keyf := config.TLSCertFile, config.TLSKeyFile
	var err error
	if certf != "" && keyf != "" {
		err = http.ListenAndServeTLS(addr, certf, keyf, mux)
	} else {
		err = http.ListenAndServe(addr, mux)
	}
	
	if err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("failed to start HTTP server: %v", err)
	}

	<-idleConnsClosed
}