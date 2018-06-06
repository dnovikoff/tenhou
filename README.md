# About

[![Build Status](https://travis-ci.org/dnovikoff/tenhou.svg?branch=master)](https://travis-ci.org/dnovikoff/tenhou)
[![Coverage Status](https://img.shields.io/codecov/c/github/dnovikoff/tenhou.svg)](https://codecov.io/gh/dnovikoff/tenhou)
[![Go Report Card](https://goreportcard.com/badge/github.com/dnovikoff/tenhou)](https://goreportcard.com/report/github.com/dnovikoff/tenhou)

golang package for working with tenhou net logs and protocol

This package contains functions for parsing/generating:
1. Log files
2. Client messages
3. Server messages

## Howto use proxy server
I've made this small proxy for researching tenhou.net flash client protocol content.

Download application
```
go get -u github.com/dnovikoff/tenhou/cmd/tenhou-proxy
```

Run application
```
$GOPATH/bin/tenhou-proxy
```

Add to your `hosts` file
```
127.0.0.1	b.mjv.jp
```

Login into flash client http://tenhou.net/0/ .
Application output would look like
```
2018/01/27 00:26:42 Started server on addr ':10080'. Sequence id is 'baab75ab'
2018/01/27 00:26:55 File for new connection is 'baab75ab_0001.log'
2018/01/27 00:26:55 Error: EOF
2018/01/27 00:26:55 Error: Read context done
2018/01/27 00:26:55 Done with 1
2018/01/27 00:26:55 File for new connection is 'baab75ab_0002.log'
2018/01/27 00:26:56 Error: EOF
2018/01/27 00:26:56 Error: Read context done
2018/01/27 00:26:56 Done with 2
2018/01/27 00:26:57 File for new connection is 'baab75ab_0003.log'
```

Protocol logs will appear in workdir.
Short example of log result.
```
Send: <Z />
Send: <HELO name="NoName" tid="f0" sx="M" />
Get: <HELO uname="%4E%6F%4E%61%6D%65" auth="20180127-c5f19a8e" ratingscale="PF3=1.000000&PF4=1.000000&PF01C=0.582222&PF02C=0.501632&PF03C=0.414869&PF11C=0.823386&PF12C=0.709416&PF13C=0.586714&PF23C=0.378722&PF33C=0.535594&PF1C00=8.000000"/>
Send: <AUTH val="20180127-e3afc0df"/>
Send: <PXR V="1" />
Send: <PXR V="1" />
Get: <LN n="bw1aJ1Pm1y" j="B8C4B11B1D4B8D4D24C3C2C3C1B2C2C" g="o4CM3E1Q12Co4g12BM4M12Q12D1e2P2J2G1G1P1D1G2G"/>
Send: <Z />
Get: <LN n="by1aL1Ph1BC" j="B4B4D4B8C3B5B8D4D12B12C3B5C1B4B" g="o4CM3E1Q12Co4k12BM4M12Q12D1e2M2M2G1G1P1D1J2G"/>
Send: <PXR V="129" />
```

## pimbooo example server
Pimboo (Pin-Man-Bamboo) is a small example server for a game of two players.
The server is created in demonstration purposes only and does not aim for a production.

Download example server
```
go get -u github.com/dnovikoff/tenhou/cmd/pimboo-server
```

Run application
```
$GOPATH/bin/pimboo-server
```

Add to your `hosts` file
```
127.0.0.1	b.mjv.jp
```

1. Login into flash client http://tenhou.net/0/ .
2. Click on any lobby
3. See start of the game

![Example of game](https://raw.githubusercontent.com/dnovikoff/tenhou/master/cmd/pimboo-server/example.gif)

The server suggests Ron on any opponent drop and Tsumo on any take.
If you call a Noten-agari, a Furiten-Ron or agari on a wrong tile, you will pay a penalty.

The game continues until one of the opponents will drop under zero points.