package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

var (
	inx  int
	exec string // The exec file and paramter. Must use absolute path

	flagSet     = flag.NewFlagSet("godaemonize", flag.ExitOnError)
	wdir        = flagSet.String("d", "", "Set daemon's working directory to <dir>")
	stderr      = flagSet.String("e", "", "Send daemon's stderr to file, default is <stderr>")
	stdout      = flagSet.String("o", "", "Send daemon's stdout to file, default is <stdout>")
	pidfile     = flagSet.String("p", "", "Save PID to <pidfile>")
	guser       = flagSet.String("u", "", "Run daemon as user <user>. Requires invocation as root")
	environment = flagSet.String("E", "", "Pass environment setting to daemon. like [a=b,c=d]")
)

func usage() {
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "godaemonize, version %s\n", GODAEMONIZE_VERSION)
		fmt.Fprintln(os.Stderr, "Usage: godaemonize [OPTIONS] -x file [ARGV] ...\n")
		fmt.Fprintln(os.Stderr, "OPTIONS\n")
		flagSet.PrintDefaults()
	}

	//-x 之后的参数透传给守护进程
	inx = len(os.Args)
	var haveX bool
	for i, a := range os.Args {
		if haveX == true {
			exec = a
			break
		}
		if a == "-x" {
			inx = i
			haveX = true
		}
	}

	flagSet.Parse(os.Args[1:inx])

	if len(os.Args) < 2 || exec == "" {
		flagSet.Usage()
		os.Exit(2)
	}
}

func main() {
	//设为单进程执行
	runtime.GOMAXPROCS(1)

	usage()

	files := make([]uintptr, 3)
	files[0] = os.Stdin.Fd()
	files[1] = os.Stdout.Fd()
	files[2] = os.Stderr.Fd()

	if syscall.Getppid() == 1 {
		daemon()
		//完成历史使命
		os.Exit(0)
	}

	//fork 新的进程
	attrs := syscall.ProcAttr{Files: files}
	_, err := syscall.ForkExec(os.Args[0], os.Args, &attrs)

	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create master process %s", err)
		os.Exit(2)
	}

	os.Exit(0)
}
