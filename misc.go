package greenery

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// LogField is used by the structured logging functions to pass around
// information about the specific field. Each logger can decide how to use
// these fields as they see fit.
type LogField struct {
	Key      string
	Type     uint8
	Integer  int
	Time     time.Time
	Duration time.Duration
	String   string
	Generic  interface{}
}

// MakeLogger describes a function that would return a normal logger
type MakeLogger func(Config, string, io.Writer) Logger

// MakeTraceLogger describes a function that would return a trace logger
type MakeTraceLogger func(Config) Logger

// Logger is an interface that describes a logger provided by logger.go,
// currently this is just a basic printf logger. A zap logger is available in
// zaplogger/ if a structured or more high performant logger is required.
type Logger interface {
	Custom(string, int, interface{})
	DebugStructured(string, ...LogField)
	InfoStructured(string, ...LogField)
	WarnStructured(string, ...LogField)
	ErrorStructured(string, ...LogField)

	LogString(string, string) LogField
	LogInteger(string, int) LogField
	LogTime(string, time.Time) LogField
	LogDuration(string, time.Duration) LogField
	LogGeneric(string, interface{}) LogField

	Sync() error

	// For testing purposes primarily
	DebugSkip(int, string)
}

// --------------------------------------------------------------------------
// the root command identifier as used by the command maps
const rootCommandID = ""

// Used for parsing our tag, note that sepKeyParts should not be changed as
// due to https://github.com/spf13/viper/pull/399 not being merged yet, we
// can't do things like for example having | for cfg variables as well
const sepKeyParts = "."
const sepCmdParts = "|"
const sepTag = ","
const sepCmdLevels = ">"
const sepCmdArgs = "<"
const sepMultipleCmds = "&"

// process will take an initialized configuration and do anything that needs
// to be done to make it ready for execution by the current command. Currently
// it's only setting up the normal/pretty logging. This is called after the
// configuration is loaded so all variables in cfg can be assumed to be
// correct.
func (cfg *BaseConfig) process() error {
	if cfg.LogFile != "" {
		f, err := cfg.s_fs.OpenFile(cfg.LogFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return errors.WithMessage(err,
				fmt.Sprintf("Cannot open log file %s", cfg.LogFile))
		}

		cfg.s_filesToClose = append(cfg.s_filesToClose, f)
		cfg.s_w = f
	}

	if cfg.Pretty {
		if cfg.s_makePretty != nil {
			cfg.s_log = cfg.s_makePretty(cfg, cfg.LogLevel.Value, cfg.s_w)
		}
	} else {
		if cfg.s_makeStructured != nil {
			cfg.s_log = cfg.s_makeStructured(cfg, cfg.LogLevel.Value, cfg.s_w)
		}
	}

	cfg.s_processed = true
	return nil
}

// runWrapper executes the user command after loading the configuration and
// processing it. It is possible for commands to run some last-minute
// initializations by setting a pre-exec handler.
func runWrapper(ccmd *cobra.Command,
	cfg *BaseConfig,
	xcmd func(Config, []string) error,
	args []string) (err error) {

	// NoEnv is only a cmdline parameter, and will be set already, if it is
	// let's remove our environmental variables
	if cfg.NoEnv {
		for k := range cfg.s_env {
			cfg.Tracef("Clearing environment variable %s", k)
			if err = os.Setenv(k, ""); err != nil {
				return
			}
		}
	}

	// All our command parameters, environment and config end up in cfg, so no
	// need to pass on cmd or args[], which aren't given to the wrapper
	err = cfg.load(cfg.s_cl, cfg.s_appName+".toml", ccmd, cfg.s_v)
	if err != nil {
		return
	}

	var name func(*cobra.Command) string
	name = func(cc *cobra.Command) string {
		if cc.HasParent() {
			return name(cc.Parent()) + sepCmdLevels + cc.Name()
		}
		return ""
	}
	cfg.s_currentcmd = strings.TrimLeft(name(ccmd), sepCmdLevels)

	if err = cfg.process(); err != nil {
		return
	}

	if cfg.s_preExecHandler != nil {
		if err = cfg.s_preExecHandler(cfg.s_cl, args); err != nil {
			return
		}
	}

	cfg.Trace("runWrapper end, calling the command")
	return xcmd(cfg.s_cl, args)
}

// getParentChild is a convenience function to get the parent/child names of
// the command including the root command.
func getParentChild(k string) (string, string) {
	parent := rootCommandID
	child := strings.Split(k, sepCmdLevels)
	cl := len(child)
	var cname string
	if cl > 1 {
		parent = strings.Join(child[:cl-1], sepCmdLevels)
		cname = child[cl-1]
	} else {
		cname = k
	}

	return parent, cname
}
