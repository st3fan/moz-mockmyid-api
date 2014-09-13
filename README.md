[![Build Status](https://travis-ci.org/st3fan/moz-mockmyid-api.svg?branch=master)](https://travis-ci.org/st3fan/moz-mockmyid-api) [![Coverage Status](https://coveralls.io/repos/st3fan/moz-mockmyid-api/badge.png)](https://coveralls.io/r/st3fan/moz-mockmyid-api)

moz-mockmyid-api
================

This projects provides an API for MockMyID to generate assertion. Normally these assertions are obtained by using Persona in a browser environment: after going throuh the login flow, the assertion will be returned to your application via a JavaScript callback. For testing (server-side) code that requires assertions this is not very practical because that full browser environment is usually not available.

This is where the MockMyID API comes in. You can make a simple call to obtain a valid (but short-lived) assertion for any `@mockmyid.com` address.

```
GET http://localhost:8124/login?email=stefan@mockmyid.com&audience=http://localhost:8080"

{ "email":"stefan@mockmyid.com",
  "audience":"http://localhost:8080",
  "assertion":"eyJhbGciOiJEUzEy...very-long-encoded-assertion...RLn-r9StaxpUw5g==" }
```

Building
--------

This is a Go project with no dependencies. You can simply check it out and run it.

```
git clone https://github.com/st3fan/moz-mockmyid-api.git
cd moz-mockmyid-api
go build
./moz-mockmyid-api
```

The tests require a single dependency, after which you can run the tests:

```
go get github.com/st3fan/moz-go-persona
go test -v
```

This requires an internet connection since the unit test will contact the Persona Verifier to make sure the generated assertion is correct.

Running
-------

Work in progress: a Dockerfile to more easily host this app.
