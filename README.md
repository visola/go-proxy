## go-proxy

[![Maintainability](https://api.codeclimate.com/v1/badges/0f398d3937f55ddcfc70/maintainability)](https://codeclimate.com/github/visola/go-proxy/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/0f398d3937f55ddcfc70/test_coverage)](https://codeclimate.com/github/visola/go-proxy/test_coverage)
[![Build Status](https://travis-ci.org/visola/go-proxy.svg?branch=master)](https://travis-ci.org/visola/go-proxy)
[![GitHub release](https://img.shields.io/github/release/visola/go-proxy.svg)](https://github.com/visola/go-proxy/releases)

A proxy to aid developers to run and debug multiple services or frontends locally.

<p style="text-align:center">
  <img width="1000px" src="doc/go-proxy-demo.gif" />
</p>

## Getting Started

1. Download a release for your specific operating system and architecture [here](https://github.com/visola/go-proxy/releases)
2. Unzip it and make the executable available in your path
3. Create the `~/.go-proxy` directory and add a mappings file
4. Run `go-proxy`

If you want to run using HTTPS, set the following two environment variables:

```
GO_PROXY_CERT_FILE=/path/to/server.crt
GO_PROXY_CERT_KEY_FILE=/path/to/server.key
```

If you don't have certificate and key files, the server will start using HTTP.

This is a sample output you should get:

```
$ go-proxy
2018/04/04 10:15:13 Starting proxy server...
Starting proxy at: https://localhost:33443
2018/04/04 10:15:13 Starting admin server...
Opening admin server at: http://localhost:1234
```

## Mappings

Mappings go in the `~/.go-proxy` directory. They are YAML files that get loaded during startup. You can put everything in one file or have multiple files, whatever your preference is. All `.yaml` (or `.yml`) files will be read and loaded. You can enable/disable each mapping in the admin UI.

A mapping maps a request path to some resource. The resource can be a local file or an HTTP server local or remote.

To map to a local file (or directory), you add a `static` attribute to your yaml file, like this:

```yaml
static:
  - from: /statics/some_javascript.js
    to: /some/place/local/my_javascript.js
  - from: /static_assets/
    to: /another/directory/
```

To proxy the request to another HTTP server (local or remote), you add a `proxy` attribute instead, like the following:

```yaml
proxy:
  - from: /static # from maps to anything that has this prefix
    to: http://localhost:1243 # prefix from path, will be appended here
```

You can also map a regular expression, using a regexp in either, `static` or `proxy`:

```yaml
proxy:
  - regexp: /static/js(.*\.chunk\.js)
    to: http://127.0.0.1:3000/static/js$1
```

## Building go-proxy

### Pre-requisites

Before anything, you'll need (versions I have local, haven't tested with others):
- Go >1.9.0
- node >8.7.0

Make sure you have all the dependencies installed:

```bash
# Install all Node dependencies
$ npm install

# Install all Go dependencies
$ dep ensure

# Install packr
$ go get -u github.com/gobuffalo/packr/...
```

### Local development

To develop locally, first make sure `packr` doesn't have any boxes encoded as binaries by running the following:

```
$ packr clean
```

Then, start [parcel-bundler](https://parceljs.org/) server using the following:

```
$ npm run start
```

Then develop in Go like you normally would (`go install` -> `go-proxy`). Files will be automatically picked up from the `./dist` directory  and source maps will be loaded. Also, any changes in the static files will be automatically rebuilt by `parcel`.

### Building the final package

There's a script that run the full build cycle. Just run:

```bash
$ bin/run.sh
```

Inside the `build` directory, you'll see the generated packges for each architecture and operating system combination.
