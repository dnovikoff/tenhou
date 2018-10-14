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