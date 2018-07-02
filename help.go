package greenery

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/woodensquares/greenery/internal/doc"
)

// These constants are to be used when creating string-list user
// documentation, they contain identifiers used for the section delimiters, as
// well as for all the individual options.
const (
	// DocBaseDelimiter is the delimiter to be used for the base section, this
	// will map to the toplevel of the DocSet struct.
	DocBaseDelimiter = doc.BaseDelimiter

	// DocUse contains the text to display in the use line for the application
	// itself, this is typically used if the root command has a handler which
	// takes arguments. The text to be used should follow in the next string in
	// the list. The string following will be assigned to DocSet.Use.
	DocUse = doc.Use

	// DocLong is the long description for the application as a whole. The text to
	// be used should follow in the next string in the list. The string following
	// will be assigned to DocSet.Long.
	DocLong = doc.Long

	// DocShort is the short description for the application as a whole. The text
	// to be used should follow in the next string in the list.. The string
	// following will be assigned to DocSet.Short.
	DocShort = doc.Short

	// DocExample is an example to be used for the application as a whole. The
	// text to be used should follow in the next string in the list. . The string
	// following will be assigned to DocSet.Example.
	DocExample = doc.Example

	// DocHelpFlag is the text displayed next to the help command when printing
	// the help for the application. The text to be used should follow in the next
	// string. The string following will be assigned to DocSet.HelpFlag.
	DocHelpFlag = doc.HelpFlag

	// DocCmdFlags is the text to use for "[flags]" as printed in command help
	// messages. The text to be used should follow in the next string in the
	// list. The string following will be assigned to DocSet.CmdFlags.
	DocCmdFlags = doc.CmdFlags

	// DocConfigEnvMsg1 is part of the localized text to be used for the config
	// env command. The text to be used should follow in the next string in the
	// list. The string following will be assigned to DocSet.ConfigEnvMsg1.
	DocConfigEnvMsg1 = doc.ConfigEnvMsg1

	// DocConfigEnvMsg2 is part of the localized text to be used for the config
	// env command. The text to be used should follow in the next string in the
	// list. The string following will be assigned to DocSet.ConfigEnvMsg2.
	DocConfigEnvMsg2 = doc.ConfigEnvMsg2

	// DocConfigEnvMsg3 is part of the localized text to be used for the config
	// env command. The text to be used should follow in the next string in the
	// list. The string following will be assigned to DocSet.ConfigEnvMsg3.
	DocConfigEnvMsg3 = doc.ConfigEnvMsg3

	// DocHelpDelimiter is the delimiter of the help section. This will map to the
	// Help field in the DocSet struct, which is a HelpStrings struct. These
	// strings will be used to localize the overall common headers/footers/... in
	// various help messages.
	DocHelpDelimiter = doc.HelpDelimiter

	// DocUsage controls the "Usage:" sentence. The text to be used should follow in
	// the next string in the list. The string following will be assigned to
	// HelpStrings.Usage.
	DocUsage = doc.Usage

	// DocAliases controls the "Aliases:" sentence. The text to be used should
	// follow in the next string in the list. The string following will be
	// assigned to HelpStrings.Aliases.
	DocAliases = doc.Aliases

	// DocExamples controls the "Examples:" sentence. The text to be used should
	// follow in the next string in the list. The string following will be assigned to
	// HelpStrings.Examples.
	DocExamples = doc.Examples

	// DocAvailableCommands controls the "Available Commands:" sentence. The text
	// to be used should follow in the next string in the list. The string
	// following will be assigned to HelpStrings.AvailableCommands.
	DocAvailableCommands = doc.AvailableCommands

	// DocFlags controls the "Flags:" sentence. The text to be used should follow
	// in the next string in the list. The string following will be assigned to
	// HelpStrings.Flags.
	DocFlags = doc.Flags

	// DocGlobalFlags controls the "Global Flags:" sentence. The text to be used
	// should follow in the next string in the list. The string following will be
	// assigned to HelpStrings.GlobalFlags.
	DocGlobalFlags = doc.GlobalFlags

	// DocAdditionalHelpTopics controls the "Additional help topics:"
	// sentence. The text to be used should follow in the next string in the
	// list. The string following will be assigned to
	// HelpStrings.AdditionalHelpTopics.
	DocAdditionalHelpTopics = doc.AdditionalHelpTopics

	// DocProvidesMoreInformationAboutACommand controls the "provides more
	// information about a command." sentence. The text to be used should follow
	// in the next string in the list. The string following will be assigned to
	// HelpStrings.ProvidesMoreInformationAboutACommand.
	DocProvidesMoreInformationAboutACommand = doc.ProvidesMoreInformationAboutACommand

	// DocUsageDelimiter is the delimiter for the usage section, this will map to
	// the Usage field in the DocSet struct, which is a CmdHelp struct. This is
	// used to contain help for individual commands, existing and
	// user-created. Each command is controlled by a sequence of five strings,
	// they are in order:
	//
	// identifier of the command (which is the key in the handlers map, this can
	//                            also be one of the predefined identifiers below)
	// usage string (in case the command has arguments), will map to CmdHelp.Use
	// short help for the command, will map to CmdHelp.Short
	// long help for the command, will map to CmdHelp.Long
	// example for the command, will map to CmdHelp.Example
	//
	// Internal commands have predefined constants, so they can be used instead of
	// writing things like "config>init" which might be confusing for the
	// localizers.
	DocUsageDelimiter = doc.UsageDelimiter

	// DocConfigCmd is the identifier for the "config" command.
	DocConfigCmd = doc.ConfigCmd

	// DocConfigInitCmd is the identifier for the "config init" command.
	DocConfigInitCmd = doc.ConfigInitCmd

	// DocConfigEnvCmd is the identifier for the "config env" command.
	DocConfigEnvCmd = doc.ConfigEnvCmd

	// DocConfigDisplayCmd is the identifier for the "config display" command.
	DocConfigDisplayCmd = doc.ConfigDisplayCmd

	// DocHelpCmd is the identifier for the "help" command.
	DocHelpCmd = doc.HelpCmd

	// DocVersionCmd is the identifier for the "version" command.
	DocVersionCmd = doc.VersionCmd

	// DocCmdlineDelimiter is the delimiter for the command line flags section,
	// this will map to the CmdLine field in the DocSet struct, which is a map of
	// strings. In this section the format is simply
	//
	// name of the variable corresponding to the flag
	// help information for this flag to be displayed
	//
	// as in the usage section above, predefined identifiers are included for the
	// base flags defined by the library.
	DocCmdlineDelimiter = doc.CmdlineDelimiter

	// DocLogLevel is the help information for the LogLevel flag.
	DocLogLevel = doc.LogLevel

	// DocConfFile is the help information for the ConfFile flag.
	DocConfFile = doc.ConfFile

	// DocLogFile is the help information for the LogFile flag.
	DocLogFile = doc.LogFile

	// DocPretty is the help information for the Pretty flag.
	DocPretty = doc.Pretty

	// DocNoEnv is the help information for the NoEnv flag.
	DocNoEnv = doc.NoEnv

	// DocNoCfg is the help information for the NoCfg flag.
	DocNoCfg = doc.NoCfg

	// DocVerbosity is the help information for the Verbosity flag.
	DocVerbosity = doc.Verbosity

	// DocDoTrace is the help information for the DoTrace flag.
	DocDoTrace = doc.DoTrace

	// DocCfgLocation is the help information for the CfgLocation flag.
	DocCfgLocation = doc.CfgLocation

	// DocCfgForce is the help information for the CfgForce flag.
	DocCfgForce = doc.CfgForce

	// DocConfigDelimiter is the delimiter for the config file section, this will
	// map to the ConfigFile field in the DocSet struct, which is a map of
	// strings. This section has the exact same structure as the Cmdline section
	// and should contain only entries for which the text between the command line
	// and config file help is different. In addition it also typically contains
	// the help displayed for custom configuration sections, the predefined base
	// section name follows.
	DocConfigDelimiter = doc.ConfigDelimiter

	// DocConfigHeader is the name of the base configuration file section, this is
	// going to be treated as a golang template, {{ date }} is available and will
	// be substituted with the date/time the configuration file has been generated
	// at.
	DocConfigHeader = doc.ConfigHeader

	// DocCustomDelimiter is the delimiter for the custom file section, this
	// will map to the Custom field in the DocSet struct, which is a map of
	// strings. This section is not used by the standard library but is
	// available for custom string localizations
	DocCustomDelimiter = doc.CustomDelimiter

	// errText is the string corresponding to the documentation parse error.
	errText = "Documentation parse error:"
)

// getDocset will return an applicable DocSet based on the lang string passed
// (typically the LANG environment variable) as well as what language this
// was. If no docs were found, the default docset is returned. Note this will
// also try to match on the first two letters, as typically, internally at
// least, documentation will be for lang as opposed to lang_country.
func getDocset(documentation map[string]*DocSet, lang, defLang string) (*DocSet, string, error) {
	var docs *DocSet
	var foundLanguage string

	if lang != "" {
		if exact, ok := documentation[lang]; ok {
			docs = exact
			foundLanguage = lang
		} else if len(lang) > 1 {
			cLanguage := lang[:2]
			if sub, ok := documentation[cLanguage]; ok {
				foundLanguage = cLanguage
				docs = sub
			}
		}
	}

	if foundLanguage == "" {
		var ok bool
		if docs, ok = documentation[defLang]; !ok {
			return &DocSet{}, "", fmt.Errorf("No docset present for the default language: %s", defLang)
		}
		foundLanguage = defLang
	}

	// Need to create a new deep copy of this docset before returning to avoid
	// issues especially during tests.
	hs := HelpStrings{}
	if docs.Help != nil {
		hs = *docs.Help
	}

	cmds := map[string]string{}
	cf := map[string]string{}
	custom := map[string]string{}
	usage := map[string]*CmdHelp{}
	for k, v := range docs.CmdLine {
		cmds[k] = v
	}
	for k, v := range docs.ConfigFile {
		cf[k] = v
	}
	for k, v := range docs.Custom {
		custom[k] = v
	}
	for k, v := range docs.Usage {
		h := CmdHelp{}
		if v == nil {
			return &DocSet{}, "", fmt.Errorf("Nil CmdHelp entry for command %s in language %s", k, foundLanguage)
		}
		h = *v
		usage[k] = &h
	}

	return &DocSet{
		Use:           docs.Use,
		Short:         docs.Short,
		Long:          docs.Long,
		HelpFlag:      docs.HelpFlag,
		CmdFlags:      docs.CmdFlags,
		Example:       docs.Example,
		ConfigEnvMsg1: docs.ConfigEnvMsg1,
		ConfigEnvMsg2: docs.ConfigEnvMsg2,
		ConfigEnvMsg3: docs.ConfigEnvMsg3,
		Help:          &hs,
		CmdLine:       cmds,
		ConfigFile:    cf,
		Usage:         usage,
		Custom:        custom,
	}, foundLanguage, nil
}

// convertDocset is used to go from a map of []string formatted documentations
// to a map of *DocSet ones by calling ConvertDocs as needed.
func convertDocset(unprocessedDocs map[string][]string) (map[string]*DocSet, error) {
	convertedDocs := make(map[string]*DocSet)
	for lang, udoc := range unprocessedDocs {
		converted, err := ConvertDocs(udoc)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("While processing language %s", lang))
		}

		convertedDocs[lang] = converted
	}
	return convertedDocs, nil
}

// ConvertDocs is used to convert a []string documentation to the internal
// *DocSet representation, it will return an error if there are any parsing
// issues. Note the code can't be super precise in the validation given a
// single list of strings with arbitrary names, it tries to catch at least
// obvious things like multiple sections or wrong number of lines in a section
// etc.
func ConvertDocs(udoc []string) (*DocSet, error) {
	langDoc := &DocSet{}

	var i int
	l := len(udoc)

	// Let's be a bit resilient and allow users to sort the sections
	// differently from the default.
	var foundBase, foundHelp, foundUsage, foundCmd, foundCfg, foundCustom bool

	for {
		if i >= l {
			// We are done
			break
		}

		if i == 0 {
			if udoc[i] != "1.0" {
				return nil, fmt.Errorf("Unsupported document format: %s", udoc[i])
			}
			i++
			continue
		}

		switch udoc[i] {
		case doc.BaseDelimiter:
			if foundBase {
				return nil, fmt.Errorf("%s duplicate base section at line %d", errText, i)
			}
			foundBase = true
			idx, err := helperBaseCmd(langDoc, udoc, i+1, l)
			if err != nil {
				return nil, err
			}
			i = idx
		case doc.HelpDelimiter:
			if foundHelp {
				return nil, fmt.Errorf("%s duplicate help section at line %d", errText, i)
			}
			foundHelp = true
			r, idx, err := helperHelp(udoc, i+1, l)
			if err != nil {
				return nil, err
			}
			i = idx
			langDoc.Help = r
		case doc.UsageDelimiter:
			if foundUsage {
				return nil, fmt.Errorf("%s duplicate usage section at line %d", errText, i)
			}
			foundUsage = true
			r, idx, err := usageHelp(udoc, i+1, l)
			if err != nil {
				return nil, err
			}
			i = idx
			langDoc.Usage = r
		case doc.CmdlineDelimiter:
			if foundCmd {
				return nil, fmt.Errorf("%s duplicate commandline section at line %d", errText, i)
			}
			foundCmd = true
			r, idx, err := mapMap(udoc, i+1, l, "cmd")
			if err != nil {
				return nil, err
			}
			i = idx
			langDoc.CmdLine = r
		case doc.ConfigDelimiter:
			if foundCfg {
				return nil, fmt.Errorf("%s duplicate config section at line %d", errText, i)
			}
			foundCfg = true
			r, idx, err := mapMap(udoc, i+1, l, "config")
			if err != nil {
				return nil, err
			}
			i = idx
			langDoc.ConfigFile = r
		case doc.CustomDelimiter:
			if foundCustom {
				return nil, fmt.Errorf("%s duplicate custom section at line %d", errText, i)
			}
			foundCustom = true
			r, idx, err := mapMap(udoc, i+1, l, "custom")
			if err != nil {
				return nil, err
			}
			i = idx
			langDoc.Custom = r
		default:
			return nil, fmt.Errorf("%s unknown section %s at line %d", errText, udoc[i], i)
		}

	}

	return langDoc, nil
}

// usageStruct is used for the first pass of the template preprocessing
type usageStruct struct {
	*HelpStrings
	OpenBrace  string
	CloseBrace string
	Command    string
}

// fixMagicTemplate will deal with the issue of "[flags]" being hardcoded
// inside Cobra, via a hackish search/replace, hopefully nobody creates an
// application named ☺示☺示☺示☺示☺示☺...
func fixMagicTemplate(s, uflag, dflag string) []string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if strings.Contains(l, "T☺示☺示☺示☺示☺示☺T") {
			l = strings.Replace(l, "T☺示☺示☺示☺示☺示☺T", "", 1)
			if uflag != "" {
				l = strings.Replace(l, "[flags]", uflag, 1)
			} else {
				l = strings.Replace(l, "[flags]", dflag, 1)
			}
		}
		lines[i] = l
	}
	return lines
}

// The usage template, given that we are processing it twice all the {{ and }}
// that will be used are simply switched to .Open/CloseBrace and will be
// substituted the first time.
var usageTemplate = `{{ .Usage }}{{ .OpenBrace }}if .Runnable{{ .CloseBrace }}
  T☺示☺示☺示☺示☺示☺T{{ .OpenBrace }}.UseLine{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasAvailableSubCommands{{ .CloseBrace }}
  {{ .OpenBrace }}.CommandPath{{ .CloseBrace }} {{ .Command }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if gt (len .Aliases) 0{{ .CloseBrace }}

{{ .Aliases }}
  {{ .OpenBrace }}.NameAndAliases{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasExample{{ .CloseBrace }}

{{ .Examples }}
{{ .OpenBrace }}.Example{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasAvailableSubCommands{{ .CloseBrace }}

{{ .AvailableCommands}}{{ .OpenBrace }}range .Commands{{ .CloseBrace }}{{ .OpenBrace }}if (or .IsAvailableCommand (eq .Name "help")){{ .CloseBrace }}
  {{ .OpenBrace }}rpad .Name .NamePadding {{ .CloseBrace }} {{ .OpenBrace }}.Short{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasAvailableLocalFlags{{ .CloseBrace }}

{{ .Flags }}
{{ .OpenBrace }}.LocalFlags.FlagUsages | trimTrailingWhitespaces{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasAvailableInheritedFlags{{ .CloseBrace }}

{{ .GlobalFlags }}
{{ .OpenBrace }}.InheritedFlags.FlagUsages | trimTrailingWhitespaces{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasHelpSubCommands{{ .CloseBrace }}

{{ .AdditionalHelpTopics }}{{ .OpenBrace }}range .Commands{{ .CloseBrace }}{{ .OpenBrace }}if .IsAdditionalHelpTopicCommand{{ .CloseBrace }}
  {{ .OpenBrace }}rpad .CommandPath .CommandPathPadding{{ .CloseBrace }} {{ .OpenBrace }}.Short{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}end{{ .CloseBrace }}{{ .OpenBrace }}if .HasAvailableSubCommands{{ .CloseBrace }}

"{{ .OpenBrace }}.CommandPath{{ .CloseBrace }} {{ .Command }} --help" {{.ProvidesMoreInformationAboutACommand}}{{ .OpenBrace }}end{{ .CloseBrace }}
`

var defaultDocList = ConvertOrPanic(doc.DefaultDocs)

// ConvertOrPanic is used to convert a set of documents from the []string
// format to the final map[string]*DocSet one, this function will panic in
// case there are any errors in the conversion.
func ConvertOrPanic(unprocessedDocs map[string][]string) map[string]*DocSet {
	docs, err := convertDocset(unprocessedDocs)
	if err != nil {
		panic("Internal error generating the documentation: " + err.Error())
	}
	return docs
}

// helperBaseCmd is the parser for the base section of the []string doc
// templates
func helperBaseCmd(langDoc *DocSet, docs []string, idx, l int) (int, error) {
	var i int

	var foundUse, foundShort, foundLong, foundExample, foundHelpFlag int
	var foundConfigEnvMsg1, foundConfigEnvMsg2, foundConfigEnvMsg3, foundCmdFlags int

	for i = idx; i < l; i++ {
		if docs[i] == doc.BaseDelimiter {
			i++
			break
		}

		if i == l-1 {
			return 0, fmt.Errorf("%s odd number of base section lines", errText)
		}

		switch docs[i] {
		case doc.Use:
			if foundUse > 0 {
				return 0, fmt.Errorf("%s duplicate Use declaration at line %d (previously line %d)", errText, i, foundUse)
			}
			foundUse = i
			langDoc.Use = docs[i+1]
			i++
		case doc.Short:
			if foundShort > 0 {
				return 0, fmt.Errorf("%s duplicate Short declaration at line %d (previously line %d)", errText, i, foundShort)
			}
			foundShort = i
			langDoc.Short = docs[i+1]
			i++
		case doc.Long:
			if foundLong > 0 {
				return 0, fmt.Errorf("%s duplicate Long declaration at line %d (previously line %d)", errText, i, foundLong)
			}
			foundLong = i
			langDoc.Long = docs[i+1]
			i++
		case doc.Example:
			if foundExample > 0 {
				return 0, fmt.Errorf("%s duplicate Example declaration at line %d (previously line %d)", errText, i, foundExample)
			}
			foundExample = i
			langDoc.Example = docs[i+1]
			i++
		case doc.HelpFlag:
			if foundHelpFlag > 0 {
				return 0, fmt.Errorf("%s duplicate HelpFlag declaration at line %d (previously line %d)", errText, i, foundHelpFlag)
			}
			foundHelpFlag = i
			langDoc.HelpFlag = docs[i+1]
			i++
		case doc.ConfigEnvMsg1:
			if foundConfigEnvMsg1 > 0 {
				return 0, fmt.Errorf("%s duplicate ConfigEnvMsg1 declaration at line %d (previously line %d)", errText, i, foundConfigEnvMsg1)
			}
			foundConfigEnvMsg1 = i
			langDoc.ConfigEnvMsg1 = docs[i+1]
			i++
		case doc.ConfigEnvMsg2:
			if foundConfigEnvMsg2 > 0 {
				return 0, fmt.Errorf("%s duplicate ConfigEnvMsg2 declaration at line %d (previously line %d)", errText, i, foundConfigEnvMsg2)
			}
			foundConfigEnvMsg2 = i
			langDoc.ConfigEnvMsg2 = docs[i+1]
			i++
		case doc.ConfigEnvMsg3:
			if foundConfigEnvMsg3 > 0 {
				return 0, fmt.Errorf("%s duplicate ConfigEnvMsg3 declaration at line %d (previously line %d)", errText, i, foundConfigEnvMsg3)
			}
			foundConfigEnvMsg3 = i
			langDoc.ConfigEnvMsg3 = docs[i+1]
			i++
		case doc.CmdFlags:
			if foundCmdFlags > 0 {
				return 0, fmt.Errorf("%s duplicate CmdFlags declaration at line %d (previously line %d)", errText, i, foundCmdFlags)
			}
			foundCmdFlags = i
			langDoc.CmdFlags = docs[i+1]
			i++
		default:
			return 0, fmt.Errorf("%s internal parse error, unknown base variable %s at line %d %v", errText, docs[i], i, docs)
		}
	}

	return i, nil
}

// helperHelp is the parser for the help section of the []string doc templates
func helperHelp(docs []string, idx, l int) (*HelpStrings, int, error) {
	hStrings := &HelpStrings{}
	var i int

	var foundUsage, foundAliases, foundExamples, foundAvailableCommands int
	var foundFlags, foundGlobalFlags, foundAdditionalHelpTopics int
	var foundProvidesMoreInformationAboutACommand int

	for i = idx; i < l; i++ {
		if docs[i] == doc.HelpDelimiter {
			i++
			break
		}

		if i == l-1 {
			return nil, 0, fmt.Errorf("%s odd number of help section lines", errText)
		}

		switch docs[i] {
		case doc.Usage:
			if foundUsage > 0 {
				return nil, 0, fmt.Errorf("%s duplicate Usage declaration at line %d (previously line %d)", errText, i, foundUsage)
			}
			foundUsage = i
			hStrings.Usage = docs[i+1]
			i++
		case doc.Aliases:
			if foundAliases > 0 {
				return nil, 0, fmt.Errorf("%s duplicate Aliases declaration at line %d (previously line %d)", errText, i, foundAliases)
			}
			foundAliases = i
			hStrings.Aliases = docs[i+1]
			i++
		case doc.Examples:
			if foundExamples > 0 {
				return nil, 0, fmt.Errorf("%s duplicate Examples declaration at line %d (previously line %d)", errText, i, foundExamples)
			}
			foundExamples = i
			hStrings.Examples = docs[i+1]
			i++
		case doc.AvailableCommands:
			if foundAvailableCommands > 0 {
				return nil, 0, fmt.Errorf("%s duplicate AvailableCommands declaration at line %d (previously line %d)", errText, i, foundAvailableCommands)
			}
			foundAvailableCommands = i
			hStrings.AvailableCommands = docs[i+1]
			i++
		case doc.Flags:
			if foundFlags > 0 {
				return nil, 0, fmt.Errorf("%s duplicate Flags declaration at line %d (previously line %d)", errText, i, foundFlags)
			}
			foundFlags = i
			hStrings.Flags = docs[i+1]
			i++
		case doc.GlobalFlags:
			if foundGlobalFlags > 0 {
				return nil, 0, fmt.Errorf("%s duplicate GlobalFlags declaration at line %d (previously line %d)", errText, i, foundGlobalFlags)
			}
			foundGlobalFlags = i
			hStrings.GlobalFlags = docs[i+1]
			i++
		case doc.AdditionalHelpTopics:
			if foundAdditionalHelpTopics > 0 {
				return nil, 0, fmt.Errorf("%s duplicate AdditionalHelpTopics declaration at line %d (previously line %d)", errText, i, foundAdditionalHelpTopics)
			}
			foundAdditionalHelpTopics = i
			hStrings.AdditionalHelpTopics = docs[i+1]
			i++
		case doc.ProvidesMoreInformationAboutACommand:
			if foundProvidesMoreInformationAboutACommand > 0 {
				return nil, 0, fmt.Errorf("%s duplicate ProvidesMoreInformationAboutACommand declaration at line %d (previously line %d)", errText, i, foundProvidesMoreInformationAboutACommand)
			}
			foundProvidesMoreInformationAboutACommand = i
			hStrings.ProvidesMoreInformationAboutACommand = docs[i+1]
			i++
		default:
			return nil, 0, fmt.Errorf("%s internal parse error, unknown help variable %s at line %d", errText, docs[i], i)
		}
	}

	return hStrings, i, nil
}

// usageHelp is the parser for the usage section of the []string doc templates
func usageHelp(docs []string, idx, l int) (map[string]*CmdHelp, int, error) {
	cmds := map[string]*CmdHelp{}
	var i int
	seen := map[string]int{}

	for i = idx; i < l; i++ {
		if docs[i] == doc.UsageDelimiter {
			i++
			break
		}

		if _, s := seen[docs[i]]; s {
			return nil, 0, fmt.Errorf("%s duplicate %s declaration at line %d (previously line %d)", errText, docs[i], i, seen[docs[i]])
		}
		seen[docs[i]] = i

		if i >= l-4 {
			return nil, 0, fmt.Errorf("%s number of section lines not divisible by 4 in the Usage section %d / %d, %v", errText, i, l, docs)
		}

		cmd := CmdHelp{
			Use:     docs[i+1],
			Short:   docs[i+2],
			Long:    docs[i+3],
			Example: docs[i+4],
		}
		cmds[docs[i]] = &cmd
		i += 4
	}

	return cmds, i, nil

}

// mapMap is the parser for both the cmd and config sections of the []string
// doc templates
func mapMap(docs []string, idx, l int, what string) (map[string]string, int, error) {
	mStrings := map[string]string{}
	var i int
	seen := map[string]int{}

	for i = idx; i < l; i++ {
		if what == "cmd" && docs[i] == doc.CmdlineDelimiter {
			i++
			break
		}

		if what == "config" && docs[i] == doc.ConfigDelimiter {
			i++
			break
		}

		if what == "custom" && docs[i] == doc.CustomDelimiter {
			i++
			break
		}

		if docs[i] == "" {
			return nil, 0, fmt.Errorf("%s empty variable name on line %d", errText, i)
		}

		if _, s := seen[docs[i]]; s {
			return nil, 0, fmt.Errorf("%s duplicate %s declaration at line %d (previously line %d)", errText, docs[i], i, seen[docs[i]])
		}
		seen[docs[i]] = i

		if i == l-1 {
			return nil, 0, fmt.Errorf("%s odd number of %s section lines %d / %d", errText, what, i, l)
		}

		mStrings[docs[i]] = docs[i+1]
		i++
	}
	return mStrings, i, nil
}
