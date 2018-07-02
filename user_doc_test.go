package greenery_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/woodensquares/greenery"
	"github.com/woodensquares/greenery/testhelper"
)

var rootHelp = filepath.Join("testdata", "language_test.common.rootHelp.stdout")
var rootHelpIt = filepath.Join("testdata", "language_test.common.rootHelp.it.stdout")
var rootHelpPiglatin = filepath.Join("testdata", "language_test.common.rootHelp.piglatin.stdout")
var configHelp = filepath.Join("testdata", "language_test.common.configHelp.stdout")
var configHelpIt = filepath.Join("testdata", "language_test.common.configHelp.it.stdout")
var configHelpPiglatin = filepath.Join("testdata", "language_test.common.configHelp.piglatin.stdout")
var configInitHelp = filepath.Join("testdata", "language_test.common.configInitHelp.stdout")
var configInitHelpIt = filepath.Join("testdata", "language_test.common.configInitHelp.it.stdout")
var configInitHelpPiglatin = filepath.Join("testdata", "language_test.common.configInitHelp.piglatin.stdout")
var configEnvHelp = filepath.Join("testdata", "language_test.common.configEnvHelp.stdout")
var configEnvHelpIt = filepath.Join("testdata", "language_test.common.configEnvHelp.it.stdout")
var configEnvHelpPiglatin = filepath.Join("testdata", "language_test.common.configEnvHelp.piglatin.stdout")
var configEnvOut = filepath.Join("testdata", "language_test.common.configEnvOut.stdout")
var configEnvOutIt = filepath.Join("testdata", "language_test.common.configEnvOut.it.stdout")
var configEnvOutPiglatin = filepath.Join("testdata", "language_test.common.configEnvOut.piglatin.stdout")
var configDisplayHelp = filepath.Join("testdata", "language_test.common.configDisplayHelp.stdout")
var configDisplayHelpIt = filepath.Join("testdata", "language_test.common.configDisplayHelp.it.stdout")
var configDisplayHelpPiglatin = filepath.Join("testdata", "language_test.common.configDisplayHelp.piglatin.stdout")
var helpHelp = filepath.Join("testdata", "language_test.common.helpHelp.stdout")
var helpHelpIt = filepath.Join("testdata", "language_test.common.helpHelp.it.stdout")
var helpHelpPiglatin = filepath.Join("testdata", "language_test.common.helpHelp.piglatin.stdout")
var versionHelp = filepath.Join("testdata", "language_test.common.versionHelp.stdout")
var versionHelpIt = filepath.Join("testdata", "language_test.common.versionHelp.it.stdout")
var versionHelpPiglatin = filepath.Join("testdata", "language_test.common.versionHelp.piglatin.stdout")

// Note that this list will be using strings, instead of constants, to ensure
// backwards compatibility (if any strings change the test will break).
var piglatin = []string{
	"1.0",
	"------ DELIMITER:BASE ------", // greenery.DocBaseDelimiter
	"Short", // greenery.DocShort
	"RLUay etcherfay",
	"Long", // greenery.DocLong
	"RLUay etcherfay ongerlay escriptionday",
	"Example", // greenery.DocExample
	"    Eesay ethay etgay ommandcay orfay examplesay",
	"HelpFlag", // greenery.DocHelpFlag
	"Elphay informationay orfay ethay applicationay.",
	"CmdFlags", // greenery.DocCmdFlags
	"[lagsfay]",
	"ConfigEnvMsg1", // greenery.DocConfigEnvMsg1
	`Ethay ollowingfay environmentay ariablesvay areay activeay anday ouldcay affectay ethay executionay
ofay ethay ogrampray ependingday onay ommandcay inelay argumentsay:`,
	"ConfigEnvMsg2", // greenery.DocConfigEnvMsg2
	"Onay ariablesvay",
	"ConfigEnvMsg3", // greenery.DocConfigEnvMsg3
	"Ethay ollowingfay environmentay ariablesvay areay availableay orfay isthay ogrampray:",
	"------ DELIMITER:BASE ------", // greenery.DocBaseDelimiter

	"------ DELIMITER:HELP ------", // greenery.DocHelpDelimiter
	"Usage", // greenery.DocUsage
	"Sageuay:",
	"Aliases", // greenery.DocAliases
	"Liasesaay:",
	"Examples", // greenery.DocExamples
	"Ampleseay:",
	"AvailableCommands", // greenery.DocAvailableCommands
	"Vailableaay Ommandscay:",
	"Flags", // greenery.DocFlags
	"Lagsfay:",
	"GlobalFlags", // greenery.DocGlobalFlags
	"Lobalgay Lagsfay:",
	"AdditionalHelpTopics", // greenery.DocAdditionalHelpTopics
	"Dditionalaay elphay opicstay:",
	"ProvidesMoreInformationAboutACommand", // greenery.DocProvidesMoreInformationAboutACommand
	"povidespay oremay informationay aboutay aay ommandcay.",
	"------ DELIMITER:HELP ------", // greenery.DocHelpDelimiter

	"------ DELIMITER:USAGE ------", // greenery.DocUsageDelimiter
	"help", // greenery.DocHelpCmd
	"[ommandcay]",
	"Elphay aboutay anyay command",
	`Elphay ovidesrpay helpay orfay anyay ommandcay inay ethay applicationay.
Implysay ypetay appnameay elphay [path otay ommandcay] orfay ullfay etailsday.`,
	"",
	"config", // greenery.DocConfigCmd
	"",
	"Onfigurationcay ilefay elatedray ommandscay",
	"",
	"",
	"config>init", // greenery.DocConfigInitCmd
	"",
	"Reatescay aay efaultday onfigcay ilefay inay wdcay oray hereway -c isay etsay",
	`Enwhay histay ommandcay isay executeday aay onfigurationcay ilefay illway ebay reatedcay
"inay hetay pecifiedsay ocationlay (in hetay urrentcay irectoryday ybay default", oray overnedgay
"ybay hetay aluevay assedpay otay hetay wdcay flag). Histay onfigurationcay ilefay illway ontaincay
"allay hetay upportedsay ariablesvay ithway heirtay efaultday values`,
	"",
	"config>env", // greenery.DocConfigEnvCmd
	"",
	"Howssay hetay activeay environmentay ariablesvay hattay ouldway impactay hetay program",
	"",
	"",
	"config>display", // greenery.DocConfigDisplayCmd
	"",
	"Howssay hetay urrentcay onfigurationcay values",
	`Howsay hetay urrentcay onfigurationcay aluesvay akingtay intoay accountay allay environmentay
"anday onfigurationcay ilefay values. Command-line agsflay invaliday orfay hetay onfigcay isplayday
"ommandcay ouldway otnay ebay displayed`,
	"",
	"version", // greenery.DocVersionCmd
	"",
	"Rintspay outay hetay ersionvay umbernay ofay hetay ogramrpay",
	"",
	"",
	"------ DELIMITER:USAGE ------", // greenery.DocUsageDelimiter

	"------ DELIMITER:COMMANDLINE ------", // greenery.DocCmdlineDelimiter
	"LogLevel",                            // greenery.DocLogLevel
	"Hetay oglay evellay ofay hetay ogrampray. Alidvay aluesvay areay \"error\", \"warn\", \"info\" anday \"debug\"",
	"ConfFile", // greenery.DocConfFile
	"Hetay onfigurationcay ilefay ocationlay",
	"LogFile", // greenery.DocLogFile
	"Hetay oglay ilefay ocationlay",
	"Pretty", // greenery.DocPretty
	"Fiay etsay hetay onsolecay outputay ofay hetay ogginglay allscay illway ebay ettifiedpray",
	"NoCfg", // greenery.DocNoCfg
	"Fiay etsay onay onfigurationay ilefay illway ebay oadedlay",
	"NoEnv", // greenery.DocNoEnv
	"Fiay etsay hetay environmentay ariablesvay illway otnay ebay onsideredcay",
	"Verbosity", // greenery.DocVerbosity
	"Hetay erbosityvay ofay hetay ogrampray, anay integeray etweenbay 0 anday 3 inclusiveay.",
	"DoTrace", // greenery.DocDoTrace
	"Enablesay acingtray",
	"CfgLocation", // greenery.DocCfgLocation
	"Hereway otay itewray hetay onfigurationcay ilefay, oneay ofay \"cwd\", \"user\" oray \"system\"",
	"CfgForce", // greenery.DocCfgForce
	"Fiay specified, anyay existingay onfigurationcay ilesfay illway ebay overwrittenay",
	"------ DELIMITER:COMMANDLINE ------", // greenery.DocCmdlineDelimiter

	// Config file variable descriptions (where different from the cmdline)
	"------ DELIMITER:CONFIG ------", // greenery.DocConfigDelimiter
	".", // greenery.DocConfigHeader
	"Autogenerateday cnfigurationcay ilefay",
	"------ DELIMITER:CONFIG ------", // greenery.DocConfigDelimiter
}

var badRoot = []string{
	"1.0",
	greenery.DocUsageDelimiter,
	"",
	"",
	"root short",
	"root long",
	"",
}

func TestBadRootSection(t *testing.T) {
	badSet, err := greenery.ConvertDocs(badRoot)
	require.NoError(t, err)

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "root help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"-h",
			},
			ExecError: "Please use the root command section to specify root command options, not the usage one",
		},
	}

	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
		UserDocList: map[string]*greenery.DocSet{
			"": badSet,
		},
	})
	require.NoError(t, err)
}

func TestUserLang(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "root help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "config help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"config",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "config init help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"config",
				"init",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configInitHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "config env help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"config",
				"env",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configEnvHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "config env output user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"config",
				"env",
			},
			GoldStdOut: &testhelper.TestFile{Source: configEnvOutPiglatin},
		},

		testhelper.TestCase{
			Name: "config display help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"config",
				"display",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configDisplayHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"help",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: helpHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "version help user language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"version",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: versionHelpPiglatin},
		},

		testhelper.TestCase{
			Name: "Localized config",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			CmdLine: []string{
				"config",
				"init",
			},
			GoldFiles: []testhelper.TestFile{
				testhelper.TestFile{Location: "simple.toml",
					Source: filepath.Join("testdata", "user_doc_test.TestUserLang.langcfg"),
					Perms:  0644,
					Custom: testhelper.CompareIgnoreTmp},
			},
			OutStdOutRegex: "^Configuration file generated at",
		},
	}

	piglatinSet, err := greenery.ConvertDocs(piglatin)
	require.NoError(t, err)

	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
		UserDocList: map[string]*greenery.DocSet{
			"zz_ZZ.UTF-8": piglatinSet,
		},
	})
	require.NoError(t, err)
}

func TestLang(t *testing.T) {
	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "empty lang, no user docs",
			ValuesValidator: func(t *testing.T, icfg greenery.Config) {
				lang, _ := icfg.GetDocs()
				require.EqualValues(t, ",en,", lang)
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "lang present, same as default, no user docs, root",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			ValuesValidator: func(t *testing.T, icfg greenery.Config) {
				lang, _ := icfg.GetDocs()
				require.EqualValues(t, "en_US.UTF-8,en,", lang)
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "lang present, other than the default, no user docs, root",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			ValuesValidator: func(t *testing.T, icfg greenery.Config) {
				lang, _ := icfg.GetDocs()
				require.EqualValues(t, "it_IT.UTF-8,it,", lang)
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelpIt},
		},
		testhelper.TestCase{
			Name: "lang not present, no user docs",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			ValuesValidator: func(t *testing.T, icfg greenery.Config) {
				lang, _ := icfg.GetDocs()
				require.EqualValues(t, "zz_ZZ.UTF-8,en,", lang)
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "lang not present, no user docs, different default language",
			Env: map[string]string{
				"LANG": "zz_ZZ.UTF-8",
			},
			ValuesValidator: func(t *testing.T, icfg greenery.Config) {
				lang, _ := icfg.GetDocs()
				require.EqualValues(t, "zz_ZZ.UTF-8,it,", lang)
			},
			ConfigDefaults: &greenery.BaseConfigOptions{
				DefaultLanguage: "it",
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelpIt},
		},

		// Help validation
		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "root help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelp},
		},
		testhelper.TestCase{
			Name: "root help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: rootHelpIt},
		},

		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "config help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"config",
			},
			GoldStdOut: &testhelper.TestFile{Source: configHelp},
		},
		testhelper.TestCase{
			Name: "config help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configHelpIt},
		},

		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "config init help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"config",
				"init",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configInitHelp},
		},
		testhelper.TestCase{
			Name: "config init help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"init",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configInitHelpIt},
		},

		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "config env help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"config",
				"env",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configEnvHelp},
		},

		testhelper.TestCase{
			Name: "config env output default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"config",
				"env",
			},
			GoldStdOut: &testhelper.TestFile{Source: configEnvOut},
		},
		testhelper.TestCase{
			Name: "config env help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"env",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configEnvHelpIt},
		},

		testhelper.TestCase{
			Name: "config env output other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"env",
			},
			GoldStdOut: &testhelper.TestFile{Source: configEnvOutIt},
		},

		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "config display help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"config",
				"display",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configDisplayHelp},
		},
		testhelper.TestCase{
			Name: "config display help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"display",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: configDisplayHelpIt},
		},

		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"help",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: helpHelp},
		},
		testhelper.TestCase{
			Name: "help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"help",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: helpHelpIt},
		},

		// ------------------------------------------------------------------
		testhelper.TestCase{
			Name: "version help default language",
			Env: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			CmdLine: []string{
				"version",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: versionHelp},
		},
		testhelper.TestCase{
			Name: "version help other language",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"version",
				"-h",
			},
			GoldStdOut: &testhelper.TestFile{Source: versionHelpIt},
		},
	}

	err := testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
	})
	require.NoError(t, err)
}

// Test adding arguments to root / config / version
func TestCommandsAdditional(t *testing.T) {
	enSet, err := greenery.ConvertDocs([]string{
		"1.0",
		greenery.DocBaseDelimiter,
		greenery.DocUse,
		"[RFORMAT]",
		greenery.DocExample,
		"example information about the various formats in general",
		greenery.DocBaseDelimiter,
		greenery.DocUsageDelimiter,
		greenery.DocVersionCmd,
		"[VFORMAT]",
		"",
		"",
		"example information about the various formats for the version",
		greenery.DocConfigCmd,
		"[CFORMAT]",
		"",
		"",
		"example information about the various formats for the config",
		greenery.DocHelpCmd,
		"",
		"",
		"",
		"example information about using the help command",
		greenery.DocUsageDelimiter,
		greenery.DocCustomDelimiter,
		"FormatString",
		"Format %s was requested\n",
		greenery.DocCustomDelimiter,
	})
	require.NoError(t, err)

	itSet, err := greenery.ConvertDocs([]string{
		"1.0",
		greenery.DocBaseDelimiter,
		greenery.DocUse,
		"[FORMATOR]",
		greenery.DocExample,
		"informazioni sui formati in genere",
		greenery.DocBaseDelimiter,
		greenery.DocUsageDelimiter,
		greenery.DocVersionCmd,
		"[FORMATOV]",
		"",
		"",
		"informazioni sui formati per la versione",
		greenery.DocConfigCmd,
		"[FORMATOC]",
		"",
		"",
		"informazioni sui formati per la configurazione",
		greenery.DocHelpCmd,
		"",
		"",
		"",
		"informazioni sul comando di aiuto",
		greenery.DocUsageDelimiter,
		greenery.DocCustomDelimiter,
		"FormatString",
		"É stato richiesto il formato %s\n",
		greenery.DocCustomDelimiter,
	})
	require.NoError(t, err)

	customH := func(cfg greenery.Config, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("Only one argument is supported")
		}

		// Just for testing purposes, for a normal program this
		// would be a flag very likely
		vf := greenery.NewEnumValue("", "TEXT", "JSON", "XML")
		if err = vf.Set(args[0]); err != nil {
			return err
		}

		// Also for testing purposes, a normal program should error check on
		// the string.
		_, docs := cfg.GetDocs()
		fmt.Printf(docs.Custom["FormatString"], vf.Value)
		return nil
	}

	tcs := []testhelper.TestCase{
		testhelper.TestCase{
			Name: "version help",
			CmdLine: []string{
				"version",
				"-h",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": customH,
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.version.stdout")},
		},
		testhelper.TestCase{
			Name: "version actual exec",
			CmdLine: []string{
				"version",
				"JSON",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": customH,
			},
			NoValidateConfigValues: true,
			OutStdOut:              "Format JSON was requested\n",
		},
		testhelper.TestCase{
			Name: "root help",
			CmdLine: []string{
				"-h",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"root": customH,
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.root.stdout")},
		},
		testhelper.TestCase{
			Name: "root actual exec",
			CmdLine: []string{
				"JSON",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"root": customH,
			},
			NoValidateConfigValues: true,
			OutStdOut:              "Format JSON was requested\n",
		},
		testhelper.TestCase{
			Name: "config help",
			CmdLine: []string{
				"config",
				"-h",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"config": customH,
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.config.stdout")},
		},
		testhelper.TestCase{
			Name: "config actual exec",
			CmdLine: []string{
				"config",
				"JSON",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"config": customH,
			},
			NoValidateConfigValues: true,
			OutStdOut:              "Format JSON was requested\n",
		},
		testhelper.TestCase{
			Name: "help help",
			CmdLine: []string{
				"help",
				"-h",
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.help.stdout")},
		},

		// Same with a different language
		testhelper.TestCase{
			Name: "version help IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"version",
				"-h",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": customH,
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.version.it.stdout")},
		},
		testhelper.TestCase{
			Name: "version actual exec IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"version",
				"JSON",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"version": customH,
			},
			NoValidateConfigValues: true,
			OutStdOut:              "É stato richiesto il formato JSON\n",
		},
		testhelper.TestCase{
			Name: "root help IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"-h",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"root": customH,
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.root.it.stdout")},
		},
		testhelper.TestCase{
			Name: "root actual exec IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"JSON",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"root": customH,
			},
			NoValidateConfigValues: true,
			OutStdOut:              "É stato richiesto il formato JSON\n",
		},
		testhelper.TestCase{
			Name: "config help IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"-h",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"config": customH,
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.config.it.stdout")},
		},
		testhelper.TestCase{
			Name: "config actual exec IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"config",
				"JSON",
			},
			OverrideBuiltinHandlers: true,
			BuiltinHandlers: map[string]greenery.Handler{
				"config": customH,
			},
			NoValidateConfigValues: true,
			OutStdOut:              "É stato richiesto il formato JSON\n",
		},
		testhelper.TestCase{
			Name: "help help IT",
			Env: map[string]string{
				"LANG": "it_IT.UTF-8",
			},
			CmdLine: []string{
				"help",
				"-h",
			},
			NoValidateConfigValues: true,
			GoldStdOut:             &testhelper.TestFile{Source: filepath.Join("testdata", cmdTestName+".TestCommandsAdditional.help.it.stdout")},
		},
	}

	err = testhelper.RunTestCases(t, tcs, testhelper.TestRunnerOptions{
		ConfigGen: testhelper.NewSimpleConfig,
		UserDocList: map[string]*greenery.DocSet{
			"en": enSet,
			"it": itSet,
		},
	})
	require.NoError(t, err)
}
