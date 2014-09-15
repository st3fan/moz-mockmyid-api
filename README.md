[![Build Status](https://travis-ci.org/st3fan/moz-mockmyid-api.svg?branch=master)](https://travis-ci.org/st3fan/moz-mockmyid-api) [![Coverage Status](https://coveralls.io/repos/st3fan/moz-mockmyid-api/badge.png)](https://coveralls.io/r/st3fan/moz-mockmyid-api)

moz-mockmyid-api
================

This project provides a simple HTTP API to generate Persona Assertions for the [MockMyID Identity Provider](https://mockmyid.com/). Normally these assertions are obtained by using Persona in a full browser environment: after going throuh the typical Persona login flow, the assertion will be returned to your (web) application via a JavaScript callback. For testing (server-side) code that requires assertions this is not very practical because that full browser environment is usually not available or difficult to interface with.

This is where the MockMyID API comes in. You can make a simple call to obtain a valid (but short-lived) assertion for any `@mockmyid.com` email address.

```
GET http://localhost:8080/assertion?email=stefan@mockmyid.com&audience=http://localhost:8080"

{ "email":"stefan@mockmyid.com",
  "audience":"http://localhost:8080",
  "assertion":"eyJhbGciOiJEUzEy...very-long-encoded-assertion...RLn-r9StaxpUw5g==" }
```

You can also request the private key for MockMyID in case you need that for your tests:

```
GET http://localhost:8080/key
{ "algorithm":"DS",
  "x":"385cb3509f086e110c5e24bdd395a84b335a09ae",
  "y":"738e...c929",
  "p":"ff60...0483",
  "q":"e21e04f911d1ed7991008ecaab3bf775984309c3",
  "g":"c52a...4a0f" }
```

(The key is not secret, see [provision.html](https://github.com/callahad/mockmyid/blob/master/public_html/browserid/provision.html))

Building
--------

This is a Go project with no dependencies. You can simply check it out and run it.

```
git clone https://github.com/st3fan/moz-mockmyid-api.git
cd moz-mockmyid-api
go build
./moz-mockmyid-api -address 127.0.0.1 -port 8080 -root /api
```

The command parameters are optional. With the above example the application would be available at `http://127.0.0.1:8080/api/`

The tests require a single dependency, after which you can run the tests:

```
go get github.com/st3fan/moz-go-persona
go test -v
```

This does require an internet connection since the unit test will contact the Persona Verifier to make sure the generated assertion is correct.

Running via Docker
------------------

The easiest way to run this app is to start a docker container. The latest version of this app is available on the Docker Hub as [st3fan/moz-mockmyid-api](https://registry.hub.docker.com/u/st3fan/moz-mockmyid-api/).

(If you want to build your own docker image, you can use the `docker-image` target in the supplied `Makefile`.)

You can boot up the API as follows:

```
docker pull st3fan/moz-mockmyid-api
docker run --name mockmyid-api --publish 8080:8080 st3fan/moz-mockmyid-api
```

You can now expose the application via your preferred front-end web server or proxy.

The Docker image can be configured with the following environment variables:

Variable | Default | Description
-------- | ------- | -------------------
ADDRESS  | 127.0.0.1 | The address to listen on
PORT     | 8080      | The port to listen on
ROOT     | /         | The URL path prefix on which to mount the API

These map directly to the application's `-address`, `-port` and `-root` command line arguments. Usually the defaults are fine.

