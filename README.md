# go-proxy

[![Release](https://img.shields.io/github/release/visola/go-proxy.svg?style=flat-square)](https://github.com/visola/go-proxy/releases/latest)
[![Build Status](https://travis-ci.com/visola/go-proxy.svg?branch=master)](https://travis-ci.com/visola/go-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/visola/go-proxy)](https://goreportcard.com/report/github.com/visola/go-proxy)
[![Maintainability](https://api.codeclimate.com/v1/badges/efcd0e7b3ca56fdd79ee/maintainability)](https://codeclimate.com/github/visola/go-proxy/maintainability)

go-proxy is a server that helps developers work faster in the world of microservices and microfrontends. You run it locally to serve all local traffic from multiple sources into one place, acting as a gateway.

<p style="text-align:center">
  <img width="600px" src="doc/go-proxy_overview.png" />
</p>

This is what it can do for you:

- serve static files from disk
- reverse proxy other http/s servers
- it has an admin UI (and API) that can be used to quickly switch upstreams
- multiple listeners in different ports that can act as different servers

# Getting Started

1. Download the latest release for your system [here](https://github.com/visola/go-proxy/releases/latest)
2. Unzip it and make the executable available in your path
3. Create a `~/.go-proxy` directory and add a mappings file
4. Run `go-proxy`

You should see something like the following:

```
$ go-proxy
2019/12/31 13:08:59 Initializing go-proxy...
2019/12/31 13:08:59 Initializing upstreams...
2019/12/31 13:08:59 Opening admin server at: http://localhost:3000
2019/12/31 13:08:59 Reading configuration directory: /Users/visola/.go-proxy
2019/12/31 13:08:59 Starting proxy at: http://localhost:33080
2019/12/31 13:08:59 Found 4 upstreams in file: /Users/visola/.go-proxy/search.yml
```

## HTTPS

If you want to run using HTTPS, set the following two environment variables:

```
GO_PROXY_CERT_FILE=/path/to/server.crt
GO_PROXY_CERT_KEY_FILE=/path/to/server.key
```

# Usage

<!-- TODO - Fill this up -->

# Building go-proxy

## Pre-requisites

Before anything, you'll need:
- Go >1.13
- Node.js >8.7.0

## Building locally

Just run the build bash script:

```
$ ./scrtips/build.sh
```

This will run the tests, install all the dependencies and build the frontend using Webpack.

If that worked, you can run the package script to generate the package for all plataforms:

```
$ ./scripts/package.sh
```

The output will be in the `build/packages` directory.

## Local development

To serve the static files, `go-proxy` uses `packr`, which can point to local files or pack everything inside the final Go binary.

You need to make sure that you have [packr](https://github.com/gobuffalo/packr) installed and available on your path. Normally the following command should do it:

```
$ go get -u github.com/gobuffalo/packr/v2/packr2
```

Run the following command to ensure that `packr` binary files do not exist:

```
$ packr clean
```

No you can start Webpack to watch the admin frontend with the following:

```
$ cd web
$ npm run start
```

The above will block, so you might want to run it from a separated console.

Now you can change the files and rerun like you would any other go application:

```
$ go run cmd/go-proxy/main.go 
```
