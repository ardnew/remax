# remax
#### Maximize serial terminal based on current window size

Often when logged into a Linux system via serial UART, user applications fail to recognize the available area of the client terminal window. In particular, if `stty size` reports incorrect dimensions, you can expect some applications to misbehave.

This application will attempt to determine the maximum size available and update the serial terminal line settings automatically.

This utility is compiled as a static executable, so no dependencies — X11, Python, or even Go (after building, of course) — are required. This makes it convenient for embedded Linux with minimal resources, just copy the executable to your target. See [Cross-Compiling](#cross-compiling) for more info.


## Usage
Use flag `-h` for a list of available options.

## Installation
Either use the Go package manager:
```sh
go get github.com/ardnew/remax
```
Or clone this repo and build/install manually:
```sh
git clone https://github.com/ardnew/remax.git $GOPATH/src/github.com/ardnew/remax
cd $GOPATH/src/github.com/ardnew/remax
go install
# or: go build && cp remax /usr/local/bin
```

## Cross-Compiling
If you want to use the utility on a barebones embedded target without installing a full Go distribution, you will need to cross-compile for that target.

For example, to build for a Raspberry Pi 3 or 4 running Raspbian (which is 32-bit only), just set the appropriate environment variables:
```sh
GOOS=linux GOARCH=arm GOARM=7 go build github.com/ardnew/remax
```

If you are targeting 64-bit ARM, use `GOARCH=arm64` and leave `GOARM` unspecified. 

Running `go tool dist list` will print a list of valid GOARCH/GOOS combinations supported by your Go installation.

Alternatively, you can view a list of available target architectures and operating systems here:
- https://golang.org/doc/install/source#environment

## Credits
Details of approach and inspiration for this utility comes from Akkana Peck ([@akkana](https://github.com/akkana), thanks!)
- http://shallowsky.com/blog/hardware/serial-24-line-terminals.html
