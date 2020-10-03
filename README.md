# ticker-go

> Real-time stock tickers from the command-line.

`ticker-go` is a simple ~shell script~ go binary using the Yahoo Finance API as a data source. It features colored output and is able to display pre- and post-market prices.

![ticker-go](https://raw.githubusercontent.com/dcluna/ticker-go/master/screenshot.png)

This is a fork of https://github.com/pstadler/ticker.sh implemented as a single executable for performance. I'm a fan of the original script but the excessive spawning of subshells became very slow when loading. Since I use this program every time I spawn a new shell, performance is important.

## Install

```sh
$ go get github.com/dcluna/ticker-go
```

Or use one of the binaries in the Releases section.

## Usage

```sh
# Single symbol:
$ ./ticker-go AAPL

# Multiple symbols:
$ ./ticker-go AAPL MSFT GOOG BTC-USD

# Read from file:
$ echo "AAPL MSFT GOOG BTC-USD" > ~/.ticker.conf
$ ./ticker-go $(cat ~/.ticker.conf)

# Disable colors:
$ NO_COLOR=1 ./ticker-go AAPL

# Update every five seconds:
$ watch -n 5 -t -c ./ticker-go AAPL MSFT GOOG BTC-USD
# Or if `watch` is not available:
$ while true; do clear; ./ticker-go AAPL MSFT GOOG BTC-USD; sleep 5; done
```

This script works well with [GeekTool](https://www.tynsoe.org/v2/geektool/) and similar software:

```sh
PATH=/usr/local/bin:$PATH
~/GitHub/ticker-go/ticker-go AAPL MSFT GOOG BTC-USD
```
