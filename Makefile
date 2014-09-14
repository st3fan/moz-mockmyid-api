# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/

# This Makefile is more used as a list of aliases for common
# commands. Not as a real build system. Maybe that should change. Not
# sure. Makefiles don't seem to be very common for Go projects.

build:
	go build .

test:
	go test

clean:
	rm -f moz-mockmyid-server

test_deps:
	go get -u github.com/st3fan/moz-go-persona

# Shortcuts for Docker

docker-image: moz-mockmyid-api
	docker build -t st3fan/moz-mockmyid-api .

docker-run:
	docker run -i -t -p "8080:8080" st3fan/moz-mockmyid-api

docker-push:
	docker push st3fan/moz-mockmyid-api
