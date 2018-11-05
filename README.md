Beget
=====

Beget is a tool for building distributed web crawler infrastructures.
Beget is written in the Go programming language and you will need the Go compiler version 1.9 or later
to build the source code.

Issue the `go build` command from this directory to compile the code:

```
$ go build cmd/beget.go
```

This will create the `beget` executable file in this directory. Run `./beget -h` for usage information
on the command.


Beget can accept batched crawl requests from external sources. It will fetch the
resources and send them over to `repositories`. The requests can be simply fed to Beget via its
standard input or pushed it over a simple HTTP REST API. The default repository is the standard
output. Beget can also be configured to use the local file system as a repository. More sophisticated
repositories can be setup as external services that accept crawled data from Beget over HTTP.
Sophisticated and scalable crawl infrastructures can be setup by composing crawl request generators,
Beget and repositories. Request generators and repositories can be implemented in any language that can
speak HTTP.

The easiest way to get started with Beget is by asking it to fetch resources over its standard input/output.
Note that Beget expects a unique "key" to be associated with each resource. This key can be used to identify
the resource in the repository. Crawl requests submitted over the standard input should be separated by newlines.
The format of a request should be `key,resource_url`. A sample interaction with Beget is given below:

```
$ ./beget
web_crawler,https://en.wikipedia.org/wiki/Web_crawler
```

This request will fetch the Wikipedia article on "Web Crawlers" and output it to the standard ouput. The output
will be in the format - `key,data_length,json_encoded_data`. For example, the above request should output:

```
web_crawler,241433,{"statusCode":200,"status":"200 OK","headers":{...},"body": "..."}
```

It is also possible to have the output redirected to files in a directory. Each resource will be written to a file
named `key`.

```
$ ./beget -repo file -repoconfig .
```

The command line argument `-repo` identifies the type of repository to use and `-repoconfig` is a string of configuration
for the repository. For the `file` repository type, the configuration is just the path to the directory where downloaded
files will be saved.

Now a request like `web_crawler,https://en.wikipedia.org/wiki/Web_crawler` will create a file `./web_crawler` and
write the `json_encoded_data` to that file.

Beget can also run in the HTTP server mode:

```
$ ./beget -httpd true
```

In the HTTP server mode, Beget will accept JSON encoded crawl requests from clients. The format of these requets should be:

```
{
 "repository": string,
 "repositoryConfig": string,
 "resources": object,
 "context": object
}
```

The attribute `repository` identifies the repository type for the request. Usually this is either `file` or `simpleHTTP`.
If it's `file`, the `repositoryConfig` should be a directory path. If it's `simpleHTTP`, `repositoryConfig` should be
the HTTP endpoint of a service to which crawled data needs to be posted.

The attribute `resources` is a map of keys to resources to crawl.
`context` is an object that contain for fine-grained configuration of the fetch requests. Its attributes are:

```
{
 "proxyURL": string (proxy to use for the requests)
 "maxRedirects": int (maximum redirects to follow defaults to 10)
 "timeoutSecs": int (request timeouts, defaults to 5secs)
 "header": object (additional HTTP headers)
}
```

A sample HTTP requestor written in python (controller.py) and an HTTP repository (httprepo.py)
is available in the sample folder. (See sample/README.md for details).