# OverMyHouse
[![Build Status](https://travis-ci.org/jsmithedin/overmyhouse.svg?branch=master)](https://travis-ci.org/jsmithedin/overmyhouse)

Plane over my house -> ADS-B -> SDR -> BEAST -> This thing -> Twitter

Heavily borrowed from https://github.com/mtigas/simurgh

Tweets to https://twitter.com/overjamieshouse

## Setup
0. go build
1. Stick a .env in the same dir as the binary containing twitter stuff: 
```shell script
consumerkey=
consumersecret=
accesstoken=
accesssecret=
```
3. ./overmyhouse
4. In another shell feed this from running dump1090 with
```shell script
nc 127.0.0.1 30005 | nc 127.0.0.1 8081
```
## Testing
Eh, soon.
