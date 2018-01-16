package main

import (
	"apiroute/logs"
	"apiroute/managers"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const Name = "apigateway"
const (
	ExitCodeOK int = 0

	ExitCodeError = 10 + iota
	ExitCodeInterrupt
	ExitCodeParseFlagsError
	ExitCodeRunnerError
	ExitCodeConfigError
)

var SignalLookup = map[string]os.Signal{
	"SIGHUP":  syscall.SIGHUP,
	"SIGKILL": syscall.SIGKILL,
	"SIGINT":  syscall.SIGINT,
	"SIGTERM": syscall.SIGTERM,
}

type CLI struct {
	sync.Mutex
	outStream, errStream io.Writer
	signalCh             chan os.Signal
	stopCh               chan struct{}
	stopped              bool
}

func NewCLI(out, err io.Writer) *CLI {
	return &CLI{
		outStream: out,
		errStream: err,
		signalCh:  make(chan os.Signal, 1),
		stopCh:    make(chan struct{}),
	}
}

func (cli *CLI) Run(args []string) int {
	// Parse the flags
	_, err := cli.ParseFlags(args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return cli.handleError(err, ExitCodeParseFlagsError)
	}

	// Initial runner
	runner := managers.NewRunner()
	runner.Init()
	go runner.Start()

	// Listen for signals
	signal.Notify(cli.signalCh)

	for {
		select {
		case err := <-runner.ErrCh:
			return cli.handleError(err, ExitCodeRunnerError)
		case s := <-cli.signalCh:
			logs.Log.Debug("cli receiving signal %q", s)

			switch s {
			case SignalLookup["SIGINT"]:
				fallthrough
			case SignalLookup["SIGTERM"]:
				fallthrough
			case SignalLookup["SIGKILL"]:
				fmt.Fprintf(cli.errStream, "Cleaning up...\n")
				runner.Stop()
				return ExitCodeInterrupt
			case SignalLookup["SIGHUP"]:
				fmt.Println("Reloading...")
				runner.Stop()
				// Configuration reload
				runner.Start()

			default:
				logs.Log.Debug("ignoring signal %q", s)
			}
		case <-cli.stopCh:
			return ExitCodeOK

		}

	}

}

func (cli *CLI) stop() {
	cli.Lock()
	defer cli.Unlock()

	if cli.stopped {
		return

	}

	close(cli.stopCh)
	cli.stopped = true

}

func (cli *CLI) ParseFlags(args []string) (*managers.Config, error) {

	cfg := managers.DefaultConfig()

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() { fmt.Fprintf(cli.errStream, usage, Name) }

	flags.Var((funcVar)(func(s string) error {
		cfg.SetServiceDiscoveryClientAddr(s)
		return nil
	}), "service-discovery-client", "")

	flags.Var((funcVar)(func(s string) error {
		cfg.SetListenPort(s)
		return nil
	}), "listen-port", "")

	// If there was a parser error, stop
	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cli *CLI) handleError(err error, status int) int {
	fmt.Fprintf(cli.errStream, "API Gateway returned errors:\n%s\n", err)
	return status

}

const usage = `
Usage: %s [options]

  An API Gateway Golang Implementation

Options:

  -service-discovery-client=<address>
      Sets the address of the Service Discovery Client
  -listen-port=<port>
      Sets the listen port of the http server
`
