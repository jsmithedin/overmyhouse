# OverMyHouse
[![Build Status](https://travis-ci.org/jsmithedin/overmyhouse.svg?branch=master)](https://travis-ci.org/jsmithedin/overmyhouse)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsmithedin/overmyhouse)](https://goreportcard.com/report/github.com/jsmithedin/overmyhouse)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1c04c8ad7be744dfaf136234db208cd7)](https://www.codacy.com/manual/jsmithedin/overmyhouse?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jsmithedin/overmyhouse&amp;utm_campaign=Badge_Grade)

Plane over my house -> ADS-B -> SDR -> BEAST -> This thing -> Twitter

Heavily borrowed from <https://github.com/mtigas/simurgh>

Tweets to <https://twitter.com/overjamieshouse>

## Setup
 1. go build
 2. Stick a .env in the same dir as the binary containing twitter stuff: 
```shell script
consumerkey=
consumersecret=
accesstoken=
accesssecret=
```
 3. ./overmyhouse
## Testing
``` shell script
go test -v
```
