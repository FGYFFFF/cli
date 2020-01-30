package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandYamlFileTest(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte("test: 15"), 0666)
	defer os.Remove("current.yaml")
	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("test")
			expect(t, val, 15)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "test", AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))
	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileTestGlobalEnvVarWins(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte("test: 15"), 0666)
	defer os.Remove("current.yaml")

	_ = os.Setenv("THE_TEST", "10")
	defer os.Setenv("THE_TEST", "")
	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("test")
			expect(t, val, 10)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "test", EnvVars: []string{"THE_TEST"}, AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))

	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileTestGlobalEnvVarWinsNested(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte(`top:
  test: 15`), 0666)
	defer os.Remove("current.yaml")

	_ = os.Setenv("THE_TEST", "10")
	defer os.Setenv("THE_TEST", "")
	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("top.test")
			expect(t, val, 10)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "top.test", EnvVars: []string{"THE_TEST"}, AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))

	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileTestSpecifiedFlagWins(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte("test: 15"), 0666)
	defer os.Remove("current.yaml")

	test := []string{"test-cmd", "--load", "current.yaml", "--test", "7"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("test")
			expect(t, val, 7)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "test", AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))

	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileTestSpecifiedFlagWinsNested(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte(`top:
  test: 15`), 0666)
	defer os.Remove("current.yaml")

	test := []string{"test-cmd", "--load", "current.yaml", "--top.test", "7"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("top.test")
			expect(t, val, 7)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "top.test", AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))

	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileTestDefaultValueFileWins(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte("test: 15"), 0666)
	defer os.Remove("current.yaml")

	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("test")
			expect(t, val, 15)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "test", Value: 7, AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))

	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileTestDefaultValueFileWinsNested(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte(`top:
  test: 15`), 0666)
	defer os.Remove("current.yaml")

	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("top.test")
			expect(t, val, 15)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "top.test", Value: 7, AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))

	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileFlagHasDefaultGlobalEnvYamlSetGlobalEnvWins(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte("test: 15"), 0666)
	defer os.Remove("current.yaml")

	_ = os.Setenv("THE_TEST", "11")
	defer os.Setenv("THE_TEST", "")

	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("test")
			expect(t, val, 11)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "test", Value: 7, EnvVars: []string{"THE_TEST"}, AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))
	err := command.Run(c)

	expect(t, err, nil)
}

func TestCommandYamlFileFlagHasDefaultGlobalEnvYamlSetGlobalEnvWinsNested(t *testing.T) {
	app := &App{}
	set := flag.NewFlagSet("test", 0)
	_ = ioutil.WriteFile("current.yaml", []byte(`top:
  test: 15`), 0666)
	defer os.Remove("current.yaml")

	_ = os.Setenv("THE_TEST", "11")
	defer os.Setenv("THE_TEST", "")

	test := []string{"test-cmd", "--load", "current.yaml"}
	_ = set.Parse(test)

	c := NewContext(app, set, nil)

	command := &Command{
		Name:        "test-cmd",
		Aliases:     []string{"tc"},
		Usage:       "this is for testing",
		Description: "testing",
		Action: func(c *Context) error {
			val := c.Int("top.test")
			expect(t, val, 11)
			return nil
		},
		Flags: []Flag{
			&IntFlag{Name: "top.test", Value: 7, EnvVars: []string{"THE_TEST"}, AlternateSource: true},
			&StringFlag{Name: "load"}},
	}
	command.Before = InitInputSourceWithContext(command.Flags, NewYamlSourceFromFlagFunc("load"))
	err := command.Run(c)

	expect(t, err, nil)
}