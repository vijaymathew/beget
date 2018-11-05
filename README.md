Beget
=====

Beget is a tool for building distributed web crawler infrastructures.
Beget is written in the Go programming language and you will need the Go compiler version 1.9 or later
to build the source code.

Beget can accept batched crawl requests from external sources. It will fetch the
resources and send them over to `repositories`. The requests can be simply fed to Beget via its
standard input or pushed it over a simple HTTP REST API. The default repository is the standard
output. Beget can also be configured to use the local file system as a repository. More sophisticated
repositories can be setup as external services that accept crawled data from Beget over HTTP.
Sophisticated and scalable crawl infrastructures can be setup by composing crawl request generators,
Beget and repositories. Request generators and repositories can be implemented in any language that can
speak HTTP.