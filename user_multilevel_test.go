package greenery_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

type multipleConfig struct {
	*greenery.BaseConfig
}

var multipleDocs = &greenery.DocSet{
	Short: "Multiple commands",
	Usage: map[string]*greenery.CmdHelp{
		"multi": &greenery.CmdHelp{
			Use:   "[ARGS]",
			Short: "command supports args",
		},
		"multiargs": &greenery.CmdHelp{
			Short: "command supports args",
		},
		"multinoargs": &greenery.CmdHelp{
			Short: "command does not support args",
		},
		"multionlyhelp": &greenery.CmdHelp{
			Short: "command has a nil handler, only help",
		},

		"multiargs>subargs": &greenery.CmdHelp{
			Short: "subcommand supports args",
		},
		"multiargs>subnoargs": &greenery.CmdHelp{
			Short: "subcommand does not support args",
		},
		"multiargs>subonlyhelp": &greenery.CmdHelp{
			Short: "subcommand has a nil handler, only help",
		},

		"multinoargs>subargs": &greenery.CmdHelp{
			Short: "subcommand supports args",
		},
		"multinoargs>subnoargs": &greenery.CmdHelp{
			Short: "subcommand does not support args",
		},
		"multinoargs>subonlyhelp": &greenery.CmdHelp{
			Short: "subcommand has a nil handler, only help",
		},

		"multionlyhelp>subargs": &greenery.CmdHelp{
			Short: "subcommand supports args",
		},
		"multionlyhelp>subnoargs": &greenery.CmdHelp{
			Short: "command only has a help handler",
		},
		"multionlyhelp>subonlyhelp": &greenery.CmdHelp{
			Short: "subcommand has a nil handler, only help",
		},
		"missingparent>something": &greenery.CmdHelp{
			Short: "fail",
		},
	},
	CmdLine: map[string]string{
		"Timeout": "the timeout to use for the fetch",
	},
}

func multipleExec(lcfg greenery.Config, args []string) error {
	cfg := lcfg.(*multipleConfig)
	fmt.Printf("Executing command %s with args %v\n", cfg.GetCurrentCommand(), args)
	return nil
}

func newMultipleConfig() greenery.Config {
	return &multipleConfig{
		BaseConfig: greenery.NewBaseConfig("multiple", map[string]greenery.Handler{
			// now various commands with/without args
			"multiargs<":    multipleExec,
			"multinoargs":   multipleExec,
			"multionlyhelp": nil,

			// and subcommands
			"multiargs>subargs<":     multipleExec,
			"multinoargs>subargs<":   multipleExec,
			"multionlyhelp>subargs<": multipleExec,

			"multiargs>subnoargs":     multipleExec,
			"multinoargs>subnoargs":   multipleExec,
			"multionlyhelp>subnoargs": multipleExec,

			"multiargs>subonlyhelp":     nil,
			"multinoargs>subonlyhelp":   nil,
			"multionlyhelp>subonlyhelp": nil,
		}),
	}
}

func TestBadCmdHandlersCmd(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name:               "nil handler for command with args",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"multionlyhelp>subargs<": nil,
			},
			ExecError: "command multionlyhelp>subargs supports arguments but has a nil handler",
		},
		testhelper.TestCase{
			Name:               "overriding the version command and not setting a handler",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"version": multipleExec,
			},
			ExecError: "Overriding of the version command is not permitted, please use cfg.SetHandler to set your desired handler",
		},
		testhelper.TestCase{
			Name:               "overriding the config command and not setting a handler",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"config": multipleExec,
			},
			ExecError: "Overriding of the config command is not permitted, please use cfg.SetHandler to set your desired handler",
		},
		testhelper.TestCase{
			Name:               "overriding the config init command and not setting a handler",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"config>init": multipleExec,
			},
			ExecError: "Overriding of the config>init command is not permitted, please use cfg.SetHandler to set your desired handler",
		},
		testhelper.TestCase{
			Name:               "overriding the config display command and not setting a handler",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"config>display": multipleExec,
			},
			ExecError: "Overriding of the config>display command is not permitted, please use cfg.SetHandler to set your desired handler",
		},
		testhelper.TestCase{
			Name:               "overriding the config env command and not setting a handler",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"config>env": multipleExec,
			},
			ExecError: "Overriding of the config>env command is not permitted, please use cfg.SetHandler to set your desired handler",
		},
		testhelper.TestCase{
			Name:                "missing usage for multiargs>subargs",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					Usage: map[string]*greenery.CmdHelp{
						"multi":                     &greenery.CmdHelp{},
						"multiargs":                 &greenery.CmdHelp{},
						"multinoargs":               &greenery.CmdHelp{},
						"multionlyhelp":             &greenery.CmdHelp{},
						"multiargs>subnoargs":       &greenery.CmdHelp{},
						"multiargs>subonlyhelp":     &greenery.CmdHelp{},
						"multinoargs>subargs":       &greenery.CmdHelp{},
						"multinoargs>subnoargs":     &greenery.CmdHelp{},
						"multinoargs>subonlyhelp":   &greenery.CmdHelp{},
						"multionlyhelp>subargs":     &greenery.CmdHelp{},
						"multionlyhelp>subnoargs":   &greenery.CmdHelp{},
						"multionlyhelp>subonlyhelp": &greenery.CmdHelp{},
						"missingparent>something":   &greenery.CmdHelp{},
					},
					CmdLine: map[string]string{
						"Timeout": "the timeout to use for the fetch",
					},
				}},
			ExecError: "Cannot find a Usage entry in the docs for command multiargs>subargs",
		},
		testhelper.TestCase{
			Name:                "nil usage for multiargs>subargs",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					Usage: map[string]*greenery.CmdHelp{
						"multi":                     &greenery.CmdHelp{},
						"multiargs":                 &greenery.CmdHelp{},
						"multinoargs":               &greenery.CmdHelp{},
						"multionlyhelp":             &greenery.CmdHelp{},
						"multiargs>subargs":         nil,
						"multiargs>subnoargs":       &greenery.CmdHelp{},
						"multiargs>subonlyhelp":     &greenery.CmdHelp{},
						"multinoargs>subargs":       &greenery.CmdHelp{},
						"multinoargs>subnoargs":     &greenery.CmdHelp{},
						"multinoargs>subonlyhelp":   &greenery.CmdHelp{},
						"multionlyhelp>subargs":     &greenery.CmdHelp{},
						"multionlyhelp>subnoargs":   &greenery.CmdHelp{},
						"multionlyhelp>subonlyhelp": &greenery.CmdHelp{},
						"missingparent>something":   &greenery.CmdHelp{},
					},
					CmdLine: map[string]string{
						"Timeout": "the timeout to use for the fetch",
					},
				}},
			ExecError: "Nil CmdHelp entry for command multiargs>subargs in language en",
		},
		testhelper.TestCase{
			Name:                "nil everything",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					Usage: nil,
					CmdLine: map[string]string{
						"Timeout": "the timeout to use for the fetch",
					},
				}},
			ExecError: "Cannot find a Usage entry in the docs for command",
		},
		testhelper.TestCase{
			Name:               "overriding the root command and not setting a handler",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"": multipleExec,
			},
			ExecError: "Overriding of the root command is not permitted, please use cfg.SetHandler to set your desired handler",
		},
		testhelper.TestCase{
			Name:               "nil handler for command with args 2",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"multiargs<": nil,
			},
			ExecError: "command multiargs supports arguments but has a nil handler",
		},
		testhelper.TestCase{
			Name:               "nil handler for command with args 2",
			OverrideHandlerMap: true,
			HandlerMap: map[string]greenery.Handler{
				"missingparent>something": nil,
			},
			ExecError: "Command missingparent, needed as a parent of missingparent>something, was not defined",
		},

		// -----------------------------------------------

		testhelper.TestCase{
			Name: "multiargs ok args",
			CmdLine: []string{
				"multiargs",
				"1",
				"2",
				"3",
			},
			OutStdOut: "Executing command multiargs with args [1 2 3]\n",
		},
		testhelper.TestCase{
			Name: "multiargs ok noargs",
			CmdLine: []string{
				"multiargs",
			},
			OutStdOut: "Executing command multiargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multiargs>subargs ok args",
			CmdLine: []string{
				"multiargs",
				"subargs",
				"1",
				"2",
				"3",
			},
			OutStdOut: "Executing command multiargs>subargs with args [1 2 3]\n",
		},
		testhelper.TestCase{
			Name: "multiargs>subargs ok noargs",
			CmdLine: []string{
				"multiargs",
				"subargs",
			},
			OutStdOut: "Executing command multiargs>subargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multiargs>subnoargs ok args",
			CmdLine: []string{
				"multiargs",
				"subnoargs",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multiargs>subnoargs ok noargs",
			CmdLine: []string{
				"multiargs",
				"subnoargs",
			},
			OutStdOut: "Executing command multiargs>subnoargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multiargs>subonlyhelp ok args",
			CmdLine: []string{
				"multiargs",
				"subonlyhelp",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multiargs>subonlyhelp ok onlyhelp",
			CmdLine: []string{
				"multiargs",
				"subonlyhelp",
			},
			OutStdOutRegex: "command has a nil handler, only help\n\n",
		},

		// --------------------------------------------------------
		testhelper.TestCase{
			Name: "multinoargs ok args",
			CmdLine: []string{
				"multinoargs",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multinoargs ok noargs",
			CmdLine: []string{
				"multinoargs",
			},
			OutStdOut: "Executing command multinoargs with args []\n",
		},

		testhelper.TestCase{
			Name: "multinoargs>subargs ok args",
			CmdLine: []string{
				"multinoargs",
				"subargs",
				"1",
				"2",
				"3",
			},
			OutStdOut: "Executing command multinoargs>subargs with args [1 2 3]\n",
		},
		testhelper.TestCase{
			Name: "multinoargs>subargs ok noargs",
			CmdLine: []string{
				"multinoargs",
				"subargs",
			},
			OutStdOut: "Executing command multinoargs>subargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multinoargs>subnoargs ok args",
			CmdLine: []string{
				"multinoargs",
				"subnoargs",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multinoargs>subnoargs ok noargs",
			CmdLine: []string{
				"multinoargs",
				"subnoargs",
			},
			OutStdOut: "Executing command multinoargs>subnoargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multinoargs>subonlyhelp ok args",
			CmdLine: []string{
				"multinoargs",
				"subonlyhelp",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multinoargs>subonlyhelp ok onlyhelp",
			CmdLine: []string{
				"multinoargs",
				"subonlyhelp",
			},
			OutStdOutRegex: "command has a nil handler, only help\n\n",
		},

		// ---------------------------------------------------------------
		testhelper.TestCase{
			Name: "multionlyhelp ok 1",
			CmdLine: []string{
				"multionlyhelp",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multionlyhelp ok",
			CmdLine: []string{
				"multionlyhelp",
			},
			OutStdOutRegex: "command has a nil handler, only help\n\n",
		},
		testhelper.TestCase{
			Name: "multionlyhelp>subargs ok args",
			CmdLine: []string{
				"multionlyhelp",
				"subargs",
				"1",
				"2",
				"3",
			},
			OutStdOut: "Executing command multionlyhelp>subargs with args [1 2 3]\n",
		},
		testhelper.TestCase{
			Name: "multionlyhelp>subargs ok noargs",
			CmdLine: []string{
				"multionlyhelp",
				"subargs",
			},
			OutStdOut: "Executing command multionlyhelp>subargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multionlyhelp>subnoargs ok args",
			CmdLine: []string{
				"multionlyhelp",
				"subnoargs",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multionlyhelp>subnoargs ok noargs",
			CmdLine: []string{
				"multionlyhelp",
				"subnoargs",
			},
			OutStdOut: "Executing command multionlyhelp>subnoargs with args []\n",
		},
		testhelper.TestCase{
			Name: "multionlyhelp>subonlyhelp ok args",
			CmdLine: []string{
				"multionlyhelp",
				"subonlyhelp",
				"1",
				"2",
				"3",
			},
			ExecError: "unknown command \"1\" for",
		},
		testhelper.TestCase{
			Name: "multionlyhelp>subonlyhelp ok onlyhelp",
			CmdLine: []string{
				"multionlyhelp",
				"subonlyhelp",
			},
			OutStdOutRegex: "command has a nil handler, only help\n\n",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newMultipleConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": multipleDocs}})

	require.NoError(t, err)
}
