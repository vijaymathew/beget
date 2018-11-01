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

package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"bufio"
	"strings"
	"time"
	"github.com/vijaymathew/beget"
)

func main() {
	httpMode := flag.Bool("httpd", false, "if true, start the HTTP server and receive crawl requests from clients")
	inter := flag.String("nwinterface", "", "network interface for the HTTP server")
	port := flag.Int("port", 8080, "HTTP server port")
	tlsCertFile := flag.String("certfile", "", "TLS certificate file for the HTTP server")
	tlsKeyFile := flag.String("keyfile", "", "TLS key file for the HTTP server")
	logFile := flag.String("logfile", "stdout", "log file name")
	concurrentCrawls := flag.Int("concurrentcrawls", 128, "maximum number of concurrent crawls")
	repo := flag.String("repo", "file", "the repository name")
	repoConfig := flag.String("repoconfig", "", "repository configuration")
	proxy := flag.String("proxyurl", "", "proxy url for the crawler")
	maxRedirs := flag.Int("maxredirects", 10, "maximum number of redirects to follow")
	timeoutSecs := flag.Duration("timeoutsecs", 5 * time.Second, "timeout for crawl requests in seconds")
	flag.Parse()

	var logger *log.Logger
	if *logFile == "stdout" {
		logger = log.New(os.Stdout, "beget: ", log.Lshortfile)
	} else {
		f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return
		}
		logger = log.New(f, "beget: ", log.Lshortfile)
	}

	ctx := beget.NewCrawlContext(*concurrentCrawls, logger)
	if *httpMode {
		startHTTPServer(*inter, *port, *tlsCertFile, *tlsKeyFile, ctx)
	} else {
		httpReqCtx := beget.HTTPRequestContext{ProxyURL: *proxy,
			MaxRedirects: *maxRedirs,
			TimeoutSecs: *timeoutSecs}
		crctx := beget.NewCrawlRequestContext(*repo, *repoConfig, &httpReqCtx)
		stdinRequests(ctx, crctx)
	}
}

func startHTTPServer(inter string, port int, tlsCertFile string, tlsKeyFile string, ctx *beget.CrawlContext) {
	cfg := beget.HTTPServerConfig{Interface: inter, Port: 8080, TLSCertFile: tlsCertFile, TLSKeyFile: tlsKeyFile}
	err := beget.StartHTTPServer(cfg, ctx)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
}

func stdinRequests(ctx *beget.CrawlContext, crctx *beget.CrawlRequestContext) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.Index(line, ",")
		if i > 0 {
			key := line[0:i]
			res := line[i+1:]
			ctx.GetOneResource(key, res, crctx)
		}
	}
}
