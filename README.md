# OverMyHouse
[![Go Report Card](https://goreportcard.com/badge/github.com/jsmithedin/overmyhouse)](https://goreportcard.com/report/github.com/jsmithedin/overmyhouse)
[![Coverage Status](https://coveralls.io/repos/github/jsmithedin/overmyhouse/badge.svg?branch=main)](https://coveralls.io/github/jsmithedin/overmyhouse?branch=main)

Plane over my house -> ADS-B -> SDR -> BEAST -> This thing -> Twitter

Heavily borrowed from <https://github.com/mtigas/simurgh>

Tweets to <https://twitter.com/overjamieshouse>

## Setup
1.  go build
2.  Stick a .env in the same dir as the binary containing twitter stuff:
```shell script
consumerkey=
consumersecret=
accesstoken=
accesssecret=
slackwebhook=
```
3.  ./overmyhouse -notify=both # twitter, slack, or both
## Testing
``` shell script
go test -v
```
