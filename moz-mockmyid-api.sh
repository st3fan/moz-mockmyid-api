#!/bin/sh

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/

if [ -z "$ADDRESS" ]; then
    ADDRESS="0.0.0.0"
fi

if [ -z "$ROOT" ]; then
    ROOT="/"
fi

if [ -z "$PORT" ]; then
    PORT="8080"
fi

exec /usr/local/bin/moz-mockmyid-api -address "$ADDRESS" -port "$PORT" -root "$ROOT"
