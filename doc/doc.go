package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/burgesQ/webfmwk/v4/log"
)

type (
	Command interface {
		Init([]string) error
		Run()
		Name() string
		Description() string
	}

	command struct {
		fs          *flag.FlagSet
		fn          func()
		description string
	}
)

func newCommand(fn func(), name string, description ...string) *command {
	c := &command{
		fs: flag.NewFlagSet(name, flag.ContinueOnError),
		fn: fn,
	}

	if len(description) > 0 {
		c.description = description[0]
	}

	return c
}

func (c *command) Name() string {
	return c.fs.Name()
}
func (c *command) Description() string {
	return c.description
}

func (c *command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *command) Run() {
	fmt.Printf("running %s (%s)\n", c.Name(), c.Description())
	c.fn()
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Command{
		newCommand(hello_world, "hello_world", "hello worl :)"),
		newCommand(url_param, "url_param", "url params"),
		newCommand(query_param, "query_param", "extend the query param struct"),
		newCommand(routes, "routes", "register routes"),
		newCommand(tls, "tls", "boot in tls mode"),
		newCommand(swagger, "swagger", "generate a swagger doc"),
		newCommand(post_content, "post_content", "post and validate content (query param, form & url)"),
		newCommand(handlers, "handlers", "use extra handlers"),
		newCommand(request_id, "request_id", "add a uuid to each request"),
		newCommand(custom_context, "custom_context", "extend, register and use a custom context"),
		newCommand(custom_worker, "custom_worker", "register extra worker"),
		newCommand(panic_to_error, "panic_to_error", "use panic to handle some error case"),
	}

	subcommand := os.Args[1]

	log.SetLogLevel(log.LogDEBUG)

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			cmd.Run()
			return nil
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

// go run . (filename)
// Example :
//   go run . panic_to_error
//   running panic_to_error (use panic to handle some error case)
//   ! ERR  : http server :4242 (*net.OpError): listen tcp :4242: bind: address already in use
func main() {

	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
