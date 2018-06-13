package doc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Validate the strings are what they should be to avoid breaking
// compatibility inadvertently. Users should be using
// greenery.DocBaseDelimiter if the strings are declared in .go files, but if
// they are using simple text files to read in we have to make sure these are
// not changed.
func TestUser(t *testing.T) {
	require.Equal(t, BaseDelimiter, "------ DELIMITER:BASE ------")
	require.Equal(t, Use, "Use")
	require.Equal(t, Long, "Long")
	require.Equal(t, Short, "Short")
	require.Equal(t, Example, "Example")
	require.Equal(t, HelpFlag, "HelpFlag")
	require.Equal(t, CmdFlags, "CmdFlags")
	require.Equal(t, ConfigEnvMsg1, "ConfigEnvMsg1")
	require.Equal(t, ConfigEnvMsg2, "ConfigEnvMsg2")
	require.Equal(t, ConfigEnvMsg3, "ConfigEnvMsg3")
	require.Equal(t, HelpDelimiter, "------ DELIMITER:HELP ------")
	require.Equal(t, Usage, "Usage")
	require.Equal(t, Aliases, "Aliases")
	require.Equal(t, Examples, "Examples")
	require.Equal(t, AvailableCommands, "AvailableCommands")
	require.Equal(t, Flags, "Flags")
	require.Equal(t, GlobalFlags, "GlobalFlags")
	require.Equal(t, AdditionalHelpTopics, "AdditionalHelpTopics")
	require.Equal(t, ProvidesMoreInformationAboutACommand, "ProvidesMoreInformationAboutACommand")
	require.Equal(t, UsageDelimiter, "------ DELIMITER:USAGE ------")
	require.Equal(t, ConfigCmd, "config")
	require.Equal(t, ConfigInitCmd, "config>init")
	require.Equal(t, ConfigEnvCmd, "config>env")
	require.Equal(t, ConfigDisplayCmd, "config>display")
	require.Equal(t, HelpCmd, "help")
	require.Equal(t, VersionCmd, "version")
	require.Equal(t, CmdlineDelimiter, "------ DELIMITER:COMMANDLINE ------")
	require.Equal(t, LogLevel, "LogLevel")
	require.Equal(t, ConfFile, "ConfFile")
	require.Equal(t, LogFile, "LogFile")
	require.Equal(t, Pretty, "Pretty")
	require.Equal(t, NoEnv, "NoEnv")
	require.Equal(t, Verbosity, "Verbosity")
	require.Equal(t, DoTrace, "DoTrace")
	require.Equal(t, CfgLocation, "CfgLocation")
	require.Equal(t, CfgForce, "CfgForce")
	require.Equal(t, ConfigDelimiter, "------ DELIMITER:CONFIG ------")
	require.Equal(t, ConfigHeader, ".")
	require.Equal(t, CustomDelimiter, "------ DELIMITER:CUSTOM ------")
}
