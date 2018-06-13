package greenery_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/woodensquares/greenery"
)

type exampleMinimalConfig struct {
	*greenery.BaseConfig
	Timeout *greenery.IntValue `greenery:"get|timeout|t,      minimal-app.timeout,   TIMEOUT"`
}

func exampleNewMinimalConfig() *exampleMinimalConfig {
	cfg := &exampleMinimalConfig{
		BaseConfig: greenery.NewBaseConfig("minimal", map[string]greenery.Handler{
			"get<": exampleMinimalGetter,
		}),
		Timeout: greenery.NewIntValue("Timeout", 0, 1000),
	}

	if err := cfg.Timeout.SetInt(400); err != nil {
		panic("Could not initialize the timeout to its default")
	}

	return cfg
}

var exampleMinimalDocs = map[string]*greenery.DocSet{
	"en": &greenery.DocSet{
		Short: "URL fetcher",
		Usage: map[string]*greenery.CmdHelp{
			"get": &greenery.CmdHelp{
				Use:   "[URI to fetch]",
				Short: "Retrieves the specified page",
			},
		},
		CmdLine: map[string]string{
			"Timeout": "the timeout to use for the fetch",
		},
	},
}

func exampleMinimalGetter(lcfg greenery.Config, args []string) error {
	cfg := lcfg.(*exampleMinimalConfig)

	if len(args) != 1 {
		return fmt.Errorf("Invalid number of command line arguments")
	}

	cfg.Debugf("fetching %s, timeout %d", args[0], cfg.Timeout.Value)
	fmt.Printf("Will fetch %s with timeout %d milliseconds\n\n", args[0], cfg.Timeout.Value)
	return nil
}

func exampleMinimalMain() {
	cfg := exampleNewMinimalConfig()
	defer cfg.Cleanup()

	if err := cfg.Execute(cfg, exampleMinimalDocs); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing: %s\n", err)
		os.Exit(1)
	}
}

// Example_minimal shows a minimal application using the greery
// library. Invocations mimicking actual command-line usage setting the
// parameter value in different ways are provided for reference.
func Example_minimal() {
	fmt.Println("----------------------------------------------------------")

	// Execute with no arguments, will use the default timeout
	os.Args = []string{"minimal", "get", "http://127.0.0.1"}
	fmt.Println("Executing with the default timeout")
	exampleMinimalMain()

	// Execute with a commandline timeout argument
	os.Args = []string{"minimal", "get", "-t", "300", "http://127.0.0.1"}
	fmt.Println("Executing with a commandline timeout")
	exampleMinimalMain()

	// Create a configuration file with a timeout value
	cfgFile, err := ioutil.TempFile("", "minimal")
	if err != nil {
		fmt.Println("Cannot create a temporary file")
		os.Exit(1)
	}
	defer func() {
		_ = os.Remove(cfgFile.Name())
	}()

	if _, err = cfgFile.Write([]byte(`# Sample config file
[minimal-app]
timeout = 550
`)); err != nil {
		fmt.Println("Cannot write the config file")
		os.Exit(1)
	}
	if err := cfgFile.Close(); err != nil {
		fmt.Println("Cannot close the config file")
		os.Exit(1)
	}

	// Execute with the created configuration file
	os.Args = []string{"minimal", "get", "-c", cfgFile.Name(), "http://127.0.0.1"}
	fmt.Println("Executing with a configuration file")
	exampleMinimalMain()

	// Execute with an environment variable argument, together with the
	// configuration file, note the environment variable will have precedence
	if err := os.Setenv("MINIMAL_TIMEOUT", "200"); err != nil {
		fmt.Println("Cannot set an environment variable")
		os.Exit(1)
	}

	os.Args = []string{"minimal", "get", "-c", cfgFile.Name(), "http://127.0.0.1"}
	fmt.Println("Executing with a configuration and environment variable timeout")
	exampleMinimalMain()

	// Executing with commandline, config and environment timeout argument,
	// the commandline will take precedence
	os.Args = []string{"minimal", "get", "-c", cfgFile.Name(), "-t", "300", "http://127.0.0.1"}
	fmt.Println("Executing with a commandline, configuration and environment timeout")
	exampleMinimalMain()

	// Output: ----------------------------------------------------------
	// Executing with the default timeout
	// Will fetch http://127.0.0.1 with timeout 400 milliseconds
	//
	// Executing with a commandline timeout
	// Will fetch http://127.0.0.1 with timeout 300 milliseconds
	//
	// Executing with a configuration file
	// Will fetch http://127.0.0.1 with timeout 550 milliseconds
	//
	// Executing with a configuration and environment variable timeout
	// Will fetch http://127.0.0.1 with timeout 200 milliseconds
	//
	// Executing with a commandline, configuration and environment timeout
	// Will fetch http://127.0.0.1 with timeout 300 milliseconds
	//
}
