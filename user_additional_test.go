package greenery_test

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

func TestBasic(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "default values",
			ExpectedValues: map[string]testhelper.Comparer{
				"CfgForce":          testhelper.Comparer{Value: false},
				"CfgLocation":       testhelper.Comparer{Value: "cwd", Accessor: "GetTyped"},
				"NoEnv":             testhelper.Comparer{Value: false},
				"NoCfg":             testhelper.Comparer{Value: false},
				"Pretty":            testhelper.Comparer{Value: false},
				"Verbosity":         testhelper.Comparer{Value: 1, Accessor: "GetTyped"},
				"VersionFull":       testhelper.Comparer{Value: ""},
				"VersionMajor":      testhelper.Comparer{Value: "0"},
				"VersionMinor":      testhelper.Comparer{Value: "0"},
				"VersionPatchlevel": testhelper.Comparer{Value: ""},
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "constraints on default values",
			ExpectedValues: map[string]testhelper.Comparer{
				"CfgLocation": testhelper.Comparer{Value: greenery.NewDefaultEnumValue("CfgLocation", "cwd", "cwd", "user", "system")},
				"LogLevel":    testhelper.Comparer{Value: greenery.NewDefaultEnumValue("LogLevel", "error", "debug", "info", "warn", "error")},
				"Verbosity":   testhelper.Comparer{Value: greenery.NewDefaultIntValue("Verbosity", 1, 0, 3)},
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
	}

	cfg := testhelper.NewSimpleConfig()
	require.Equal(t, "en", cfg.GetDefaultLanguage())

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestHandlers(t *testing.T) {
	okhandlers := map[string]greenery.Handler{
		"-pre-exec-handler": func(cfg greenery.Config, args []string) error {
			fmt.Println("pre-exec called")
			return nil
		},
		"version": func(cfg greenery.Config, args []string) error {
			fmt.Println("version called")
			return nil
		},
		"config": func(cfg greenery.Config, args []string) error {
			fmt.Println("config called")
			return nil
		},
		"root": func(cfg greenery.Config, args []string) error {
			fmt.Println("root called")
			return nil
		},
	}

	failpreexec := map[string]greenery.Handler{
		"-pre-exec-handler": func(cfg greenery.Config, args []string) error {
			return fmt.Errorf("pre-exec fail")
		},
	}

	failhandlers := map[string]greenery.Handler{
		"-pre-exec-handler": func(cfg greenery.Config, args []string) error {
			fmt.Println("pre-exec called")
			return nil
		},
		"version": func(cfg greenery.Config, args []string) error {
			return fmt.Errorf("fail")
		},
		"config": func(cfg greenery.Config, args []string) error {
			return fmt.Errorf("fail")
		},
		"root": func(cfg greenery.Config, args []string) error {
			return fmt.Errorf("fail")
		},
	}

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "version handler",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			OutStdOut:               "pre-exec called\nversion called\n",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         okhandlers,
		},
		testhelper.TestCase{
			Name: "root handler",
			NoValidateConfigValues:  true,
			OutStdOut:               "pre-exec called\nroot called\n",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         okhandlers,
		},
		testhelper.TestCase{
			Name: "config handler",
			CmdLine: []string{
				"config",
			},
			NoValidateConfigValues:  true,
			OutStdOut:               "pre-exec called\nconfig called\n",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         okhandlers,
		},

		testhelper.TestCase{
			Name: "fail pre-exec handler",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			ExecError:               "pre-exec fail",
			ExecErrorOutput:         true,
			OutStdOut:               "",
			OutStdErrRegex:          "Usage:",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         failpreexec,
		},

		testhelper.TestCase{
			Name: "version handler",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			OutStdOut:               "pre-exec called\n",
			OutStdErrRegex:          "Usage:",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         failhandlers,
			ExecError:               "fail",
			ExecErrorOutput:         true,
		},
		testhelper.TestCase{
			Name: "root handler",
			NoValidateConfigValues:  true,
			OutStdOut:               "pre-exec called\n",
			OutStdErrRegex:          "Usage:",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         failhandlers,
			ExecError:               "fail",
			ExecErrorOutput:         true,
		},
		testhelper.TestCase{
			Name: "config handler",
			CmdLine: []string{
				"config",
			},
			NoValidateConfigValues:  true,
			OutStdOut:               "pre-exec called\n",
			OutStdErrRegex:          "Usage:",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers:         failhandlers,
			ExecError:               "fail",
			ExecErrorOutput:         true,
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

func TestSetOptions(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "set options in a handler",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			ExecError:               "Configuration options can be changed only before executing",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": func(cfg greenery.Config, args []string) error {
					return cfg.SetOptions(greenery.BaseConfigOptions{
						VersionMajor:      "1",
						VersionMinor:      "2",
						VersionPatchlevel: "3a-testing",
					})
				},
			},
		},
	}
	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

// --------------------------------------------------
type currentCmdConfig struct {
	*greenery.BaseConfig
}

func caller(lcfg greenery.Config, args []string) error {
	fmt.Printf("Called %s\n", lcfg.GetCurrentCommand())
	return nil
}

func newCurrentCmdConfig() greenery.Config {
	cfg := &currentCmdConfig{
		BaseConfig: greenery.NewBaseConfig("currentcmd", map[string]greenery.Handler{
			"base1":                caller,
			"base1>sub1":           caller,
			"base1>sub2":           caller,
			"base1>sub2>sub1":      caller,
			"base1>sub2>sub2":      caller,
			"base1>sub2>sub2>sub1": caller,
			"base2":                caller,
			"base2>sub1":           caller,
		}),
	}

	return cfg
}

var currentCmdDocs = &greenery.DocSet{
	Usage: map[string]*greenery.CmdHelp{
		"base1": &greenery.CmdHelp{
			Short: "base1",
		},
		"base1>sub1": &greenery.CmdHelp{
			Short: "base1>sub1",
		},
		"base1>sub2": &greenery.CmdHelp{
			Short: "base1>sub2",
		},
		"base1>sub2>sub1": &greenery.CmdHelp{
			Short: "base1>sub2>sub1",
		},
		"base1>sub2>sub2": &greenery.CmdHelp{
			Short: "base1>sub2>sub2",
		},
		"base1>sub2>sub2>sub1": &greenery.CmdHelp{
			Short: "base1>sub2>sub2>sub1",
		},
		"base2": &greenery.CmdHelp{
			Short: "base2",
		},
		"base2>sub1": &greenery.CmdHelp{
			Short: "base2>sub1",
		},
	},
}

func TestGetCurrentCommand(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "caller: base1",
			CmdLine: []string{
				"base1",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base1\n",
		},
		testhelper.TestCase{
			Name: "caller: base1>sub1",
			CmdLine: []string{
				"base1",
				"sub1",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base1>sub1\n",
		},
		testhelper.TestCase{
			Name: "caller: base1>sub2",
			CmdLine: []string{
				"base1",
				"sub2",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base1>sub2\n",
		},
		testhelper.TestCase{
			Name: "caller: base1>sub2>sub1",
			CmdLine: []string{
				"base1",
				"sub2",
				"sub1",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base1>sub2>sub1\n",
		},
		testhelper.TestCase{
			Name: "caller: base1>sub2>sub1",
			CmdLine: []string{
				"base1",
				"sub2",
				"sub2",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base1>sub2>sub2\n",
		},
		testhelper.TestCase{
			Name: "caller: base1>sub2>sub2>sub1",
			CmdLine: []string{
				"base1",
				"sub2",
				"sub2",
				"sub1",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base1>sub2>sub2>sub1\n",
		},
		testhelper.TestCase{
			Name: "caller: base2",
			CmdLine: []string{
				"base2",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base2\n",
		},
		testhelper.TestCase{
			Name: "caller: base2>sub1",
			CmdLine: []string{
				"base2",
				"sub1",
			},

			NoValidateConfigValues: true,
			OutStdOutRegex:         "Called base2>sub1\n",
		},
		testhelper.TestCase{
			Name: "caller: config",
			CmdLine: []string{
				"config",
			},

			NoValidateConfigValues:  true,
			OutStdOutRegex:          "Called config\n",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": caller,
				"config":  caller,
				"root":    caller,
			},
		},
		testhelper.TestCase{
			Name: "caller: version",
			CmdLine: []string{
				"version",
			},

			NoValidateConfigValues:  true,
			OutStdOutRegex:          "Called version\n",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": caller,
				"config":  caller,
				"root":    caller,
			},
		},
		testhelper.TestCase{
			Name: "caller: root",
			NoValidateConfigValues:  true,
			OutStdOutRegex:          "Called \n",
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": caller,
				"config":  caller,
				"root":    caller,
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newCurrentCmdConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": currentCmdDocs},
	})
	require.NoError(t, err)
}

// --------------------------------------------
type sharingConfig struct {
	*greenery.BaseConfig

	CmdEnvConf    string `greenery:"apple|cmdenvconf1|&orange|cmdenvconf2|, .cmdenvconf,   CMDENVCONF"`
	CmdEnv        string `greenery:"apple|cmdenv1|&orange|cmdenv2|,                      , CMDENV"`
	CmdOnly       string `greenery:"apple|cmdconf1|&orange|cmdconf2|,,"`
	EnvA          string `greenery:"apple|enva1|&orange|enva2|,                          , CMDENVCOMMON"`
	EnvB          string `greenery:"apple|envb1|&orange|envb2|,                          , CMDENVCOMMON"`
	Single1Env    string `greenery:"apple|single1env|,                                   , CMDENVSINGLE"`
	Single2Env    string `greenery:"apple|single2env|,                                   , CMDENVSINGLE"`
	Single3Env    string `greenery:"apple|single3env|,                                   , CMDENVSINGLE"`
	Single1Cfg    string `greenery:"apple|single1cfg|,                      .singlecfg   ,"`
	Single2Cfg    string `greenery:"apple|single2cfg|,                      .singlecfg   ,"`
	Single3Cfg    string `greenery:"apple|single3cfg|,                      .singlecfg   ,"`
	Single1EnvCfg string `greenery:"apple|single1envcfg|,                   .singleenvcfg, CMDENVCFGSINGLE"`
	Single2EnvCfg string `greenery:"apple|single2envcfg|,                   .singleenvcfg, CMDENVCFGSINGLE"`
	Single3EnvCfg string `greenery:"apple|single3envcfg|,                   .singleenvcfg, CMDENVCFGSINGLE"`
}

func newSharingConfig() greenery.Config {
	cfg := &sharingConfig{
		BaseConfig: greenery.NewBaseConfig("sharing", map[string]greenery.Handler{
			"apple":  testhelper.NopNoArgs,
			"orange": testhelper.NopNoArgs,
		}),
		CmdEnv:     "default",
		CmdOnly:    "default",
		CmdEnvConf: "default",
		EnvA:       "default",
		EnvB:       "default",
	}

	return cfg
}

var sharedDocs = &greenery.DocSet{
	Usage: map[string]*greenery.CmdHelp{
		"apple": &greenery.CmdHelp{
			Short: "test",
		},
		"orange": &greenery.CmdHelp{
			Short: "test",
		},
	},
	CmdLine: map[string]string{
		"CmdEnvConf":    "test parameter",
		"CmdEnv":        "test parameter",
		"CmdOnly":       "test parameter",
		"EnvA":          "test parameter",
		"EnvB":          "test parameter",
		"Single1Env":    "test parameter",
		"Single2Env":    "test parameter",
		"Single3Env":    "test parameter",
		"Single1Cfg":    "test parameter",
		"Single2Cfg":    "test parameter",
		"Single3Cfg":    "test parameter",
		"Single1EnvCfg": "test parameter",
		"Single2EnvCfg": "test parameter",
		"Single3EnvCfg": "test parameter",
	},
}

func TestSharing(t *testing.T) {
	tcs := []testhelper.TestCase{
		// Cmdline only
		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "Custom cmd and env sharing, defaults apple",
			CmdLine: []string{
				"apple",
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env sharing, defaults orange",
			CmdLine: []string{
				"orange",
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env variable, cmd1 apple",
			CmdLine: []string{
				"apple",
				"--cmdenv1",
				"cmd",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnv": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env variable, cmd2 orange",
			CmdLine: []string{
				"orange",
				"--cmdenv2",
				"cmd",
			},

			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnv": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env variable, env apple",
			CmdLine: []string{
				"apple",
			},
			Env: map[string]string{
				"SHARING_CMDENV": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnv": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env variable, env orange",
			CmdLine: []string{
				"orange",
			},
			Env: map[string]string{
				"SHARING_CMDENV": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnv": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env variable, env and cmd apple",
			CmdLine: []string{
				"apple",
				"--cmdenv1",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENV": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnv": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and env variable, env and cmd orange",
			CmdLine: []string{
				"orange",
				"--cmdenv2",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENV": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnv": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and shared env variable, env",
			CmdLine: []string{
				"apple",
			},
			Env: map[string]string{
				"SHARING_CMDENVCOMMON": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"EnvA": testhelper.Comparer{Value: "env"},
				"EnvB": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and shared env variable, env and cmd1",
			CmdLine: []string{
				"apple",
				"--enva1",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVCOMMON": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"EnvA": testhelper.Comparer{Value: "cmd"},
				"EnvB": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and shared env variable, env and cmd2",
			CmdLine: []string{
				"apple",
				"--envb1",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVCOMMON": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"EnvA": testhelper.Comparer{Value: "env"},
				"EnvB": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and shared env variable, env and both cmd",
			CmdLine: []string{
				"apple",
				"--enva1",
				"cmd1",
				"--envb1",
				"cmd2",
			},
			Env: map[string]string{
				"SHARING_CMDENVCOMMON": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"EnvA": testhelper.Comparer{Value: "cmd1"},
				"EnvB": testhelper.Comparer{Value: "cmd2"},
			},
		},

		// cmd + conf + env, make sure precedence is still ok
		testhelper.TestCase{
			Name: "Custom cmd and conf and env variable, cmd and env and conf apple",
			CmdLine: []string{
				"apple",
				"--cmdenvconf1",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVCONF": "env",
			},
			CfgContents: `cmdenvconf = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnvConf": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and conf and env variable, cmd and env and conf orange",
			CmdLine: []string{
				"orange",
				"--cmdenvconf2",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVCONF": "env",
			},
			CfgContents: `cmdenvconf = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnvConf": testhelper.Comparer{Value: "cmd"},
			},
		},

		testhelper.TestCase{
			Name: "Custom cmd and conf and env variable, env and conf apple",
			CmdLine: []string{
				"apple",
			},
			Env: map[string]string{
				"SHARING_CMDENVCONF": "env",
			},
			CfgContents: `cmdenvconf = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnvConf": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and conf and env variable, env and conf orange",
			CmdLine: []string{
				"orange",
			},
			Env: map[string]string{
				"SHARING_CMDENVCONF": "env",
			},
			CfgContents: `cmdenvconf = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnvConf": testhelper.Comparer{Value: "env"},
			},
		},

		testhelper.TestCase{
			Name: "Custom cmd and conf and env variable, conf apple",
			CmdLine: []string{
				"apple",
			},
			CfgContents: `cmdenvconf = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnvConf": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Custom cmd and conf and env variable, conf orange",
			CmdLine: []string{
				"orange",
			},
			CfgContents: `cmdenvconf = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"CmdEnvConf": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared env, cmd + env",
			CmdLine: []string{
				"apple",
				"--single1env",
				"cmd",
				"--single2env",
				"cmd",
				"--single3env",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVSINGLE": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1Env": testhelper.Comparer{Value: "cmd"},
				"Single2Env": testhelper.Comparer{Value: "cmd"},
				"Single3Env": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared env, cmd + env partial",
			CmdLine: []string{
				"apple",
				"--single1env",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVSINGLE": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1Env": testhelper.Comparer{Value: "cmd"},
				"Single2Env": testhelper.Comparer{Value: "env"},
				"Single3Env": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared env, env only",
			CmdLine: []string{
				"apple",
			},
			Env: map[string]string{
				"SHARING_CMDENVSINGLE": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1Env": testhelper.Comparer{Value: "env"},
				"Single2Env": testhelper.Comparer{Value: "env"},
				"Single3Env": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg, cmd + cfg",
			CmdLine: []string{
				"apple",
				"--single1cfg",
				"cmd",
				"--single2cfg",
				"cmd",
				"--single3cfg",
				"cmd",
			},
			CfgContents: `singlecfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1Cfg": testhelper.Comparer{Value: "cmd"},
				"Single2Cfg": testhelper.Comparer{Value: "cmd"},
				"Single3Cfg": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg, cmd + cfg partial",
			CmdLine: []string{
				"apple",
				"--single1cfg",
				"cmd",
			},
			CfgContents: `singlecfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1Cfg": testhelper.Comparer{Value: "cmd"},
				"Single2Cfg": testhelper.Comparer{Value: "cfg"},
				"Single3Cfg": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg, cfg only",
			CmdLine: []string{
				"apple",
			},
			CfgContents: `singlecfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1Cfg": testhelper.Comparer{Value: "cfg"},
				"Single2Cfg": testhelper.Comparer{Value: "cfg"},
				"Single3Cfg": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg and env, cmd + cfg + env",
			CmdLine: []string{
				"apple",
				"--single1envcfg",
				"cmd",
				"--single2envcfg",
				"cmd",
				"--single3envcfg",
				"cmd",
			},
			CfgContents: `singleenvcfg = "cfg"`,
			Env: map[string]string{
				"SHARING_CMDENVCFGSINGLE": "env",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1EnvCfg": testhelper.Comparer{Value: "cmd"},
				"Single2EnvCfg": testhelper.Comparer{Value: "cmd"},
				"Single3EnvCfg": testhelper.Comparer{Value: "cmd"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg and env, cmd + cfg + env partial",
			CmdLine: []string{
				"apple",
				"--single1envcfg",
				"cmd",
			},
			Env: map[string]string{
				"SHARING_CMDENVCFGSINGLE": "env",
			},
			CfgContents: `singleenvcfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1EnvCfg": testhelper.Comparer{Value: "cmd"},
				"Single2EnvCfg": testhelper.Comparer{Value: "env"},
				"Single3EnvCfg": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg and env, cmd + cfg partial",
			CmdLine: []string{
				"apple",
				"--single1envcfg",
				"cmd",
			},
			CfgContents: `singleenvcfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1EnvCfg": testhelper.Comparer{Value: "cmd"},
				"Single2EnvCfg": testhelper.Comparer{Value: "cfg"},
				"Single3EnvCfg": testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg cfg and env, env + cfg",
			CmdLine: []string{
				"apple",
			},
			Env: map[string]string{
				"SHARING_CMDENVCFGSINGLE": "env",
			},
			CfgContents: `singleenvcfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1EnvCfg": testhelper.Comparer{Value: "env"},
				"Single2EnvCfg": testhelper.Comparer{Value: "env"},
				"Single3EnvCfg": testhelper.Comparer{Value: "env"},
			},
		},
		testhelper.TestCase{
			Name: "Normal cmd, shared cfg cfg and env, cfg only",
			CmdLine: []string{
				"apple",
			},
			CfgContents: `singleenvcfg = "cfg"`,
			ExpectedValues: map[string]testhelper.Comparer{
				"Single1EnvCfg": testhelper.Comparer{Value: "cfg"},
				"Single2EnvCfg": testhelper.Comparer{Value: "cfg"},
				"Single3EnvCfg": testhelper.Comparer{Value: "cfg"},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newSharingConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": sharedDocs},
	})
	require.NoError(t, err)
}

// -------------
// Characters coming from http://www.columbia.edu/~fdc/utf8/index.html

type unicodeConfig struct {
	*greenery.BaseConfig

	Aᛖᚻᚹᛦᛚᚳᚢᛗ string `greenery:"☺☺☺>☺☺☺|ᛖᚻᚹᛦᛚᚳᚢᛗ|, γλώσσα.ᚠᛇᚻ,  "`
	T示示示      string `greenery:"☺☺☺>☹☹☹|示示示|, γλώσσα.示示示,"`
}

func newUnicodeConfig() greenery.Config {
	cfg := &unicodeConfig{
		BaseConfig: greenery.NewBaseConfig("sharing", map[string]greenery.Handler{
			"☺☺☺":     nil,
			"☺☺☺>☺☺☺": caller,
			"☺☺☺>☹☹☹": caller,
		}),
		// To avoid issues with exporting start with uppercase characters just
		// in case
		Aᛖᚻᚹᛦᛚᚳᚢᛗ: "ಬಾ ಇಲ್ಲಿ ಸಂಭವಿಸು ಇಂದೆನ್ನ ಹೃದಯದಲಿ",
		T示示示:      "शक्नोम्यत्तुम्",
	}

	return cfg
}

var unicodeDocs = &greenery.DocSet{
	Usage: map[string]*greenery.CmdHelp{
		"☺☺☺": &greenery.CmdHelp{
			Short: "",
		},
		"☺☺☺>☺☺☺": &greenery.CmdHelp{
			Short: "smiley",
		},
		"☺☺☺>☹☹☹": &greenery.CmdHelp{
			Short: "frowny",
		},
	},
	CmdLine: map[string]string{
		"T示示示":      "kanji",
		"Aᛖᚻᚹᛦᛚᚳᚢᛗ": "runes",
	},
	ConfigFile: map[string]string{
		"Aᛖᚻᚹᛦᛚᚳᚢᛗ": "runesc",
		"T示示示":      "kanjic",
	},
}

func TestUnicode(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Unicode command and parameters, cfg file",
			CmdLine: []string{
				"☺☺☺",
				"☺☺☺",
			},
			CfgContents: `[γλώσσα]
ᚠᛇᚻ = "cfg"
示示示 = "cfg"`,
			OutStdOutRegex: "Called ☺☺☺>☺☺☺\n",
			ExpectedValues: map[string]testhelper.Comparer{
				"Aᛖᚻᚹᛦᛚᚳᚢᛗ": testhelper.Comparer{Value: "cfg"},
				"T示示示":      testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Unicode command and parameters, cfg and cmd",
			CmdLine: []string{
				"☺☺☺",
				"☺☺☺",
				"--ᛖᚻᚹᛦᛚᚳᚢᛗ",
				"மொழிகளிலே",
			},
			CfgContents: `[γλώσσα]
ᚠᛇᚻ = "cfg"
示示示 = "cfg"`,
			OutStdOutRegex: "Called ☺☺☺>☺☺☺\n",
			ExpectedValues: map[string]testhelper.Comparer{
				"Aᛖᚻᚹᛦᛚᚳᚢᛗ": testhelper.Comparer{Value: "மொழிகளிலே"},
				"T示示示":      testhelper.Comparer{Value: "cfg"},
			},
		},
		testhelper.TestCase{
			Name: "Unicode command and parameters, cfg and cmd",
			CmdLine: []string{
				"☺☺☺",
				"☹☹☹",
				"--示示示",
				"மொழிகளிலே",
			},
			CfgContents: `[γλώσσα]
ᚠᛇᚻ = "cfg"
示示示 = "cfg"`,
			OutStdOutRegex: "Called ☺☺☺>☹☹☹\n",
			ExpectedValues: map[string]testhelper.Comparer{
				"Aᛖᚻᚹᛦᛚᚳᚢᛗ": testhelper.Comparer{Value: "cfg"},
				"T示示示":      testhelper.Comparer{Value: "மொழிகளிலே"},
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newUnicodeConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": unicodeDocs},
	})
	require.NoError(t, err)
}

func TestGetFS(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "check for the file being created",
			CmdLine: []string{
				"version",
			},
			NoValidateConfigValues: true,
			PrecreateFiles: []testhelper.TestFile{
				testhelper.TestFile{Contents: []byte("123"), Location: "/tmp/whatever"},
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{"version": func(cfg greenery.Config, args []string) error {
				af := cfg.GetFs()
				c, lerr := afero.ReadFile(af, "/tmp/whatever")
				if lerr != nil {
					return lerr
				}

				if len(c) != 3 || c[0] != '1' || c[1] != '2' || c[2] != '3' {
					return fmt.Errorf("no match")
				}

				fmt.Println("ok")
				return nil
			}},
			OutStdOut: "ok\n",
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}
