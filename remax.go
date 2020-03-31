package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardnew/version"
	"golang.org/x/sys/unix"
)

const commandName = "remax"

func init() {
	version.ChangeLog = []version.Change{{
		Package: commandName,
		Version: "0.1.0",
		Date:    "2020 Mar 30",
		Description: []string{
			"initial revision",
		},
	}}
}

func printChangeLog() {
	version.PrintChangeLog()
}

func printVersion() {
	fmt.Printf("%s version %s\n", commandName, version.String())
}

func main() {

	var (
		argChangeLog bool
		argVersion   bool
		argPrint     bool
		argQuiet     bool
		argTimeout   time.Duration
	)

	flag.BoolVar(&argChangeLog, "changelog", false, "display change history")
	flag.BoolVar(&argVersion, "version", false, "display version information")
	flag.BoolVar(&argPrint, "p", false, "print terminal size without changing it")
	flag.BoolVar(&argQuiet, "q", false, "suppress all non-error output")
	flag.DurationVar(&argTimeout, "t", 2*time.Second, "read timeout in response to ANSI sequence")
	flag.Parse()

	if argChangeLog {
		printChangeLog()
	} else if argVersion {
		printVersion()
	} else {

		log.SetFlags(0)

		fatal := func(term *Terminal, fmt string, arg ...interface{}) {
			term.Restore()
			log.Fatalf(fmt, arg...)
		}

		term, err := RawTerminal(os.Stdin)
		if err != nil {
			log.Fatalf("failed to get raw terminal: %+v", err)
		}

		chs := make(chan string)
		go func(t *Terminal, c chan string) {
			r := bufio.NewReader(os.Stdin)
			if s, err := r.ReadString('R'); err == nil {
				c <- s
			}
		}(term, chs)

		if _, err := fmt.Fprintf(os.Stdin, "\033[9999;9999H\033[6n"); err != nil {
			fatal(term, "failed to print ANSI escape sequence: %+v", err)
		}

		var rows, cols uint
		select {
		case <-time.After(argTimeout):
			term.Restore()
			log.Fatal("failed to read terminal response")
		case s := <-chs:
			term.Restore()
			fmt.Sscanf(s[2:], `%d;%dR`, &rows, &cols)
		}

		if !argPrint {
			if err := term.SetSize(rows, cols); err != nil {
				log.Fatalf("failed to set terminal size: %+v", err)
			}
		}

		if !argQuiet || argPrint {
			log.Println()
			log.Printf("terminal size: rows=%d, cols=%d", rows, cols)
		}
	}
}

type Terminal struct {
	fildes uintptr
	status unix.Termios
	backup unix.Termios
}

func GetTerminal(tty *os.File) (*Terminal, error) {
	if tty == nil {
		return nil, fmt.Errorf("null terminal reference")
	}
	t := &Terminal{fildes: tty.Fd()}
	if err := t.ioctlGet(unix.TCGETS); err != nil {
		return nil, err
	}
	t.backup = t.status
	return t, nil
}

func RawTerminal(tty *os.File) (*Terminal, error) {
	t, err := GetTerminal(tty)
	if err != nil {
		return nil, err
	}
	if err := t.RawMode(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Terminal) ioctlGet(req uint) (err error) {
	var arg *unix.Termios
	if arg, err = unix.IoctlGetTermios(int(t.fildes), req); err != nil {
		return fmt.Errorf("IoctlGetTermios(%d): %+v", t.fildes, err)
	}
	// update the receiver only on success
	t.status = *arg
	return nil
}

func (t *Terminal) ioctlSet(req uint, arg unix.Termios) (err error) {
	if err = unix.IoctlSetTermios(int(t.fildes), req, &arg); err != nil {
		return fmt.Errorf("IoctlSetTermios(%d): %+v", t.fildes, err)
	}
	// update the receiver only on success
	t.status = arg
	return nil
}

func (t *Terminal) ioctlGetWinsize() (rows uint, cols uint, err error) {
	var arg *unix.Winsize
	if arg, err = unix.IoctlGetWinsize(int(t.fildes), unix.TIOCGWINSZ); err != nil {
		return 0, 0, fmt.Errorf("IoctlGetWinsize(%d): %+v", t.fildes, err)
	}
	return uint(arg.Row), uint(arg.Col), nil
}

func (t *Terminal) ioctlSetWinsize(rows uint, cols uint) (err error) {
	arg := unix.Winsize{Row: uint16(rows), Col: uint16(cols)}
	if err = unix.IoctlSetWinsize(int(t.fildes), unix.TIOCSWINSZ, &arg); err != nil {
		return fmt.Errorf("IoctlSetWinsize(%d): %+v", t.fildes, err)
	}
	return nil
}

func (t *Terminal) Restore() error {
	return t.ioctlSet(unix.TCSETS, t.backup)
}

func (t *Terminal) RawMode() error {
	// source: https://github.com/golang/crypto/blob/master/ssh/terminal/util.go#L37:6
	status := t.status
	status.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	status.Oflag &^= unix.OPOST
	status.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	status.Cflag &^= unix.CSIZE | unix.PARENB
	status.Cflag |= unix.CS8
	status.Cc[unix.VMIN] = 1
	status.Cc[unix.VTIME] = 0
	return t.ioctlSet(unix.TCSETS, status)
}

func (t *Terminal) GetSize() (rows uint, cols uint, err error) {
	return t.ioctlGetWinsize()
}

func (t *Terminal) SetSize(rows uint, cols uint) error {
	return t.ioctlSetWinsize(rows, cols)
}

func (t *Terminal) String() string {
	return fmt.Sprintf("fd=0x%08X, status=%#v", t.fildes, t.status)
}
