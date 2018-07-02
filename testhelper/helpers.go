package testhelper

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
)

// CompareFunc is an interface that is fulfilled by compare functions, used
// when verifying config struct values
type CompareFunc func(*testing.T, string, interface{}, interface{}) error

// Comparer is a struct containing a single expected value, how to retrieve
// it, and how to compare it.
type Comparer struct {
	Value        interface{}
	Accessor     string
	Compare      CompareFunc
	DefaultValue bool
}

// LogLine contains an expected log line to compare to. Note that lines pass
// through json which will deserialize things like time.Time as strings etc.
type LogLine struct {
	Level       string
	Msg         string
	MsgRegex    string
	Custom      map[string]interface{}
	CustomRegex map[string]string
}

// TestFile describes a file that is expected to be created by the application
// being tested.
type TestFile struct {
	Location string
	Source   string
	Contents []byte
	Perms    os.FileMode
	Custom   func(*testing.T, []byte, []byte) bool
}

// TestCase represent a single test case.
type TestCase struct {
	// Name contains the name of the test, will be used as t.Run()'s name
	Name string

	// Env is a map of strings containing the environment for this test case
	Env map[string]string

	// CmdLine is the command line passed to our executable
	CmdLine []string

	// CfgContents is a string containing the configuration file contents for
	// this test
	CfgContents string

	// ExecError is the expected exec function error (if any) this will be
	// matched as a substring
	ExecError string

	// ExecErrorRegex is the expected exec function error (if any) this will be
	// matched as a regex
	ExecErrorRegex string

	// ExecErrorOutput will cause the output of the program (stdout/err/logs)
	// to be checked even if the execution had an error
	ExecErrorOutput bool

	// ExpectedValues is a map of interface{} containing the values that are
	// to be expected to be set in the config object after the environment,
	// command line and configuration file are taken into account. Any config
	// object fields not included here will be validated to remain their
	// default value as set by the config generator function. This will be
	// ignored if Validate is set
	ExpectedValues map[string]Comparer

	// ValuesValidator contains a custom validation function, which will be used
	// instead of the internal configuration values validation. It will be
	// passed the configuration, it is expected to fail the test itself if the
	// validation is unsuccessful.
	ValuesValidator func(*testing.T, greenery.Config)

	// Trace is used to output trace log messages, trace is an undocumented
	// log level used primarily during development.
	Trace bool

	// CmdlineCfgName contains a name to be passed as the config file name to
	// the test, it is typically used to validate a non-existent configuration
	// file, because configuration files are passed via CfgContents
	CmdlineCfgName string

	// ConfigGen is the function that will return the configuration struct to
	// be used
	ConfigGen func() greenery.Config

	// OverrideHandlerMap will enable overriding of the test handler map via
	// HandlerMap
	OverrideHandlerMap bool

	// HandlerMap is a map of handlers, that will be used on the passed in
	// configuration.
	HandlerMap map[string]greenery.Handler

	// OverrideBuiltinHandlers will enable overriding the test default
	// handlers for the configuration
	OverrideBuiltinHandlers bool

	// BuiltinHandlers is a map of override builtin handlers, that will be
	// used on the passed in configuration. Note the key of the map is a
	// string, not an OverrideHandler type to allow arbitrary test overrides.
	BuiltinHandlers map[string]greenery.Handler

	// OverrideUserDocList allows overriding of the test user doc list
	OverrideUserDocList bool

	// UserDocList is a map with the application user documents, if the map
	// contains a "" key, that will be meant to indicate the default language
	// docset.
	UserDocList map[string]*greenery.DocSet
	CfgFile     string

	// GoldStdOut contains the gold expected standard output for this test case
	GoldStdOut *TestFile

	// GoldStdErr contains the gold expected standard error for this test case
	GoldStdErr *TestFile

	// GoldStdOut contains the gold expected logfile for this test case
	GoldLog string

	// NoValidateConfigValues if set will cause the test to not look at the
	// configuration values
	NoValidateConfigValues bool

	// ConfigDefaults allows specifying config defaults via SetOptions on a
	// per-testcase basis
	ConfigDefaults *greenery.BaseConfigOptions

	// GoldFiles is a list of expected gold files
	GoldFiles []TestFile

	// PrecreateFiles is a list of files to be created in the test environment
	// before executing the test
	PrecreateFiles []TestFile

	// RealFilesystem controls whether the test executes on the filesystem or
	// on an in-memory cache
	RealFilesystem bool

	// RemoveFiles contains a list of files to be removed after the test
	// execution
	RemoveFiles []string

	// CustomVars is a list of expected custom variables in the config file
	CustomVars []string

	// CustomParser contains a custom parser to be used for custom variables
	CustomParser func(greenery.Config, map[string]interface{}) ([]string, error)

	// OverrideTraceLogger controls overriding of the test default trace
	// logger creation function
	OverrideTraceLogger bool

	// TraceLogger contains the MakeTraceLogger function to use
	TraceLogger greenery.MakeTraceLogger

	// OverridePrettyLogger controls overriding of the test default pretty
	// logger creation function
	OverridePrettyLogger bool

	// PrettyLogger contains the MakePrettyLogger function to use
	PrettyLogger greenery.MakeLogger

	// OverrideStructuredLogger controls overriding of the test default
	// structured logger creation function
	OverrideStructuredLogger bool

	// StructuredLogger contains the MakeStructuredLogger function to use
	StructuredLogger greenery.MakeLogger

	// OutStdOut is the expected stdout output from the executed command
	OutStdOut string

	// OutStdOutRegex will be matched against the output from the executed command
	OutStdOutRegex string

	// OutStdErr is the expected stderr output from the executed command
	OutStdErr string

	// OutStdErrRegex will be matched against the stderr output from the executed command
	OutStdErrRegex string

	// OutLogAllLines means one regex per expected log line will be tested
	OutLogAllLines bool

	// OutLogLines contains the expected log lines. All regexes must match at
	// least one of the log lines, done in sequence so go through the log
	// lines until the first regex matches, then the second and so on
	OutLogLines []LogLine

	// OutLogRegex will be matched against the log output from the executed command
	OutLogRegex string

	// Pretty is set to true if the expected logging output is 'pretty'
	// otherwise it's assumed to be standard JSON
	Pretty bool
	// Parallel makes this test case run parallel, parallel tests do not
	// have their stdout/stderr captured, and only OutLog directives are taken
	// into account.
	Parallel bool

	// Internals
	af     afero.Fs
	pretty bool
}

// Some commandline optional flags, users can check these as needed in their
// test logic

// GoldUpdate is the flag controlling whether the gold files should be updated
var GoldUpdate = flag.Bool("test.update-gold", false, "update golden files")

// ForceTrace is the flag controlling if the test should be run with tracing
// enabled for test debugging purposes
var ForceTrace = flag.Bool("test.force-trace", false, "force trace")

// ForceNoParallel is a flag that if set will disable parallel test execution
var ForceNoParallel = flag.Bool("test.force-no-parallel", false, "force no parallel")

func sanityCheck(t *testing.T, tc *TestCase) {
	if tc.Parallel && len(tc.Env) != 0 {
		require.FailNow(t, "Tests run with Parallel=true do not support setting the environment")
	}

	for _, x := range tc.CmdLine {
		if x == "--pretty" {
			tc.pretty = true
		}

		if x == "--log-file" {
			require.FailNow(t, "Tests are run with --log-file automatically, please use OutLogLines / OutLogRegex / GoldLog to test")
		}

		if x == "-c" {
			require.FailNow(t, "-c is not supported, please use ConfigFile / ConfigContents")
		}
	}

	if tc.pretty && (tc.OutLogLines != nil || tc.GoldLog != "") {
		require.FailNow(t, "--pretty is not supported when structured logfile testing is in effect, please use OutLogRegex instead")
	}

	if (len(tc.CustomVars) != 0 && tc.CustomParser == nil) ||
		(len(tc.CustomVars) == 0 && tc.CustomParser != nil) {
		require.FailNow(t, "Either both of CustomVars/CustomParser must be set, or neither")
	}

	if (tc.OutStdOut != "" && (tc.GoldStdOut != nil || tc.OutStdOutRegex != "")) ||
		(tc.GoldStdOut != nil && tc.OutStdOutRegex != "") {
		require.FailNow(t, "Only one of OutStdOut, OutStdOutRegex and GoldStdOut should be set")
	}

	if (tc.OutStdErr != "" && (tc.GoldStdErr != nil || tc.OutStdErrRegex != "")) ||
		(tc.GoldStdErr != nil && tc.OutStdErrRegex != "") {
		require.FailNow(t, "Only one of OutStdErr, OutStdErrRegex and GoldStdErr should be set")
	}

	if tc.ExecError != "" && tc.ExecErrorRegex != "" {
		require.FailNow(t, "Only one of ExecError and ExecErrorRegex should be set")
	}

	for _, v := range tc.GoldFiles {
		if v.Contents != nil {
			require.FailNow(t, "Gold files should not have Contents set, only Location (gold file %s)", v.Location)
		}
	}
}

// CompareIgnoreTmp is a function that will compare two files ignoring any
// temporary files listed in them. This is typically used when comparing
// configuration files where we want to ignore the log file line, say, is
// /tmp/xxx in the currently executing test and /tmp/yyy in the gold file
func CompareIgnoreTmp(t *testing.T, gold, found []byte) bool {
	goldS := strings.Split(string(gold), "\n")
	foundS := strings.Split(string(found), "\n")
	if len(goldS) != len(foundS) {
		t.Logf("Different file lengths: expected %d, found %d", len(goldS), len(foundS))
		return false
	}

	rx := regexp.MustCompile("(.*?)" + filepath.Join(os.TempDir(), "[a-zA-Z0-9.]+") + "(.*?)")
	for i := range goldS {
		g := rx.ReplaceAllString(goldS[i], "$1[FILE]$2")
		f := rx.ReplaceAllString(foundS[i], "$1[FILE]$2")

		if g != f {
			t.Logf("Expected: %s", g)
			t.Logf("Found: %s", f)
			return false
		}
	}
	return true
}

// TestRunnerOptions is a struct containing options for this test runner, see
// the TestCase documentation for information on the fields available.
type TestRunnerOptions struct {
	ConfigGen                func() greenery.Config
	Parallel                 bool
	HandlerMap               map[string]greenery.Handler
	OverrideHandlerMap       bool
	UserDocList              map[string]*greenery.DocSet
	OverrideTraceLogger      bool
	TraceLogger              greenery.MakeTraceLogger
	OverridePrettyLogger     bool
	PrettyLogger             greenery.MakeLogger
	OverrideStructuredLogger bool
	StructuredLogger         greenery.MakeLogger
	CompareMap               map[string]CompareFunc
	OverrideBuiltinHandlers  bool
	BuiltinHandlers          map[string]greenery.Handler
}

// LoggingTester is a greenery handler that will output a log message at all
// possible levels / verbosities.
func LoggingTester(cfg greenery.Config, args []string) error {
	fields := []greenery.LogField{
		cfg.LogInteger("int", 111),
		cfg.LogString("string", "string"),
		cfg.LogGeneric("generic", true),
		cfg.LogTime("time", time.Date(2018, time.June, 1, 12, 15, 31, 5e8, time.UTC)),
		cfg.LogDuration("duration", time.Minute),
	}

	cfg.Debugq("Debugq 1")
	cfg.Debugqf("Debugqf %d", 1)
	cfg.Debugqs("Debugqs 1", fields...)

	cfg.Debug("Debug 2")
	cfg.Debugf("Debugf %d", 2)
	cfg.Debugs("Debugs 2", fields...)

	cfg.Debugv("Debugv 3")
	cfg.Debugvf("Debugvf %d", 3)
	cfg.Debugvs("Debugvs 3", fields...)

	cfg.Infoq("Infoq 1")
	cfg.Infoqf("Infoqf %d", 1)
	cfg.Infoqs("Infoqs 1", fields...)

	cfg.Info("Info 2")
	cfg.Infof("Infof %d", 2)
	cfg.Infos("Infos 2", fields...)

	cfg.Infov("Infov 3")
	cfg.Infovf("Infovf %d", 3)
	cfg.Infovs("Infovs 3", fields...)

	cfg.Warn("Warning")
	cfg.Warnf("Warning %d", 1)
	cfg.Warns("Warning", fields...)

	cfg.Error("Error")
	cfg.Errorf("Error %d", 1)
	cfg.Errors("Error", fields...)

	cfg.Trace("Trace")
	cfg.Tracef("Trace %d", 1)
	return nil
}

// RunTestCases will run the declared test cases in .Run() subtests
func RunTestCases(t *testing.T, tcs []TestCase, global TestRunnerOptions) error {
	mtx := &sync.Mutex{}
	removeList := &([]string{})

	for _, tc := range tcs {
		sanityCheck(t, &tc)

		// Capture tc and run
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			// If updating even if users wants parallel, run sequential to
			// avoid any possible clobbering issues with the gold files
			// updating. Test performance during gold files updating should
			// not matter.
			if (tc.Parallel || global.Parallel) && !*GoldUpdate && !*ForceNoParallel {
				t.Parallel()
			}

			var err error
			compFuncs := map[string]CompareFunc{
				"CfgLocation": CompareGetterToGetter,
				"Verbosity":   CompareGetterToGetter,
				"LogLevel":    CompareGetterToGetter,
			}

			for kk, vv := range global.CompareMap {
				compFuncs[kk] = vv
			}

			// Test-local overrides
			cgen := global.ConfigGen
			if tc.ConfigGen != nil {
				cgen = tc.ConfigGen
			}
			require.NotNil(t, cgen, "Either the per-test ConfigGenerator or the global generator must not be nil")

			userDocList := global.UserDocList
			if tc.OverrideUserDocList {
				userDocList = tc.UserDocList
			}

			if tc.RealFilesystem {
				// This should be part of t.Run() because otherwise if tests
				// are parallel it would be run before the tests have
				// completed. It's ok to call it multiple times with the same
				// removeFiles variable.
				defer removeFiles(t, removeList, mtx)

				tc.af = afero.NewOsFs()
				if len(tc.RemoveFiles) != 0 {
					t.Logf("Add user-provided files %v to the cleanup list", tc.RemoveFiles)
					addCleanupFiles(t, tc.RemoveFiles, removeList, mtx)
				}
			} else {
				// Note unfortunately this does not yet support O_EXCL
				// https://github.com/spf13/afero/pull/102
				// so some tests might want to run on the real filesystem
				// instead. Obviously no need to clean files here.
				tc.af = afero.NewMemMapFs()
			}

			// And create any needed files
			created := precreateFiles(t, tc.af, tc.PrecreateFiles)
			if tc.RealFilesystem {
				if len(created) != 0 {
					t.Logf("Add pre-created files %v to the cleanup list", created)
					addCleanupFiles(t, created, removeList, mtx)
				}
			}

			// Start from a new configuration
			cleanCfg := cgen()
			defer func() {
				// cleanCfg doesn't really have tracing turned on, so clear
				// the flag to avoid panic on cleanup.
				cleanCfg.TestHelper("set-do-trace", []string{"false"})
				cleanCfg.Cleanup()
			}()

			logs := cleanCfg.TestHelper("get-loggers", nil).([3]interface{})
			structuredLogger := logs[0].(greenery.MakeLogger)
			prettyLogger := logs[1].(greenery.MakeLogger)
			traceLogger := logs[2].(greenery.MakeTraceLogger)

			if tc.OverrideStructuredLogger {
				structuredLogger = tc.StructuredLogger
			} else {
				if global.OverrideStructuredLogger {
					structuredLogger = global.StructuredLogger
				}
			}
			if tc.OverridePrettyLogger {
				prettyLogger = tc.PrettyLogger
			} else {
				if global.OverridePrettyLogger {
					prettyLogger = global.PrettyLogger
				}

			}

			if tc.OverrideTraceLogger {
				traceLogger = tc.TraceLogger
			} else {
				if global.OverrideTraceLogger {
					traceLogger = global.TraceLogger
				}
			}

			require.NoError(t, cleanCfg.SetLoggers(structuredLogger, prettyLogger, traceLogger))
			cleanCfg.SetFs(tc.af)
			cfg := cgen()
			defer cfg.Cleanup()
			require.NoError(t, cfg.SetLoggers(structuredLogger, prettyLogger, traceLogger))
			cfg.SetFs(tc.af)

			if global.OverrideHandlerMap || tc.OverrideHandlerMap {
				fmap := cleanCfg.TestHelper("get-fmap", nil).(map[string]greenery.Handler)

				nfmap := map[string]greenery.Handler{}
				for kk, vv := range fmap {
					nfmap[kk] = vv
				}

				if global.OverrideHandlerMap {
					for kk, vv := range global.HandlerMap {
						nfmap[kk] = vv
					}
				}

				if tc.OverrideHandlerMap {
					for kk, vv := range tc.HandlerMap {
						nfmap[kk] = vv
					}
				}

				cleanCfg.TestHelper("set-fmap", nfmap)
				cfg.TestHelper("set-fmap", nfmap)
			}

			if tc.CustomParser != nil {
				cfg.RegisterExtraParse(tc.CustomParser, tc.CustomVars)
			}

			if tc.ConfigDefaults != nil {
				require.NoError(t, cleanCfg.SetOptions(*tc.ConfigDefaults))
				require.NoError(t, cfg.SetOptions(*tc.ConfigDefaults))
			}

			// If the user has given us docs with "" it means they want them
			// to be applied to the default language
			if empty, ok := userDocList[""]; ok {
				userDocList[cfg.GetDefaultLanguage()] = empty
			}

			// For debugging purposes, turn on tracing if asked
			if *ForceTrace || tc.Trace {
				cfg.StartTracing()
				cleanCfg.TestHelper("set-do-trace", []string{"true"})
			}

			// clean our environment and set any environment test variables.
			os.Clearenv()
			for k, v := range tc.Env {
				err = os.Setenv(k, v)
				require.NoError(t, err)
			}

			// If the user is passing us a config file, create it and pass it
			// in the test command line
			var cmdLine []string
			if tc.CfgContents != "" || tc.CfgFile != "" {
				var f string
				cts := tc.CfgContents
				if tc.CfgFile != "" {
					var ctsb []byte
					ctsb, err = ioutil.ReadFile(tc.CfgFile)
					require.NoError(t, err)
					cts = string(ctsb)
				}

				f, err = TempFileT(t, tc.af, "tcfg", ".toml", cts, removeList, mtx, tc.RealFilesystem)
				require.NoError(t, err)

				if tc.CmdlineCfgName != "" {
					// Typically used for failure unit tests
					f = tc.CmdlineCfgName
				}

				cmdLine = []string{
					"-c",
					f,
				}

				cmdLine = append(cmdLine, tc.CmdLine...)
				cleanCfg.TestHelper("set-cfg-file", []string{f})
			} else {
				cmdLine = tc.CmdLine
			}

			// Always log to make it easier to validate stdout / stderr /
			// logging
			var logFile string
			logFile, cmdLine = setLogging(t, &tc, cmdLine, removeList, mtx)
			cleanCfg.TestHelper("set-log-file", []string{logFile})

			// Set up the final command line
			cfg.TestHelper("set-root-args", cmdLine)

			// Start the stdout/stderr grabbers if needed
			grabberOut := NewGrabber()
			grabberErr := NewGrabber()
			if !tc.Parallel {
				require.NoError(t, grabberOut.Start(&os.Stdout))
				defer func() {
					_, _ = grabberOut.Stop()
				}()

				require.NoError(t, grabberErr.Start(&os.Stderr))
				defer func() {
					_, _ = grabberErr.Stop()
				}()
			}

			overrideMap := map[string]greenery.OverrideHandler{
				"root":              greenery.OverrideRootHandler,
				"version":           greenery.OverrideVersionHandler,
				"config":            greenery.OverrideConfigHandler,
				"-pre-exec-handler": greenery.OverridePreExecHandler,
			}

			// Override the handlers if the user requested this
			if global.OverrideBuiltinHandlers {
				for kk, vv := range global.BuiltinHandlers {
					ch, ok := overrideMap[kk]
					require.True(t, ok, "Invalid override handler %s", kk)
					err = cfg.SetHandler(ch, vv)
					require.NoError(t, err)
				}

			}

			if tc.OverrideBuiltinHandlers {
				for kk, vv := range tc.BuiltinHandlers {
					ch, ok := overrideMap[kk]
					require.True(t, ok, "Invalid override handler %s", kk)
					err = cfg.SetHandler(ch, vv)
					require.NoError(t, err)
				}

			}

			// Grab the clean configuration values in a map to make it easier
			// to compare later on
			cleanValues := getValues(cleanCfg)

			// And go
			t.Logf("Calling execute with cmd line %v", cmdLine)
			execErr := cfg.Execute(cfg, userDocList)

			// Validate any errors expected in Execute
			execHadErr := execErr != nil
			if tc.ExecError != "" {
				require.True(t, execErr != nil, "Was expecting error substring %s, but Execute operated successfully", tc.ExecError)
				require.True(t, strings.Contains(execErr.Error(), tc.ExecError), "Error \"%s\" did not contain the error substring \"%s\"", execErr.Error(), tc.ExecError)
				execErr = nil
			}

			if tc.ExecErrorRegex != "" {
				require.True(t, execErr != nil, "Was expecting error regex %s, but Execute operated successfully", tc.ExecErrorRegex)
				require.True(t, regexpMatches(tc.ExecErrorRegex, execErr.Error()), "Error \"%s\" did not match error regex \"%s\"", execErr.Error(), tc.ExecErrorRegex)
				execErr = nil
			}

			require.NoError(t, execErr)

			if execHadErr && !tc.ExecErrorOutput {
				return
			}

			// Validate stdout/stderr
			if !tc.Parallel {
				outS, errOut := grabberOut.Stop()
				errS, errErr := grabberErr.Stop()

				require.NoError(t, errOut)
				require.NoError(t, errErr)

				if tc.Trace {
					t.Logf("Tracing stdout: %s\n", outS)
					t.Logf("Tracing stderr: %s\n", errS)
				} else {
					verifyStdOutErr(t, &tc, outS, errS)
				}
			}

			// Validate the logfile if required
			if !tc.pretty {
				if tc.GoldLog != "" {
					t.Log("Calling verify gold logs")
					verifyGoldLog(t, &tc, logFile)
				}

				if tc.OutLogLines != nil {
					t.Log("Calling verify logs")
					verifyLog(t, tc.af, logFile, tc.OutLogLines, tc.OutLogAllLines)
				}
			}

			if tc.OutLogRegex != "" {
				t.Log("Base log regex")
				logB, err := afero.ReadFile(tc.af, logFile)
				require.NoError(t, err)
				require.True(t, regexpMatches(tc.OutLogRegex, string(logB)), "Log output \"%s\" did not match log regex \"%s\"", string(logB), tc.OutLogRegex)
			}

			verifyGoldFiles(t, &tc)

			if tc.ValuesValidator != nil {
				t.Logf("Test has a validator, calling it")
				tc.ValuesValidator(t, cfg)

				// All done
				return
			}

			if tc.NoValidateConfigValues {
				return
			}

			// If there is no validate check for the correct values, the set
			// ones set to what they should be, the others set to the default
			// (whatever it is) in the "clean" configuration.
			verifyValues(t, &tc, cfg, cleanCfg, cleanValues, compFuncs)

		})
	}

	return nil
}
