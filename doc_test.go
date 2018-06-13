package greenery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var docsetStruct = &DocSet{
	Use:           "Use",
	Short:         "Short",
	Long:          "Long",
	Example:       "Example",
	HelpFlag:      "HelpFlag",
	CmdFlags:      "CmdFlags",
	ConfigEnvMsg1: "ConfigEnvMsg1",
	ConfigEnvMsg2: "ConfigEnvMsg2",
	ConfigEnvMsg3: "ConfigEnvMsg3",
	Help: &HelpStrings{
		Usage:                                "HUsage",
		Aliases:                              "HAliases",
		Examples:                             "HExamples",
		AvailableCommands:                    "HAvailableCommands",
		Flags:                                "HFlags",
		GlobalFlags:                          "HGlobalFlags",
		AdditionalHelpTopics:                 "HAdditionalHelpTopics",
		ProvidesMoreInformationAboutACommand: "HProvidesMoreInformationAboutACommand",
	},
	CmdLine: map[string]string{
		"cmd1": "valuecmd1",
		"cmd2": "valuecmd2",
	},
	ConfigFile: map[string]string{
		"cfg1": "valuecfg1",
		"cfg2": "valuecfg2",
	},
	Usage: map[string]*CmdHelp{
		"usage1": &CmdHelp{
			Use:     "Use1",
			Short:   "Short1",
			Long:    "Long1",
			Example: "Example1",
		},
		"usage2": &CmdHelp{
			Use:     "Use2",
			Short:   "Short2",
			Long:    "Long2",
			Example: "Example2",
		},
	},
	Custom: map[string]string{
		"k1": "v1",
		"k2": "v2",
	},
}

var baseSection = []string{
	DocBaseDelimiter,
	DocUse,
	"Use",
	DocShort,
	"Short",
	DocLong,
	"Long",
	DocExample,
	"Example",
	DocHelpFlag,
	"HelpFlag",
	DocCmdFlags,
	"CmdFlags",
	DocConfigEnvMsg1,
	"ConfigEnvMsg1",
	DocConfigEnvMsg2,
	"ConfigEnvMsg2",
	DocConfigEnvMsg3,
	"ConfigEnvMsg3",
	DocBaseDelimiter,
}

var helpSection = []string{
	DocHelpDelimiter,
	DocUsage,
	"HUsage",
	DocAliases,
	"HAliases",
	DocExamples,
	"HExamples",
	DocAvailableCommands,
	"HAvailableCommands",
	DocFlags,
	"HFlags",
	DocGlobalFlags,
	"HGlobalFlags",
	DocAdditionalHelpTopics,
	"HAdditionalHelpTopics",
	DocProvidesMoreInformationAboutACommand,
	"HProvidesMoreInformationAboutACommand",
	DocHelpDelimiter,
}

var usageSection = []string{
	DocUsageDelimiter,
	"usage1",
	"Use1",
	"Short1",
	"Long1",
	"Example1",
	"usage2",
	"Use2",
	"Short2",
	"Long2",
	"Example2",
	DocUsageDelimiter,
}

var cmdSection = []string{
	DocCmdlineDelimiter,
	"cmd1",
	"valuecmd1",
	"cmd2",
	"valuecmd2",
	DocCmdlineDelimiter,
}

var cfgSection = []string{
	DocConfigDelimiter,
	"cfg1",
	"valuecfg1",
	"cfg2",
	"valuecfg2",
	DocConfigDelimiter,
}

var customSection = []string{
	DocCustomDelimiter,
	"k1",
	"v1",
	"k2",
	"v2",
	DocCustomDelimiter,
}

func appender(elements ...[]string) []string {
	var res []string
	for _, v := range elements {
		res = append(res, v...)
	}
	return res
}

func TestConvertDocset(t *testing.T) {
	type docTest struct {
		Name    string
		Strings []string
		Set     *DocSet
		Error   string
	}

	expectedErr := "Documentation parse error: "

	sections := [][]string{
		baseSection,
		helpSection,
		usageSection,
		cmdSection,
		cfgSection,
		customSection,
	}

	t.Logf("Boo")
	tcs := []docTest{
		docTest{
			Name:    "Standard order",
			Strings: appender([]string{"1.0"}, appender(sections...)),
			Set:     docsetStruct,
		},
	}

	// Iterative heap's algorithm to get all possible permutations for our
	// sections.
	c := []int{0, 0, 0, 0, 0, 0}
	i := 0
	ti := 0
	for i < 6 {
		if c[i] < i {
			if i%2 == 0 {
				sections[0], sections[i] = sections[i], sections[0]
			} else {
				sections[c[i]], sections[i] = sections[i], sections[c[i]]
			}

			ti++
			c[i]++
			i = 0

			tcs = append(tcs, docTest{
				Name:    fmt.Sprintf("Swap %d", ti),
				Strings: appender([]string{"1.0"}, appender(sections...)),
				Set:     docsetStruct,
			})
		} else {
			c[i] = 0
			i++
		}
	}

	tcs = append(tcs, []docTest{
		docTest{
			Name: "Duplicate 1",
			Strings: appender(
				[]string{"1.0"},
				baseSection,
				baseSection,
				helpSection,
				usageSection,
				cmdSection,
				cfgSection,
				customSection,
			),
			Error: expectedErr + "duplicate base section",
		},
		docTest{
			Name: "Duplicate 2",
			Strings: appender(
				[]string{"1.0"},
				baseSection,
				helpSection,
				helpSection,
				usageSection,
				cmdSection,
				cfgSection,
				customSection,
			),
			Error: expectedErr + "duplicate help section",
		},
		docTest{
			Name: "Duplicate 3",
			Strings: appender(
				[]string{"1.0"},
				baseSection,
				helpSection,
				usageSection,
				usageSection,
				cmdSection,
				cfgSection,
				customSection,
			),
			Error: expectedErr + "duplicate usage section",
		},
		docTest{
			Name: "Duplicate 4",
			Strings: appender(
				[]string{"1.0"},
				baseSection,
				helpSection,
				usageSection,
				cmdSection,
				cmdSection,
				cfgSection,
				customSection,
			),
			Error: expectedErr + "duplicate commandline section",
		},
		docTest{
			Name: "Duplicate 5",
			Strings: appender(
				[]string{"1.0"},
				baseSection,
				helpSection,
				usageSection,
				cmdSection,
				cfgSection,
				cfgSection,
				customSection,
			),
			Error: expectedErr + "duplicate config section",
		},
		docTest{
			Name: "Duplicate 6",
			Strings: appender(
				[]string{"1.0"},
				baseSection,
				helpSection,
				usageSection,
				cmdSection,
				cfgSection,
				customSection,
				customSection,
			),
			Error: expectedErr + "duplicate custom section",
		},

		// Wrong lines tests, all the same for different sections
		docTest{
			// Help section last, odd number of lines
			Name: "Bad help 1",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:4],
			),
			Error: "odd number of help section lines",
		},
		docTest{
			// Cut the helpsection short, so it gets a wrongly named variable,
			Name: "Bad help 2",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:5],
				[]string{"bad", "value"},
			),
			Error: "unknown help variable bad",
		},
		docTest{
			Name: "Bad base 1",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:4],
			),
			Error: "odd number of base section lines",
		},
		docTest{
			Name: "Bad base 2",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:5],
				[]string{"bad", "value"},
			),
			Error: "unknown base variable bad",
		},
		docTest{
			Name: "Bad usage",
			Strings: appender(
				[]string{"1.0"},
				usageSection[:5],
			),
			Error: "not divisible by 4",
		},
		docTest{
			Name: "Bad cmd and cfg",
			Strings: appender(
				[]string{"1.0"},
				cfgSection[:4],
			),
			Error: "odd number of config section lines",
		},
		docTest{
			Name: "Unknown section type",
			Strings: appender(
				[]string{"1.0"},
				cfgSection,
				[]string{"UNKNOWN"},
			),
			Error: "unknown section UNKNOWN",
		},
		docTest{
			Name: "Cfg empty variable",
			Strings: appender(
				[]string{"1.0"},
				cfgSection[:3],
				[]string{"", "empty"},
			),
			Error: "empty variable name on line",
		},
		docTest{
			Name: "Cmd empty variable",
			Strings: appender(
				[]string{"1.0"},
				cmdSection[:3],
				[]string{"", "empty"},
			),
			Error: "empty variable name on line",
		},
		docTest{
			Name: "Custom empty variable",
			Strings: appender(
				[]string{"1.0"},
				customSection[:3],
				[]string{"", "empty"},
			),
			Error: "empty variable name on line",
		},
		docTest{
			Name: "Bad version",
			Strings: appender(
				[]string{"2.0"},
				customSection,
			),
			Error: "Unsupported document format: 2.0",
		},
		docTest{
			Name: "Duplicate base 1",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocUse, "Use", DocBaseDelimiter},
			),
			Error: "duplicate Use declaration at line 20 (previously line 2)",
		},
		docTest{
			Name: "Duplicate base 2",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocShort, "Short", DocBaseDelimiter},
			),
			Error: "duplicate Short declaration at line 20 (previously line 4)",
		},
		docTest{
			Name: "Duplicate base 3",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocLong, "Long", DocBaseDelimiter},
			),
			Error: "duplicate Long declaration at line 20 (previously line 6)",
		},
		docTest{
			Name: "Duplicate base 4",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocExample, "Example", DocBaseDelimiter},
			),
			Error: "duplicate Example declaration at line 20 (previously line 8)",
		},
		docTest{
			Name: "Duplicate base 5",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocHelpFlag, "HelpFlag", DocBaseDelimiter},
			),
			Error: "duplicate HelpFlag declaration at line 20 (previously line 10)",
		},
		docTest{
			Name: "Duplicate base 6",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocConfigEnvMsg1, "ConfigEnvMsg1", DocBaseDelimiter},
			),
			Error: "duplicate ConfigEnvMsg1 declaration at line 20 (previously line 14)",
		},
		docTest{
			Name: "Duplicate base 7",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocConfigEnvMsg2, "ConfigEnvMsg2", DocBaseDelimiter},
			),
			Error: "duplicate ConfigEnvMsg2 declaration at line 20 (previously line 16)",
		},
		docTest{
			Name: "Duplicate base 8",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocConfigEnvMsg3, "ConfigEnvMsg3", DocBaseDelimiter},
			),
			Error: "duplicate ConfigEnvMsg3 declaration at line 20 (previously line 18)",
		},
		docTest{
			Name: "Duplicate base 3",
			Strings: appender(
				[]string{"1.0"},
				baseSection[:len(baseSection)-1],
				[]string{DocCmdFlags, "CmdFlags", DocBaseDelimiter},
			),
			Error: "duplicate CmdFlags declaration at line 20 (previously line 12)",
		},
		docTest{
			Name: "Duplicate help 1",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocUsage, "DocUsage", DocHelpDelimiter},
			),
			Error: "duplicate Usage declaration at line 18 (previously line 2)",
		},
		docTest{
			Name: "Duplicate help 2",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocAliases, "DocAliases", DocHelpDelimiter},
			),
			Error: "duplicate Aliases declaration at line 18 (previously line 4)",
		},
		docTest{
			Name: "Duplicate help 3",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocExamples, "DocExamples", DocHelpDelimiter},
			),
			Error: "duplicate Examples declaration at line 18 (previously line 6)",
		},
		docTest{
			Name: "Duplicate help 4",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocAvailableCommands, "DocAvailableCommands", DocHelpDelimiter},
			),
			Error: "duplicate AvailableCommands declaration at line 18 (previously line 8)",
		},
		docTest{
			Name: "Duplicate help 5",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocFlags, "DocFlags", DocHelpDelimiter},
			),
			Error: "duplicate Flags declaration at line 18 (previously line 10)",
		},
		docTest{
			Name: "Duplicate help 6",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocGlobalFlags, "DocGlobalFlags", DocHelpDelimiter},
			),
			Error: "duplicate GlobalFlags declaration at line 18 (previously line 12)",
		},
		docTest{
			Name: "Duplicate help 7",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocAdditionalHelpTopics, "DocAdditionalHelpTopics", DocHelpDelimiter},
			),
			Error: "duplicate AdditionalHelpTopics declaration at line 18 (previously line 14)",
		},
		docTest{
			Name: "Duplicate help 8",
			Strings: appender(
				[]string{"1.0"},
				helpSection[:len(helpSection)-1],
				[]string{DocProvidesMoreInformationAboutACommand, "DocProvidesMoreInformationAboutACommand", DocHelpDelimiter},
			),
			Error: "duplicate ProvidesMoreInformationAboutACommand declaration at line 18 (previously line 16)",
		},
		docTest{
			Name: "Duplicate usage",
			Strings: appender(
				[]string{"1.0"},
				usageSection[:len(usageSection)-1],
				[]string{"usage1",
					"Use1",
					"Short1",
					"Long1",
					"Example1",
					DocUsageDelimiter},
			),
			Error: "duplicate usage1 declaration at line 12 (previously line 2)",
		},
		docTest{
			Name: "Duplicate cmd",
			Strings: appender(
				[]string{"1.0"},
				cmdSection[:len(cmdSection)-1],
				[]string{"cmd1", "valuecmd1", DocCmdlineDelimiter},
			),
			Error: "duplicate cmd1 declaration at line 6 (previously line 2)",
		},
		docTest{
			Name: "Duplicate cfg",
			Strings: appender(
				[]string{"1.0"},
				cfgSection[:len(cfgSection)-1],
				[]string{"cfg1", "valuecfg1", DocConfigDelimiter},
			),
			Error: "duplicate cfg1 declaration at line 6 (previously line 2)",
		},
		docTest{
			Name: "Duplicate custom",
			Strings: appender(
				[]string{"1.0"},
				customSection[:len(customSection)-1],
				[]string{"k2", "v2", DocCustomDelimiter},
			),
			Error: "duplicate k2 declaration at line 6 (previously line 4)",
		},
	}...)

	for _, tt := range tcs {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			set, err := ConvertDocs(tt.Strings)
			if tt.Error != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.Error)
			} else {
				require.NoError(t, err)
			}

			if tt.Set != nil {
				require.Equal(t, set, tt.Set)
			}
		})
	}

	badSet := map[string][]string{
		"ok":     appender([]string{"1.0"}, baseSection),
		"not-ok": appender([]string{"1.0"}, baseSection[:4]),
	}
	set, err := convertDocset(badSet)
	require.Nil(t, set)
	require.Contains(t, err.Error(), "While processing language not-ok")

	require.PanicsWithValue(t, "Internal error generating the documentation: While processing language not-ok: Documentation parse error: odd number of base section lines", func() {
		ConvertOrPanic(badSet)
	})
}

func TestGetDocset(t *testing.T) {
	en := &DocSet{Use: "en"}
	it := &DocSet{Use: "it"}
	itFull := &DocSet{Use: "it-full"}
	docs := map[string]*DocSet{
		"en":          en,
		"it":          it,
		"it_IT.UTF-8": itFull,
		"verify":      docsetStruct,
	}

	type docTest struct {
		Name    string
		Lang    string
		Default string
		Set     *DocSet
		Found   string
	}

	tcs := []docTest{
		docTest{
			Name:    "Default, full LANG",
			Lang:    "en_US.UTF-8",
			Default: "en",
			Set:     en,
			Found:   "en",
		},
		docTest{
			Name:    "Default, partial LANG",
			Lang:    "en",
			Default: "en",
			Set:     en,
			Found:   "en",
		},
		docTest{
			Name:    "Default, empty LANG",
			Lang:    "en",
			Default: "en",
			Set:     en,
			Found:   "en",
		},
		docTest{
			Name:    "Not default, full LANG",
			Lang:    "it_IT.UTF-8",
			Default: "en",
			Set:     itFull,
			Found:   "it_IT.UTF-8",
		},
		docTest{
			Name:    "Not default, partial LANG",
			Lang:    "it_",
			Default: "en",
			Set:     it,
			Found:   "it",
		},
		docTest{
			Name:    "Verify deepcopy",
			Lang:    "verify",
			Default: "en",
			Set:     docsetStruct,
			Found:   "verify",
		},
	}

	for _, tt := range tcs {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			set, found, err := getDocset(docs, tt.Lang, tt.Default)
			require.NoError(t, err)
			require.Equal(t, tt.Set.Use, set.Use)
			require.Equal(t, found, tt.Found)

			if tt.Found == "verify" {
				// Kinda verify this was a deepcopy, values should be the
				// same, but not the actual returned struct, unfortunately
				// can't compare the maps directly
				require.Equal(t, set, tt.Set)
				require.True(t, set != tt.Set)
				require.True(t, set.Usage["usage1"] != tt.Set.Usage["usage1"])
				require.True(t, set.Usage["usage2"] != tt.Set.Usage["usage2"])
			}
		})
	}
}
