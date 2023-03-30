package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/burgesQ/webfmwk/v5/log"
	wtls "github.com/burgesQ/webfmwk/v5/tls"
)

type (
	command interface {
		Init([]string) error
		Run()
		Name() string
		Description() string
	}

	cmd struct {
		fs          *flag.FlagSet
		fn          func()
		description string
	}
)

func newCommand(fn func() *webfmwk.Server, name string, description ...string) *cmd {
	c := &cmd{
		fs: flag.NewFlagSet(name, flag.ContinueOnError),
		fn: func() {
			fmt.Println("loading server ...")

			s := fn()

			fmt.Println("starting server ...")

			defer s.WaitAndStop()

			if name == "tls" {
				// start asynchronously on :4242
				s.StartTLS(":4242", wtls.Config{
					Cert:     "/path/to/cert",
					Key:      "/path/to/key",
					Insecure: false,
				})
			} else {
				s.Start(":4242")
			}
		},
	}

	if len(description) > 0 {
		c.description = description[0]
	}

	return c
}

func (c *cmd) Name() string {
	return c.fs.Name()
}

func (c *cmd) Description() string {
	return c.description
}

func (c *cmd) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *cmd) Run() {
	fmt.Printf("running %s (%s)\n", c.Name(), c.Description())
	c.fn()
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command (matching a file name)")
	}

	cmds := []command{
		newCommand(helloWorld, "hello_world", "hello worl :)"),
		newCommand(urlParam, "url_param", "url params"),
		newCommand(queryParam, "query_param", "extend the query param struct"),
		newCommand(routes, "routes", "register routes"),
		newCommand(tls, "tls", "boot in tls mode"),
		newCommand(swagger, "swagger", "generate a swagger doc"),
		newCommand(postContent, "post_content", "post and validate content (query param, form & url)"),
		newCommand(handlers, "handlers", "use extra handlers"),
		newCommand(customContext, "custom_context", "extend, register and use a custom context"),
		newCommand(customWorker, "custom_worker", "register extra worker"),
		newCommand(panicToError, "panic_to_error", "use panic to handle some error case"),
		newCommand(logMe, "logging", "use panic to handle some error case"),
	}

	subcommand := os.Args[1]

	log.SetLogLevel(log.LogDebug)

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			cmd.Run()

			return nil
		}
	}

	return fmt.Errorf("Unknown subcommand: %q", subcommand)
}

// go run . (filename)
// Example :
//
//	go run . panic_to_error
//	running panic_to_error (use panic to handle some error case)
//	! ERR  : http server :4242 (*net.OpError): listen tcp :4242: bind: address already in use
func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
