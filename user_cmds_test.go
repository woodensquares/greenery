package greenery_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

type cmdsBaseConfig struct {
	*greenery.BaseConfig
}

var cmdTestName = "cmds_test"

func newCmdsBaseConfig() greenery.Config {
	cfg := &cmdsBaseConfig{
		BaseConfig: greenery.NewBaseConfig(cmdTestName, nil),
	}

	return cfg
}

func TestVersion(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Default version string",
			CmdLine: []string{
				"version",
			},
			OutStdOut: "0.0\n",
		},
		testhelper.TestCase{
			Name: "x.x.x version string",
			CmdLine: []string{
				"version",
			},
			OutStdOut: "1.2.3a-testing\n",
			ConfigDefaults: &greenery.BaseConfigOptions{
				VersionMajor:      "1",
				VersionMinor:      "2",
				VersionPatchlevel: "3a-testing",
			},
		},
		testhelper.TestCase{
			Name: "x.x version string",
			CmdLine: []string{
				"version",
			},
			OutStdOut: "1.2-testing\n",
			ConfigDefaults: &greenery.BaseConfigOptions{
				VersionMajor: "1",
				VersionMinor: "2-testing",
			},
		},
		testhelper.TestCase{
			Name: "x version string",
			CmdLine: []string{
				"version",
			},
			OutStdOut: "3\n",
			ConfigDefaults: &greenery.BaseConfigOptions{
				VersionMajor: "3",
			},
		},
		testhelper.TestCase{
			Name: "custom version string (precedence over x.x.x)",
			CmdLine: []string{
				"version",
			},
			OutStdOut: "custom build xxxx on yyyy-mm-dd\n",
			ConfigDefaults: &greenery.BaseConfigOptions{
				VersionFull:       "custom build xxxx on yyyy-mm-dd",
				VersionMajor:      "1",
				VersionMinor:      "2",
				VersionPatchlevel: "3",
			},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newCmdsBaseConfig,
	})

	require.NoError(t, err)
}

func TestConfigInitDefault(t *testing.T) {
	goldDefault := filepath.Join("testdata", cmdTestName+".TestConfig.defaultcfg")
	goldDefaultV0 := filepath.Join("testdata", cmdTestName+".TestConfig.defaultv0cfg")

	cwd, err := os.Getwd()
	require.NoError(t, err)
	cwdConf := filepath.Join(cwd, cmdTestName+".toml")
	cfDir := configdir.New(cmdTestName, cmdTestName+".toml")
	userConf := cfDir.QueryFolders(configdir.Global)[0].Path
	systemConf := cfDir.QueryFolders(configdir.System)[0].Path

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Default config in cwd",
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefault, Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "^Configuration file generated at " + cwdConf + "\n$",
		},
		testhelper.TestCase{
			Name: "Default config in cwd, no output",
			CmdLine: []string{
				"config",
				"init",
				"-v",
				"0",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 0, Accessor: "GetTyped"},
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefaultV0, Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
		},
		// Note using afero it allows us to do without worrying about
		// permissions in /etc for example.
		testhelper.TestCase{
			Name: "Default config in system",
			CmdLine: []string{
				"config",
				"init",
				"--location",
				"system",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CfgLocation": testhelper.Comparer{Value: "system", Accessor: "GetTyped"},
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: systemConf, Source: goldDefault, Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "^Configuration file generated at " + systemConf + "\n$",
		},
		testhelper.TestCase{
			Name: "Default config in user",
			CmdLine: []string{
				"config",
				"init",
				"--location",
				"user",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CfgLocation": testhelper.Comparer{Value: "user", Accessor: "GetTyped"},
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: userConf, Source: goldDefault, Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "^Configuration file generated at " + userConf + "\n$",
		},
		testhelper.TestCase{
			Name: "Existing conf file",
			CmdLine: []string{
				"config",
				"init",
			},
			PrecreateFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefault, Perms: 0644},
			},
			ExecError:      "Cannot create config file",
			RealFilesystem: true,
		},
		testhelper.TestCase{
			Name: "Existing conf file force",
			CmdLine: []string{
				"config",
				"init",
				"--force",
			},
			PrecreateFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: goldDefault, Perms: 0644},
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"CfgForce": testhelper.Comparer{Value: true},
			},
			OutStdOutRegex: "^Configuration file generated at " + cwdConf + "\n$",
			RealFilesystem: true,
		},
		testhelper.TestCase{
			Name: "Display a generated config",
			CmdLine: []string{
				"config",
				"display",
			},
			CfgFile: goldDefault,
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displaydefaultcfg"),
				Custom: testhelper.CompareIgnoreTmp},
		},

		testhelper.TestCase{
			Name: "Conf variable, duplicate, only once in cfg file, ok",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					Test1 string `greenery:"||none, .param,"`
					Test2 string `greenery:"||none, .param,"`
				}{
					BaseConfig: greenery.NewBaseConfig("cmds_test", nil),
					Test1:      "same",
					Test2:      "same",
				}
			},
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: filepath.Join("testdata", cmdTestName+".TestConfig.duplicatecfg"), Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex:      "^Configuration file generated at " + cwdConf + "\n$",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					ConfigFile: map[string]string{
						greenery.DocConfigHeader: "Config generated while testing",
						"Test1":                  "same param doc",
						"Test2":                  "same param doc",
					},
				},
			},
		},
		testhelper.TestCase{
			Name: "Conf variable, duplicate, only once in cfg file but different defaults",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					Test1 string `greenery:"||none, .param,"`
					Test2 string `greenery:"||none, .param,"`
				}{
					BaseConfig: greenery.NewBaseConfig("cmds_test", nil),
					Test1:      "same",
					Test2:      "different",
				}
			},
			CmdLine: []string{
				"config",
				"init",
			},
			ExecError:           "More than one configuration variable corresponds to config file variable .param with different defaults: \"same\" and \"different\" for example",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					ConfigFile: map[string]string{
						"Test1": "same param doc",
						"Test2": "same param doc",
					},
				},
			},
		},
		testhelper.TestCase{
			Name: "Conf variable, duplicate, only once in cfg file but different docs",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					Test1 string `greenery:"||none, .param,"`
					Test2 string `greenery:"||none, .param,"`
				}{
					BaseConfig: greenery.NewBaseConfig("cmds_test", nil),
					Test1:      "same",
					Test2:      "same",
				}
			},
			CmdLine: []string{
				"config",
				"init",
			},
			ExecError:           "More than one configuration variable corresponds to config file variable .param with different doc lines: \"same param doc\" and \"different param doc\" for example",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					ConfigFile: map[string]string{
						"Test1": "same param doc",
						"Test2": "different param doc",
					},
				},
			},
		},
		testhelper.TestCase{
			Name: "Conf variable, no docs",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					Test1 string `greenery:"|someParam|, .param,"`
				}{
					BaseConfig: greenery.NewBaseConfig("cmds_test", nil),
					Test1:      "same",
				}
			},
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: filepath.Join("testdata", cmdTestName+".TestConfig.nodoccfg"), Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex:      "^Configuration file generated at " + cwdConf + "\n$",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					CmdLine: map[string]string{
						"Test1": "some doc",
					},
					ConfigFile: map[string]string{
						greenery.DocConfigHeader: "Config generated while testing",
						"Test1":                  "",
					},
				},
			},
		},
		testhelper.TestCase{
			Name: "Conf variable, no docs no cmdline",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					Test1 string `greenery:"||none, .param,"`
				}{
					BaseConfig: greenery.NewBaseConfig("cmds_test", nil),
					Test1:      "same",
				}
			},
			CmdLine: []string{
				"config",
				"init",
			},
			ExecError:           "Config file variable Test1 has no documentation, set it to empty if it's meant to not have it",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					ConfigFile: map[string]string{
						greenery.DocConfigHeader: "Config generated while testing",
					},
				},
			},
		},
		testhelper.TestCase{
			Name: "Conf variable, fallback docs",
			ConfigGen: func() greenery.Config {
				return &struct {
					*greenery.BaseConfig
					Test1 string `greenery:"|someParam|, .param,"`
				}{
					BaseConfig: greenery.NewBaseConfig("cmds_test", nil),
					Test1:      "same",
				}
			},
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml", Source: filepath.Join("testdata", cmdTestName+".TestConfig.fallbackdoccfg"), Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex:      "^Configuration file generated at " + cwdConf + "\n$",
			OverrideUserDocList: true,
			UserDocList: map[string]*greenery.DocSet{
				"": &greenery.DocSet{
					CmdLine: map[string]string{
						"Test1": "some doc",
					},
					ConfigFile: map[string]string{
						greenery.DocConfigHeader: "Config generated while testing",
					},
				},
			},
		},
	}

	// Override this to make the generated config file header constant
	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: newCmdsBaseConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": &greenery.DocSet{
				ConfigFile: map[string]string{
					greenery.DocConfigHeader: "Config generated while testing",
				},
			},
		}})

	require.NoError(t, err)
}

func TestConfigInitExtra(t *testing.T) {
	goldDefault := filepath.Join("testdata", cmdTestName+".TestConfig.extracfg")

	cwd, err := os.Getwd()
	require.NoError(t, err)
	cwdConf := filepath.Join(cwd, "extra.toml")
	ptime := time.Date(2014, time.June, 3, 12, 8, 32, 454, time.UTC)
	ptime2 := time.Date(1978, time.May, 27, 7, 32, 00, 00, time.UTC)

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Default extra config",
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: "extra.toml", Source: goldDefault, Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "^Configuration file generated at " + cwdConf + "\n$",
		},
		testhelper.TestCase{
			Name: "help for custom flags",
			CmdLine: []string{
				"flag",
				"-h",
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.helpflag.stdout")},
		},
		testhelper.TestCase{
			Name: "Display a generated extra config with custom deserialization",
			CmdLine: []string{
				"config",
				"display",
			},
			CustomVars:   testhelper.ExtraConfigCustomVars,
			CustomParser: testhelper.ExtraConfigCustomParse,
			CfgFile:      goldDefault,
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displayextracfg"),
				Custom: testhelper.CompareIgnoreTmp},
		},
		testhelper.TestCase{
			Name: "Display a generated extra config with custom deserialization, validate different values in config from default",
			CmdLine: []string{
				"config",
				"display",
			},
			CustomVars:   testhelper.ExtraConfigCustomVars,
			CustomParser: testhelper.ExtraConfigCustomParse,
			CfgFile:      filepath.Join("testdata", cmdTestName+".TestConfig.extracustomcfg"),
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.displayextracustomcfg"),
				Custom: testhelper.CompareIgnoreTmp},
			ExpectedValues: map[string]testhelper.Comparer{
				"LogLevel":  testhelper.Comparer{Value: "info", Accessor: "GetTyped"},
				"NoEnv":     testhelper.Comparer{Value: true},
				"Pretty":    testhelper.Comparer{Value: true},
				"Verbosity": testhelper.Comparer{Value: 2, Accessor: "GetTyped"},
				"Bool":      testhelper.Comparer{Value: false},
				"NameValue": testhelper.Comparer{Value: [][]string{
					[]string{"j1", "v1"},
					[]string{"j2", "v2"},
					[]string{"j3", "v3"},
				}},
				"NameEnum": testhelper.Comparer{Value: [][]string{
					[]string{"j1", "a"},
					[]string{"j2", "b"},
					[]string{"j3", "c"},
				}},
				"FlagIP":      testhelper.Comparer{Value: "192.168.1.1", Accessor: "GetTyped"},
				"FlagPort":    testhelper.Comparer{Value: greenery.PortValue(443)},
				"FlagCString": testhelper.Comparer{Value: "HIHIHI", Accessor: "GetTyped"},
				"FlagEnum":    testhelper.Comparer{Value: "c", Accessor: "GetTyped"},
				"FlagInt":     testhelper.Comparer{Value: 500, Accessor: "GetTyped"},
				"Float32":     testhelper.Comparer{Value: float32(13.34)},
				"Float64":     testhelper.Comparer{Value: float64(-57.78)},
				"Int":         testhelper.Comparer{Value: -11},
				"Int8":        testhelper.Comparer{Value: int8(-12)},
				"Int16":       testhelper.Comparer{Value: int16(-13)},
				"Int32":       testhelper.Comparer{Value: int32(-14)},
				"Int64":       testhelper.Comparer{Value: int64(-15)},
				"Uint":        testhelper.Comparer{Value: uint(11)},
				"Uint8":       testhelper.Comparer{Value: uint8(12)},
				"Uint16":      testhelper.Comparer{Value: uint16(13)},
				"Uint32":      testhelper.Comparer{Value: uint32(14)},
				"Uint64":      testhelper.Comparer{Value: uint64(15)},
				"String":      testhelper.Comparer{Value: "other"},
				"PTime":       testhelper.Comparer{Value: &ptime},
				"Time":        testhelper.Comparer{Value: time.Date(2015, time.June, 3, 12, 8, 32, 454, time.UTC)},
				"Duration":    testhelper.Comparer{Value: time.Hour*48 + time.Minute*16 + time.Second*32 + time.Millisecond*44},

				"SliceString": testhelper.Comparer{Value: []string{"hi", "there"}},
				"SliceInt":    testhelper.Comparer{Value: []int{5, 6, 7}},
			},
		},
		testhelper.TestCase{
			Name: "Check env variable overrides for all extra variables",
			CmdLine: []string{
				"version",
			},
			Env: map[string]string{
				"EXTRA_BOOL": "false",

				"EXTRA_FLAGIP":      "10.0.0.1",
				"EXTRA_FLAGPORT":    "333",
				"EXTRA_FLAGENUM":    "a",
				"EXTRA_FLAGCSTRING": "HI",
				"EXTRA_FLAGINT":     "600",

				"EXTRA_FLOAT32": "11.11",
				"EXTRA_FLOAT64": "-11.11",

				"EXTRA_INT":      "-10",
				"EXTRA_INT8":     "-11",
				"EXTRA_INT16":    "-12",
				"EXTRA_INT32":    "-13",
				"EXTRA_INT64":    "-14",
				"EXTRA_UINT":     "10",
				"EXTRA_UINT8":    "11",
				"EXTRA_UINT16":   "12",
				"EXTRA_UINT32":   "13",
				"EXTRA_UINT64":   "14",
				"EXTRA_PTIME":    "1978-05-27T07:32:00Z",
				"EXTRA_TIME":     "1979-05-27T07:32:00Z",
				"EXTRA_DURATION": "115892045000000",

				"EXTRA_STRING": "test",
			},
			CustomVars:   testhelper.ExtraConfigCustomVars,
			CustomParser: testhelper.ExtraConfigCustomParse,
			OutStdOut:    "0.0\n",
			ExpectedValues: map[string]testhelper.Comparer{
				"Bool":        testhelper.Comparer{Value: false},
				"FlagIP":      testhelper.Comparer{Value: "10.0.0.1", Accessor: "GetTyped"},
				"FlagPort":    testhelper.Comparer{Value: greenery.PortValue(333)},
				"FlagCString": testhelper.Comparer{Value: "HI", Accessor: "GetTyped"},
				"FlagEnum":    testhelper.Comparer{Value: "a", Accessor: "GetTyped"},
				"FlagInt":     testhelper.Comparer{Value: 600, Accessor: "GetTyped"},
				"Float32":     testhelper.Comparer{Value: float32(11.11)},
				"Float64":     testhelper.Comparer{Value: float64(-11.11)},
				"Int":         testhelper.Comparer{Value: -10},
				"Int8":        testhelper.Comparer{Value: int8(-11)},
				"Int16":       testhelper.Comparer{Value: int16(-12)},
				"Int32":       testhelper.Comparer{Value: int32(-13)},
				"Int64":       testhelper.Comparer{Value: int64(-14)},
				"Uint":        testhelper.Comparer{Value: uint(10)},
				"Uint8":       testhelper.Comparer{Value: uint8(11)},
				"Uint16":      testhelper.Comparer{Value: uint16(12)},
				"Uint32":      testhelper.Comparer{Value: uint32(13)},
				"Uint64":      testhelper.Comparer{Value: uint64(14)},
				"String":      testhelper.Comparer{Value: "test"},
				"PTime":       testhelper.Comparer{Value: &ptime2},
				"Time":        testhelper.Comparer{Value: time.Date(1979, time.May, 27, 7, 32, 00, 00, time.UTC)},
				"Duration":    testhelper.Comparer{Value: time.Hour*32 + time.Minute*11 + time.Second*32 + time.Millisecond*45},

				"SliceString": testhelper.Comparer{Value: []string{"first", "second", "third"}},
				"SliceInt":    testhelper.Comparer{Value: []int{1, 2, 3, 4}},

				// No env for these or base values, so default
				"LogLevel":  testhelper.Comparer{Value: "error", Accessor: "GetTyped"},
				"Pretty":    testhelper.Comparer{Value: false},
				"Verbosity": testhelper.Comparer{Value: 1, Accessor: "GetTyped"},
				"NameValue": testhelper.Comparer{Value: [][]string{
					[]string{"k1", "v1"},
					[]string{"k2", "v2"},
					[]string{"k3", "v3"},
				}},
				"NameEnum": testhelper.Comparer{Value: [][]string{
					[]string{"k1", "a"},
					[]string{"k2", "b"},
					[]string{"k3", "c"},
				}},
			},
		},
	}

	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)
}

func TestConfigEnv(t *testing.T) {
	goldDefaultExtra := filepath.Join("testdata", cmdTestName+".TestConfig.extracfg")

	tcs1 := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Config env no env affecting noextra",
			CmdLine: []string{
				"config",
				"env",
			},
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfigEnv.default.noenv")},
		},
		testhelper.TestCase{
			Name: "Config env with env affecting noextra",
			CmdLine: []string{
				"config",
				"env",
			},
			Env: map[string]string{
				"SIMPLE_VERBOSITY": "0",
				"SIMPLE_LOGLEVEL":  "debug",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 0, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "debug", Accessor: "GetTyped"},
			},
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfigEnv.default.env")},
		},
	}

	err := testhelper.RunTestCases(t, tcs1, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)

	tcs2 := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Extra config env no env affecting extra",
			CmdLine: []string{
				"config",
				"env",
			},
			CustomVars:   testhelper.ExtraConfigCustomVars,
			CustomParser: testhelper.ExtraConfigCustomParse,
			CfgFile:      goldDefaultExtra,
			GoldStdOut:   &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfigEnv.extra.noenv")},
		},
		testhelper.TestCase{
			Name: "Config env with env affecting extra",
			CmdLine: []string{
				"config",
				"env",
			},
			CustomVars:   testhelper.ExtraConfigCustomVars,
			CustomParser: testhelper.ExtraConfigCustomParse,
			CfgFile:      goldDefaultExtra,
			Env: map[string]string{
				"EXTRA_VERBOSITY": "0",
				"EXTRA_LOGLEVEL":  "debug",
				"EXTRA_STRING":    "hi",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 0, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "debug", Accessor: "GetTyped"},
				"String":    testhelper.Comparer{Value: "hi"},
			},
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfigEnv.extra.env")},
		},
	}

	err = testhelper.RunTestCases(t, tcs2, testhelper.TestRunnerOptions{
		ConfigGen:  testhelper.NewExtraConfig,
		CompareMap: testhelper.ExtraConfigCustomComparers,
		UserDocList: map[string]*greenery.DocSet{
			"": &testhelper.ExtraDocs}})
	require.NoError(t, err)

	// Doc only for conf, not cmd, check that env still works
	confEnvConf := func() greenery.Config {
		return &struct {
			*greenery.BaseConfig
			TestParam string `greenery:"test|testParam|, , TESTPARAM"`
		}{
			BaseConfig: greenery.NewBaseConfig("partial", map[string]greenery.Handler{
				"test": testhelper.NopNoArgs,
			}),
			TestParam: "default",
		}
	}

	docs := map[string]*greenery.DocSet{
		"": &greenery.DocSet{
			Usage: map[string]*greenery.CmdHelp{
				"test": &greenery.CmdHelp{
					Short: "test",
				},
			},
			CmdLine: map[string]string{
				"TestParam": "test parameter",
			},
		}}

	tcs3 := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Partial doc no env",
			CmdLine: []string{
				"config",
				"env",
			},
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfigEnv.partial.noenv")},
		},
		testhelper.TestCase{
			Name: "Partial doc with env affecting",
			CmdLine: []string{
				"config",
				"env",
			},
			Env: map[string]string{
				"PARTIAL_VERBOSITY": "0",
				"PARTIAL_LOGLEVEL":  "debug",
				"PARTIAL_TESTPARAM": "hi",
			},
			ExpectedValues: map[string]testhelper.Comparer{
				"Verbosity": testhelper.Comparer{Value: 0, Accessor: "GetTyped"},
				"LogLevel":  testhelper.Comparer{Value: "debug", Accessor: "GetTyped"},
				"TestParam": testhelper.Comparer{Value: "hi"},
			},
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfigEnv.partial.env")},
		},
	}

	err = testhelper.RunTestCases(t, tcs3, testhelper.TestRunnerOptions{
		ConfigGen:   confEnvConf,
		UserDocList: docs,
	})
	require.NoError(t, err)
}

// Custom flag with only the minimal number of methods required, Set, String
// and UnmarshalText. Note since Type is not present the help will not be as
// good. This will end up serializing in the config file as a string
type minimalValue uint16

func (p *minimalValue) Set(e string) error {
	var i int

	if e == "" {
		i = 0
	} else {
		var err error
		i, err = strconv.Atoi(e)
		if err != nil {
			return fmt.Errorf("%s, cannot be converted to an integer", e)
		}
	}

	if i < 0 || i > 65535 {
		return fmt.Errorf("Minimal flags must be between 0 and 65535, have %d", i)
	}

	*p = minimalValue(i)
	return nil
}
func (p *minimalValue) String() string                  { return strconv.Itoa(int(*p)) }
func (p *minimalValue) UnmarshalText(text []byte) error { return p.Set(string(text)) }

func TestConfigInitCustom(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "Generate the config",
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: cmdTestName + ".toml",
					Source: filepath.Join("testdata", cmdTestName+".TestConfig.customcfg"), Perms: 0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "Configuration file generated",
		},
		testhelper.TestCase{
			Name: "Display it",
			CmdLine: []string{
				"config",
				"display",
			},
			CfgFile: filepath.Join("testdata", cmdTestName+".TestConfig.customcfg"),
			GoldStdOut: &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestConfig.customdisplaycfg"),
				Custom: testhelper.CompareIgnoreTmp},
		},
	}

	// Override this to make the generated config file header constant
	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: func() greenery.Config {
			cfg := &struct {
				*greenery.BaseConfig
				Minimal minimalValue `greenery:"||none,   .minimal, "`
				Special complex128   `greenery:"||custom, .special, "`
			}{
				BaseConfig: greenery.NewBaseConfig(cmdTestName, nil),
				Minimal:    minimalValue(100),
			}
			return cfg
		},
		UserDocList: map[string]*greenery.DocSet{
			"": &greenery.DocSet{
				ConfigFile: map[string]string{
					greenery.DocConfigHeader: "Config generated while testing",
					"Minimal":                "minimal",
					"Special":                "",
				},
			}},
	})

	require.NoError(t, err)
}
