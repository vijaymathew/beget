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
	"log"
	"encoding/json"
)

type CrawlContext struct {
	tokens chan struct{}
	Logger *log.Logger
}

type CrawlRequestContext struct {
	Repo Repository
	Context *HTTPRequestContext
}

var repoRegistry = map[string]NewRepo{"file": NewFileRepository,
	"simpleHTTP": NewSimpleHTTPRepository}

func NewCrawlRequestContext(repoName string, repoConfig string, httprctx *HTTPRequestContext) (*CrawlRequestContext) {
	repo := repoRegistry[repoName](repoConfig)
	ctx := CrawlRequestContext{Repo: repo, Context: httprctx}
	return &ctx
}

func NewCrawlContext(maxTokens int, logger *log.Logger) (*CrawlContext) {
	crawlCtx := CrawlContext{tokens: make(chan struct{}, maxTokens), Logger: logger}
	return &crawlCtx
}

func IsValidRepo(repoName string) bool {
	_, ok := repoRegistry[repoName]
	return ok
}

func (crawlCtx *CrawlContext) AcquireToken() {
	crawlCtx.tokens <- struct{}{}
}

func (crawlCtx *CrawlContext) ReleaseToken() {
	<-crawlCtx.tokens
}

func (crawlCtx *CrawlContext) GetResource(resources map[string]string, ctx *CrawlRequestContext) {
	for key, res := range resources {
		go getResource(key, res, ctx, crawlCtx)
	}
}

func (crawlCtx *CrawlContext) GetOneResource(resourceKey string, resource string, ctx *CrawlRequestContext) {
	go getResource(resourceKey, resource, ctx, crawlCtx)
}

func getResource(resourceKey string, resource string, ctx *CrawlRequestContext, crawlCtx *CrawlContext) {
	crawlCtx.AcquireToken()
	defer crawlCtx.ReleaseToken()
	response, err := Get(resource, ctx.Context)
	if err != nil {
		crawlCtx.Logger.Printf("getResource: %v", err)
		return
	}
	saveToRepository(resourceKey, &response, ctx.Repo, crawlCtx)
}

func saveToRepository(resourceKey string, response *HTTPResponse, store Repository, crawlCtx *CrawlContext) {
	bytes, err := json.Marshal(response)
	if err != nil {
		crawlCtx.Logger.Printf("saveToRepository: %v", err)
		return
	}
	ok, err := store.Put(resourceKey, bytes)
	if err != nil {
		crawlCtx.Logger.Printf("saveToRepository: %v", err)
	} else if !ok {
		crawlCtx.Logger.Printf("saveToRepository: key already exists: %s", resourceKey)
	}
}


