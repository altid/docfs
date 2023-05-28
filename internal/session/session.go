package session

import (
	"context"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/altid/libs/config/types"
	"github.com/altid/libs/markup"
	"github.com/altid/libs/service/commander"
	"github.com/altid/libs/service/controller"
)

type ctlItem int
const (
	ctlCommand ctlItem = iota
	ctlSucceed
	ctlErr
)

type Session struct {
	ctx			context.Context
	cancel		context.CancelFunc
	ctrl		controller.Controller
	Defaults	*Defaults
	Verbose		bool
	debug		func(ctlItem, ...any)
}

type Defaults struct {
	Logdir	types.Logdir `altid:"logdir,no_prompt"`
	TLSCert	string		 `altid:"tlscert,no_prompt"`
	TLSKey	string		 `altid:"tlskey,no_prompt"`
}

func (s *Session) Parse() {
	s.debug = func(ctlItem, ...any) {}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	if s.Verbose {
		s.debug = ctlLogging
	}
}

func (s *Session) Connect(Username string) error {
	return nil
}

func (s *Session) Run(c controller.Controller, cmd *commander.Command) error {
	switch cmd.Name {
	case "open":
		newfile := strings.Join(cmd.Args, " ")
		c.CreateBuffer(path.Base(newfile))
		return parseDocument(c, newfile)
	case "close":
		return c.DeleteBuffer(strings.Join(cmd.Args, " "))
	case "fetch":
		// TODO: Pull the remote resource in
		return nil
	default:
		return errors.New("command not supported")
	}
}

func (s *Session) Start(c controller.Controller) error {
	s.ctrl = c
	return nil
}

func (s *Session) Listen(c controller.Controller) {
	s.Start(c)
	<-s.ctx.Done()
}

func (s *Session) Command(cmd *commander.Command) error {
	s.debug(ctlCommand, cmd)
	return s.Run(s.ctrl, cmd)
}

func (s *Session) Quit() {
	s.cancel()
}

func (s *Session) Handle(bufname string, l *markup.Lexer) error {
	// No op really at the moment, though in the future modification of the buffer makes perfect sense
	return nil
}

func ctlLogging(ctl ctlItem, args ...any) {
	l := log.New(os.Stdout, "docfs ", 0)
	switch ctl {
	case ctlCommand:
		m := args[0].(*commander.Command)
		l.Printf("command name=\"%s\" heading=\"%d\" sender=\"%s\" args=\"%s\" from=\"%s\"", m.Name, m.Heading, m.Sender, m.Args, m.From)
	case ctlSucceed:
		l.Printf("%s succeeded\n", args[0])
	case ctlErr:
		l.Printf("error: err=\"%v\"\n", args[0])
	}
}