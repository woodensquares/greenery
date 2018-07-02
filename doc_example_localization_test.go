package greenery_test

import (
	"fmt"
	"log"
	"os"

	"github.com/woodensquares/greenery"
)

type exampleLocalizedConfig struct {
	*greenery.BaseConfig
	Timeout *greenery.IntValue `greenery:"get|timeout|t,      localized-app.timeout,   TIMEOUT"`
}

func exampleNewLocalizedConfig() *exampleLocalizedConfig {
	cfg := &exampleLocalizedConfig{
		BaseConfig: greenery.NewBaseConfig("localized", map[string]greenery.Handler{
			"get<": exampleLocalizedGetter,
		}),
		Timeout: greenery.NewIntValue("Timeout", 0, 1000),
	}

	if err := cfg.Timeout.SetInt(400); err != nil {
		panic("Could not initialize the timeout to its default")
	}

	return cfg
}

// Sample localization in pig latin, using the []strings format for
// compactness and ease of editing even for non-golang developers
var exampleLocalizedPiglatin = []string{
	"1.0",
	greenery.DocBaseDelimiter,
	// Root command values
	greenery.DocShort,
	"RLUay etcherfay",
	greenery.DocLong,
	"RLUay etcherfay ongerlay escriptionday",
	greenery.DocExample,
	"    Eesay ethay etgay ommandcay orfay examplesay",
	greenery.DocHelpFlag,
	"Elphay informationay orfay ethay applicationay.",
	greenery.DocCmdFlags,
	"[lagsfay]",
	greenery.DocConfigEnvMsg1,
	`Ethay ollowingfay environmentay ariablesvay areay activeay anday ouldcay affectay ethay executionay
ofay ethay ogrampray ependingday onay ommandcay inelay argumentsay:`,
	greenery.DocConfigEnvMsg2,
	"Onay ariablesvay",
	greenery.DocConfigEnvMsg3,
	"Ethay ollowingfay environmentay ariablesvay areay availableay orfay isthay ogrampray:",
	greenery.DocBaseDelimiter,

	// General help template strings
	greenery.DocHelpDelimiter,
	greenery.DocUsage,
	"Sageuay:",
	greenery.DocAliases,
	"Liasesaay:",
	greenery.DocExamples,
	"Ampleseay:",
	greenery.DocAvailableCommands,
	"Vailableaay Ommandscay:",
	greenery.DocFlags,
	"Lagsfay:",
	greenery.DocGlobalFlags,
	"Lobalgay Lagsfay:",
	greenery.DocAdditionalHelpTopics,
	"Dditionalaay elphay opicstay:",
	greenery.DocProvidesMoreInformationAboutACommand,
	"povidespay oremay informationay aboutay aay ommandcay.",
	greenery.DocHelpDelimiter,

	// Individual command usage strings, system first
	greenery.DocUsageDelimiter,
	greenery.DocHelpCmd,
	"[ommandcay]",
	"Elphay aboutay anyay command",
	`Elphay ovidesrpay helpay orfay anyay ommandcay inay ethay applicationay.
Implysay ypetay appnameay elphay [path otay ommandcay] orfay ullfay etailsday.`,
	"",
	greenery.DocConfigCmd,
	"",
	"Onfigurationcay ilefay elatedray ommandscay",
	"",
	"",
	greenery.DocConfigInitCmd,
	"",
	"Reatescay aay efaultday onfigcay ilefay inay wdcay oray hereway -c isay etsay",
	`Enwhay histay ommandcay isay executeday aay onfigurationcay ilefay illway ebay reatedcay
"inay hetay pecifiedsay ocationlay (in hetay urrentcay irectoryday ybay default", oray overnedgay
"ybay hetay aluevay assedpay otay hetay wdcay flag). Histay onfigurationcay ilefay illway ontaincay
"allay hetay upportedsay ariablesvay ithway heirtay efaultday values`,
	"",
	greenery.DocConfigEnvCmd,
	"",
	"Howssay hetay activeay environmentay ariablesvay hattay ouldway impactay hetay program",
	"",
	"",
	greenery.DocConfigDisplayCmd,
	"",
	"Howssay hetay urrentcay onfigurationcay values",
	`Howsay hetay urrentcay onfigurationcay aluesvay akingtay intoay accountay allay environmentay
"anday onfigurationcay ilefay values. Command-line agsflay invaliday orfay hetay onfigcay isplayday
"ommandcay ouldway otnay ebay displayed`,
	"",
	greenery.DocVersionCmd,
	"",
	"Rintspay outay hetay ersionvay umbernay ofay hetay ogramrpay",
	"",
	"",

	// Our command
	"get",
	"[UEI otay fetch]",
	"Etrievesray hetay pecifiedsay agepay",
	"",
	`    Otay etrieveray a agepay availableay onay localhost atay /hello.html implysay unray

    localized get http://localhost/hello.html`,
	greenery.DocUsageDelimiter,

	// Cmdline variable descriptions, system first
	greenery.DocCmdlineDelimiter,
	greenery.DocLogLevel,
	"Hetay oglay evellay ofay hetay ogrampray. Alidvay aluesvay areay \"error\", \"warn\", \"info\" anday \"debug\"",
	greenery.DocConfFile,
	"Hetay onfigurationcay ilefay ocationlay",
	greenery.DocLogFile,
	"Hetay oglay ilefay ocationlay",
	greenery.DocPretty,
	"Fiay etsay hetay onsolecay outputay ofay hetay ogginglay allscay illway ebay ettifiedpray",
	greenery.DocNoCfg,
	"Fiay etsay onay onfigurationay ilefay illway ebay oadedlay",
	greenery.DocNoEnv,
	"Fiay etsay hetay environmentay ariablesvay illway otnay ebay onsideredcay",
	greenery.DocVerbosity,
	"Hetay erbosityvay ofay hetay ogrampray, anay integeray etweenbay 0 anday 3 inclusiveay.",
	greenery.DocDoTrace,
	"Enablesay acingtray",
	greenery.DocCfgLocation,
	"Hereway otay itewray hetay onfigurationcay ilefay, oneay ofay \"cwd\", \"user\" oray \"system\"",
	greenery.DocCfgForce,
	"Fiay specified, anyay existingay onfigurationcay ilesfay illway ebay overwrittenay",
	// our flag
	"Timeout",
	"Hetay imeouttay otay useay orfay hetay ETGay operationay",
	greenery.DocCmdlineDelimiter,

	// Config file variable descriptions (where different from the cmdline)
	greenery.DocConfigDelimiter,
	greenery.DocConfigHeader,
	"Onfigurationcay eneratedgay onay {{ date }}",
	greenery.DocConfigDelimiter,

	// Custom message
	greenery.DocCustomDelimiter,
	"Message",
	"Illway etchfay %s ithway imeouttay %d illisecondsmay\n",
	greenery.DocCustomDelimiter,
}

// Same in English, this will of course contain only values specific to our
// application.
var exampleLocalizedEnglish = []string{
	"1.0",
	greenery.DocBaseDelimiter,

	// Root command
	greenery.DocShort,
	"URL fetcher",
	greenery.DocExample,
	"    See the get command for examples",
	greenery.DocBaseDelimiter,

	// Additional commands, get in this case
	greenery.DocUsageDelimiter,
	"get",
	"[URI to fetch]",
	"Retrieves the specified page",
	"",
	`    To retrieve a page available on localhost at /hello.html simply run

    localized get http://localhost/hello.html`,
	greenery.DocUsageDelimiter,

	// Cmdline flag documentation
	greenery.DocCmdlineDelimiter,
	"Timeout",
	"the timeout, in milliseconds, to use for the fetch",
	greenery.DocCmdlineDelimiter,

	// Custom message
	greenery.DocCustomDelimiter,
	"Message",
	"Will fetch %s with timeout %d milliseconds\n",
	greenery.DocCustomDelimiter,
}

var exampleLocalizedDocs = greenery.ConvertOrPanic(map[string][]string{
	"en":     exampleLocalizedEnglish,
	"zz.ZZZ": exampleLocalizedPiglatin,
})

func exampleLocalizedGetter(lcfg greenery.Config, args []string) error {
	cfg := lcfg.(*exampleLocalizedConfig)

	if len(args) != 1 {
		return fmt.Errorf("Invalid number of command line arguments")
	}

	cfg.Debugf("fetching %s, timeout %d", args[0], cfg.Timeout.Value)
	_, docs := cfg.GetDocs()
	fmt.Printf(docs.Custom["Message"], args[0], cfg.Timeout.Value)
	return nil
}

func exampleLocalizedMain() {
	cfg := exampleNewLocalizedConfig()
	defer cfg.Cleanup()

	cfg.VersionMajor = "1"
	cfg.VersionMinor = "0"

	if err := cfg.Execute(cfg, exampleLocalizedDocs); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing: %s\n", err)
		os.Exit(1)
	}
}

// Example_localized is just like Example_minimal, but supporting a fully
// localized additional language, pig latin, if called with LANG set to
// zz.ZZZ. A sample invocation of the help message setting the LANG
// environment variable to two different values is provided.
func Example_localized() {
	fmt.Println("----------------------------------------------------------")
	if err := os.Setenv("LANG", "en_US.UTF-8"); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Default help for the get command\n\n")
	os.Args = []string{"localized", "get", "-h"}
	exampleLocalizedMain()

	fmt.Printf("Default command output\n\n")
	os.Args = []string{"localized", "get", "http://127.0.0.1"}
	exampleLocalizedMain()

	if err := os.Setenv("LANG", "zz.ZZZ"); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("----------------------------------------------------------")
	fmt.Printf("\nLocalized help for the get command\n\n")
	os.Args = []string{"localized", "get", "-h"}
	exampleLocalizedMain()

	fmt.Printf("Localized command output\n\n")
	os.Args = []string{"localized", "get", "http://127.0.0.1"}
	exampleLocalizedMain()

	// Output: ----------------------------------------------------------
	// Default help for the get command
	//
	// Retrieves the specified page
	//
	// Usage:
	//   localized get [URI to fetch] [flags]
	//
	// Examples:
	//     To retrieve a page available on localhost at /hello.html simply run
	//
	//     localized get http://localhost/hello.html
	//
	// Flags:
	//   -t, --timeout int   the timeout, in milliseconds, to use for the fetch (default 400)
	//
	// Global Flags:
	//   -c, --config string      The configuration file location
	//       --help               help information for the application.
	//       --log-file string    The log file location
	//   -l, --log-level string   The log level of the program. Valid values are "error", "warn", "info" and "debug" (default "error")
	//       --no-cfg             If set no configuration file will be loaded
	//       --no-env             If set the environment variables will not be considered
	//       --pretty             If set the console output of the logging calls will be prettified
	//   -v, --verbosity int      The verbosity of the program, an integer between 0 and 3 inclusive. (default 1)
	//
	// Default command output
	//
	// Will fetch http://127.0.0.1 with timeout 400 milliseconds
	// ----------------------------------------------------------
	//
	// Localized help for the get command
	//
	// Etrievesray hetay pecifiedsay agepay
	//
	// Sageuay:
	//   localized get [UEI otay fetch] [lagsfay]
	//
	// Ampleseay:
	//     Otay etrieveray a agepay availableay onay localhost atay /hello.html implysay unray
	//
	//     localized get http://localhost/hello.html
	//
	// Lagsfay:
	//   -t, --timeout int   Hetay imeouttay otay useay orfay hetay ETGay operationay (default 400)
	//
	// Lobalgay Lagsfay:
	//   -c, --config string      Hetay onfigurationcay ilefay ocationlay
	//       --help               Elphay informationay orfay ethay applicationay.
	//       --log-file string    Hetay oglay ilefay ocationlay
	//   -l, --log-level string   Hetay oglay evellay ofay hetay ogrampray. Alidvay aluesvay areay "error", "warn", "info" anday "debug" (default "error")
	//       --no-cfg             Fiay etsay onay onfigurationay ilefay illway ebay oadedlay
	//       --no-env             Fiay etsay hetay environmentay ariablesvay illway otnay ebay onsideredcay
	//       --pretty             Fiay etsay hetay onsolecay outputay ofay hetay ogginglay allscay illway ebay ettifiedpray
	//   -v, --verbosity int      Hetay erbosityvay ofay hetay ogrampray, anay integeray etweenbay 0 anday 3 inclusiveay. (default 1)
	//
	// Localized command output
	//
	// Illway etchfay http://127.0.0.1 ithway imeouttay 400 illisecondsmay
}
