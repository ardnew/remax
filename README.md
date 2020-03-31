# remax
#### Maximize serial terminal based on current window size

The `remax` command will reconfigure a serial terminal to use the full height and width of the containing window. It does not resize your window — it just modifies the terminal mode settings so that applications (`vim`, `less`, etc.) can utilize all available space.

Distributed as a statically-linked ELF executable with no dependencies (X11, Python, or even Go). It is also **very fast**.


## Installation
Either use the Go package manager if you have Go installed:
```sh
go get -v github.com/ardnew/remax
```

Or download the latest pre-compiled executable:

||Linux|
|-:|:-----:|
|**i386**|[remax-0.1.0-386.tar.gz](https://github.com/ardnew/remax/releases/download/v0.1.0/remax-0.1.0-386.tar.gz)|
|**x86_64**|[remax-0.1.0-amd64.tar.gz](https://github.com/ardnew/remax/releases/download/v0.1.0/remax-0.1.0-amd64.tar.gz)|
|**ARM**|[remax-0.1.0-arm.tar.gz](https://github.com/ardnew/remax/releases/download/v0.1.0/remax-0.1.0-arm.tar.gz)|
|**ARM64**|[remax-0.1.0-arm64.tar.gz](https://github.com/ardnew/remax/releases/download/v0.1.0/remax-0.1.0-arm64.tar.gz)|


## Usage
Run without arguments to maximize the terminal size for the current window and print its new dimensions.

See flag `-h` for other options:

```
Usage of remax:
  -changelog
        display change history
  -p    print terminal size without changing it
  -q    suppress all non-error output
  -t duration
        read timeout in response to ANSI sequence (default 2s)
  -version
        display version information
```


## Purpose

Often when logged into a Linux system via serial UART, user applications fail to recognize the available area of the client terminal window. In particular, if `stty size` reports incorrect dimensions, you can expect some applications to misbehave.

Since it is compiled as a static executable, no dependencies — X11, Python, or even Go (after building, of course) — are required. And being a true executable, it runs much faster than an interpreted script implementation. This makes it convenient for embedded Linux with minimal resources, just copy the executable to your target. See [Cross-Compiling](#cross-compiling) for more info.


## Cross-Compiling

#### GNU Make
A [Makefile](Makefile) is provided that compiles and packages a tarball for all supported platforms using the default `all` target (i.e. just run `make` without arguments).

#### Manual
If you want to build the utility yourself for a barebones embedded target without installing a full Go distribution, you will need to cross-compile for that target.

For example, to build for a Raspberry Pi 3/4 running Raspbian (32-bit ARM) from another host (where Go is installed):
```sh
git clone https://github.com/ardnew/remax.git $GOPATH/src/github.com/ardnew/remax
cd $GOPATH/src/github.com/ardnew/remax
GOOS=linux GOARCH=arm go install
# or: go build && cp remax /usr/local/bin
```


## Credits
Details of approach and inspiration for this utility comes from Akkana Peck ([@akkana](https://github.com/akkana), thanks!)
- http://shallowsky.com/blog/hardware/serial-24-line-terminals.html
