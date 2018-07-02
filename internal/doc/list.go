package doc

// BaseDelimiter is documented as part of the non-internal class
const BaseDelimiter = "------ DELIMITER:BASE ------"

// Use is documented as part of the non-internal class
const Use = "Use"

// Long is documented as part of the non-internal class
const Long = "Long"

// Short is documented as part of the non-internal class
const Short = "Short"

// Example is documented as part of the non-internal class
const Example = "Example"

// HelpFlag is documented as part of the non-internal class
const HelpFlag = "HelpFlag"

// CmdFlags is documented as part of the non-internal class
const CmdFlags = "CmdFlags"

// ConfigEnvMsg1 is documented as part of the non-internal class
const ConfigEnvMsg1 = "ConfigEnvMsg1"

// ConfigEnvMsg2 is documented as part of the non-internal class
const ConfigEnvMsg2 = "ConfigEnvMsg2"

// ConfigEnvMsg3 is documented as part of the non-internal class
const ConfigEnvMsg3 = "ConfigEnvMsg3"

// HelpDelimiter is documented as part of the non-internal class
const HelpDelimiter = "------ DELIMITER:HELP ------"

// Usage is documented as part of the non-internal class
const Usage = "Usage"

// Aliases is documented as part of the non-internal class
const Aliases = "Aliases"

// Examples is documented as part of the non-internal class
const Examples = "Examples"

// AvailableCommands is documented as part of the non-internal class
const AvailableCommands = "AvailableCommands"

// Flags is documented as part of the non-internal class
const Flags = "Flags"

// GlobalFlags is documented as part of the non-internal class
const GlobalFlags = "GlobalFlags"

// AdditionalHelpTopics is documented as part of the non-internal class
const AdditionalHelpTopics = "AdditionalHelpTopics"

// ProvidesMoreInformationAboutACommand is documented as part of the
// non-internal class
const ProvidesMoreInformationAboutACommand = "ProvidesMoreInformationAboutACommand"

// UsageDelimiter is documented as part of the non-internal class
const UsageDelimiter = "------ DELIMITER:USAGE ------"

// ConfigCmd is documented as part of the non-internal class
const ConfigCmd = "config"

// ConfigInitCmd is documented as part of the non-internal class
const ConfigInitCmd = "config>init"

// ConfigEnvCmd is documented as part of the non-internal class
const ConfigEnvCmd = "config>env"

// ConfigDisplayCmd is documented as part of the non-internal class
const ConfigDisplayCmd = "config>display"

// HelpCmd is documented as part of the non-internal class
const HelpCmd = "help"

// VersionCmd is documented as part of the non-internal class
const VersionCmd = "version"

// CmdlineDelimiter is documented as part of the non-internal class
const CmdlineDelimiter = "------ DELIMITER:COMMANDLINE ------"

// LogLevel is documented as part of the non-internal class
const LogLevel = "LogLevel"

// ConfFile is documented as part of the non-internal class
const ConfFile = "ConfFile"

// LogFile is documented as part of the non-internal class
const LogFile = "LogFile"

// Pretty is documented as part of the non-internal class
const Pretty = "Pretty"

// NoEnv is documented as part of the non-internal class
const NoEnv = "NoEnv"

// NoEnv is documented as part of the non-internal class
const NoCfg = "NoCfg"

// Verbosity is documented as part of the non-internal class
const Verbosity = "Verbosity"

// DoTrace is documented as part of the non-internal class
const DoTrace = "DoTrace"

// CfgLocation is documented as part of the non-internal class
const CfgLocation = "CfgLocation"

// CfgForce is documented as part of the non-internal class
const CfgForce = "CfgForce"

// ConfigDelimiter is documented as part of the non-internal class
const ConfigDelimiter = "------ DELIMITER:CONFIG ------"

// ConfigHeader is documented as part of the non-internal class
const ConfigHeader = "."

// CustomDelimiter is documented as part of the non-internal class
const CustomDelimiter = "------ DELIMITER:CUSTOM ------"

// DefaultDocs contains all the supported languages as a map of string lists
// following this format. As a library user one can decide to use this format,
// and call greenery.ConvertDocs, or write a DocSet directly.
//
// ------ HELP ------
// help variable name,
// translated help text,
// .....
// ------ USAGE ------
// command name, short translated help text, long translated help text
// command name, short translated,long translated
// ...
// ------ COMMANDLINE ------
// config struct variable name, translated variable commandline help text
// .....
// ------ CONFIG ------
// config block, config block help
// .....
// config struct variable name, translated variable config file help text
// .....
//
// This will be processed once the program starts into internal structs, it
// will panic if not successful so it should be easy to test.
var DefaultDocs = map[string][]string{
	"en": english,
	"it": italian,
}
