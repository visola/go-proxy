## go-proxy

[![Maintainability](https://api.codeclimate.com/v1/badges/0f398d3937f55ddcfc70/maintainability)](https://codeclimate.com/github/visola/go-proxy/maintainability)

## Local development and Building

### Pre-requisites

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

To develop locally, start [parcel-bundler](https://parceljs.org/) server using the following:

```
$ npm run start
```

Then develop in Go like you normally would. Files will be automatically picked up from the `./dist` directory
and source maps will be loaded. Also, any changes in the static files will be automatically rebuilt by `parcel`.

### Building the final package

Generate the final static distribution:

```
$ npm run bundle
```

Then run packr to embed the static code into a Go file, as binary:

```
$ packr clean && packr
```

Last, generate the Go binary:

```
$ go install
```
