# tenhou
golang package for working with tenhou net logs and protocol

# About

This package contains functions for parsing/generating:
1. Log files
2. Client messages
3. Server messages


## Example of use

Download example server
```
go get -u github.com/dnovikoff/tenhou/example_server
```

Run application
```
$GOPATH/bin/example_server
```

Add to your `hosts` file
```
127.0.0.1	b.mjv.jp
```

1. Login into flash client http://tenhou.net/0/ .
2. Click on any lobby
3. See start of the game

The example of use ends at this point