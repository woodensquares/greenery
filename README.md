
[![Build Status](https://travis-ci.org/woodensquares/greenery.svg "Travis CI status")](https://travis-ci.org/woodensquares/greenery)
[![Coverage status](https://codecov.io/github/woodensquares/greenery/branch/master/graph/badge.svg)](https://codecov.io/github/woodensquares/greenery?branch=master)
[![GoDoc](https://godoc.org/github.com/woodensquares/greenery?status.svg)](https://godoc.org/github.com/woodensquares/greenery)


# NOTE

Greenery has not been officially released yet, it is undergoing final
polishing before the 1.0 release, this means that any interfaces are still
subject to change. Please do feel free to open issues if you think anything in
the library would benefit from changing, especially in terms of interfaces and
functionality.

# Overview

Greenery is a framework that can be used to create localized CLI applications
supporting command-line, environment and configuration-file options.

It is an opinionated porcelain built on top of
[Cobra](https://github.com/spf13/cobra) and
[Viper](https://github.com/spf13/viper). In greenery, rather than via code,
the configuration variables are defined in a single configuration structure,
which is mapped to the user-visible variables via golang struct annotations.

User commands are invoked with a configuration structure set to the currently
effective configuration, taking into account the command line, the environment
and any specified configuration files.

A localizeable documentation functionality is also provided, together with
some predefined commands that can be used to generate configuration files that
can be used by the application itself.

# Table of Contents

- [Overview](#overview)
- [Installing](#installing)
- [Concepts](#concepts)
  * [Configuration](#configuration)
  * [Commands](#commands)
  * [Flags](#flags)
  * [Arguments](#arguments)
- [Examples](#examples)
  * [A minimal example](#a-minimal-example)
  * [A localized example](#a-localized-example)
- [Reference](#reference)
  * [The annotation format](#the-annotation-format)
  * [Localization](#localization)
  * [Logging](#logging)
  * [Tracing](#tracing)
  * [Configuration files](#configuration-files)

# Installing

To install greenery simply run:

    go get -u github.com/woodensquares/greenery/...

to install both the library, its test helper and zap logging back-end
additional functionality. After this you are ready to go.

# Concepts

Greenery operates in the following manner: given a command line with flags and
arguments, the environment and optionally a TOML configuration file, it will
take all supported ways to specify flags, convert those in configuration
struct values, and pass this effective configuration, together with any
command line arguments, to the function implementing the command requested by
the user.

## Configuration

The configuration struct, embedding the base configuration provided by the
greenery library, is what should control the execution of the program. Handler
functions will have access to the effective configuration.

## Commands

Commands are the commands that the user will choose on the command line to
represent an action they want to execute, which will be mapped to a handler
implementing it.

Commands in greenery are referenced both by the config struct annotation,
which controls the binding between commands and flags, by the map containing
information about which command calls which handler function, and the
documentation passed to the library.

Command names must not contain spaces or any of the & , | < characters. The <
character is allowed as the last character of the command name, to signify
this command will take arguments, see below for more information.

The > character is used to define subcommand relationships, so a command named
"hat>fedora" in a "make" program, would map to a "make hat fedora" command
line (where there could be a separate "hat>baseball" command defined, mapping
to "make hat baseball").

## Flags

Flags are options that affect user commands, in greenery flags are values that
must fulfill the flag.Value interface. In order for greenery to operate
correctly in terms of creating and parsing the configuration for non-standard
types, the encoding.TextMarshaler and encoding.TextUnmarshaler interfaces are
also taken into account respectively when generating configuration files, and
parsing configuration/environment values.

## Arguments

Arguments are additional arguments the user will put on the command line after
the specified command, if a command is meant to support arguments it should be
declared with a < character.

# Examples

These examples are also included in the godoc documentation, linked above,
which contains additional examples.

## A minimal example

The minimal example, as its name implies, shows how to use the greenery
library to create probably the simplest application possible, with one command
and one flag.

```go
package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/woodensquares/greenery"
)
```
After the preamble let's create the configuration for an application that we
assume will execute some REST GET request against a server, in this case we'd
like the application to allow us to specify a URI to connect to, and to have a
timeout for the request.

The application will have the URI as a command line parameter, and will have
the timeout as an option that can be set via the commandline, environment or
configuration file.

```go
type exampleMinimalConfig struct {
    *greenery.BaseConfig
    Timeout *greenery.IntValue `greenery:"get|timeout|t,      minimal-app.timeout,   TIMEOUT"`
}
```

our configuration struct embeds the BaseConfig struct available in the
library, and adds our parameter. The parameter is defined via a *greenery*
struct annotation in a [format described below](#the-annotation-format).

With this definition, our flag will be available as --timeout/-t on the
command line for the get parameter, via the MINIMAL_TIMEOUT environmental
variable, and the *timeout* parameter in the configuration file under the
[minimal-app] section.


```go
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
```

typically one would want to have a function that creates an initialized
configuration struct, with any default values. In our case this function will
create a struct with a default timeout value of 400, and set its limits to
between 0 and 1000 via the provided IntValue flag.

The base configuration is also initialized by passing "minimal" as our
application name, and declaring that we have a "get" command that takes a
commandline argument (the "<" character) and will call the
*exampleMinimalGetter* function if invoked.

In order for the library to be able to display help, we should now declare the
documentation for our program

```go
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
```

for simplicity's sake we are just providing the English documentation, and
we'll be using the struct format for the documentation, rather than the string
list format. See [later in the document](#localization) for an
example of a string list documentation.

Our application will just have a short help format both for itself and for the
get command, it also will reuse the command line help for the config file
help, which means we just have to specify CmdLine.

It is now time to create the actual logic for our application, which in this
case will simply validate that the user did pass an argument to the command,
and will print out what it would be doing with the effective timeout, as well
as emitting a debug message.

```go
func exampleMinimalGetter(lcfg greenery.Config, args []string) error {
    cfg := lcfg.(*exampleMinimalConfig)

    if len(args) != 1 {
        return fmt.Errorf("Invalid number of command line arguments")
    }

    cfg.Debugf("fetching %s, timeout %d", args[0], cfg.Timeout.Value)
    fmt.Printf("Will fetch %s with timeout %d milliseconds\n\n", args[0], cfg.Timeout.Value)
    return nil
}
```

The main function for our application simply creates our configuration, defers
a Cleanup call, and executes it.

```go
func main() {
    cfg := exampleNewMinimalConfig()
    defer cfg.Cleanup()

    if err := cfg.Execute(cfg, exampleMinimalDocs); err != nil {
        fmt.Fprintf(os.Stderr, "Error executing: %s\n", err)
        os.Exit(1)
    }
}
```

Let's now look at how this minimal application works

## Invocations

Let's first try to run without passing any parameter
```
~: go run minimal.go 
URL fetcher

Usage:
  minimal [command]

Available Commands:
  config      Configuration file related commands
  get         Retrieves the specified page
  help        Help about any command
  version     Prints out the version number of the program

Flags:
  -c, --config string      The configuration file location
      --help               help information for the application.
      --log-file string    The log file location
  -l, --log-level string   The log level of the program. Valid values are "error", "warn", "info" and "debug" (default "error")
      --no-cfg             If set no configuration file will be loaded
      --no-env             If set the environment variables will not be considered
      --pretty             If set the console output of the logging calls will be prettified
  -v, --verbosity int      The verbosity of the program, an integer between 0 and 3 inclusive. (default 1)

"minimal [command] --help" provides more information about a command.
```

As you can see we will be getting the help for the application, with our
**get** command as well as the available default commands and options.

Our typical invocation would be simply to execute the get command

```
~: go run minimal.go get http://somewhere.com
Will fetch http://somewhere.com with timeout 400 milliseconds
```

as you can see the timeout used will be our default, this can be changed of
course via the commandline

```
go run minimal.go get --timeout 100 http://somewhere.com
Will fetch http://somewhere.com with timeout 100 milliseconds
```

or the environment

```
~: export MINIMAL_TIMEOUT=440
~: go run minimal.go get http://somewhere.com
Will fetch http://somewhere.com with timeout 440 milliseconds
```

we can also create a configuration file and set it there, let's create a
configuration file in our current directory

```
~: go run minimal.go config init
Configuration file generated at ........../minimal.toml
~: cat minimal.toml 
# Configuration generated on 2018-06-30T14:07:46-07:00

# The log file location
log-file = ""
# The log level of the program. Valid values are "error", "warn", "info" and "debug"
log-level = "error"
# If set the environment variables will not be considered
no-env = false
# If set the console output of the logging calls will be prettified
pretty = false
# The verbosity of the program, an integer between 0 and 3 inclusive.
verbosity = 1

[minimal-app]
# the timeout to use for the fetch
timeout = 400
```

as you can see default values for all the internal flags, as well as our
timeout flag, have been added. If we now edit this file and provide a
different value, it will be honored

```
~: cat minimal.toml | grep timeout
# the timeout to use for the fetch
timeout = 300
~: go run minimal.go get http://somewhere.com
Will fetch http://somewhere.com with timeout 300 milliseconds
```

the application automatically will look for a configuration file named
applicationname.toml in the various XDG configuration directories, as well as
in the current directory.

## A localized example

The localized example is meant to be exactly the same as the minimal example,
but showing how it could be localized. For fun it is showing a pig latin
localization, that can be seen by setting LANG to *zz.ZZZ*

This sample, unlike the minimal sample above, shows the "list of strings" way
of specifying a documentation set, this can be useful if non-developers will
be in charge of the documentation directly, for example the default language
documentation rather than the struct above will be

```go
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
```

note in this case we are also showing that custom strings can be added to the
documentation for localization of code-emitted messages, the getter function
in fact in this sample has the following printf call

```go
func exampleLocalizedGetter(lcfg greenery.Config, args []string) error {

    // ...

    fmt.Printf(docs.Custom["Message"], args[0], cfg.Timeout.Value)
	return nil
}
```

the largest section of this file is the pig latin localization, once these two
localizations are created, they can be converted to the usual struct format
via 

```go
var exampleLocalizedDocs = greenery.ConvertOrPanic(map[string][]string{
	"en":     exampleLocalizedEnglish,
	"zz.ZZZ": exampleLocalizedPiglatin,
})
```

in this case we are using the ConvertOrPanic function, which is typically used
for package variable definitions.

# Reference

## The annotation format

As shown above in the sample code, greenery operates by mapping exported
struct values to config/environment/commandline flags.

Each exported struct member in the config struct is expected to contain a
greenery annotation, unexported members will be ignored, and it is not legal
to override any of the default base variables in the embedding struct: this
will generate an error.

Greenery struct tags are composed of three parts separated by commas

```go
    `greenery:"get|timeout|t,      minimal-app.timeout,   TIMEOUT"`
```

### Commandline

The first part of the annotation controls the commandline behavior of the
flag, it is composed by three parts separated by pipe characters: the first
part is the command this flag belongs to, the second is the long name of the
option, and the third is a single letter used as a shortname.

If the long or short name for the options are not present, they will not be
available. If the command name is not present, in general the flag is
configuration/environment only, in that case the first part should be
specified as **||none** or **||custom** depending if the flag in question
requires custon TOML deserialization.

### Config file

The second part of the annotation controls the name of the option in the
configuration file, it is composed by two parts separated by a period. The
first part is the header of the section this option belongs to, in this case a
section named *[minimal-app]*, while the second is the name of the
configuration file variable corresponding to the flag.

If this part of the annotation is not present, the flag is not going to be
accessible via the configuration file. If the section name is not present, the
flag will be part of the base section of the configuration file together with
the other base flags like log level and so on.

### Environment

The third part of the annotation controls the name of the environmental
variable corresponding to the option. It will be accessible on the environment
via NAMEOFTHEAPP_[name] where NAMEOFTHEAPP is the uppercased name of the
application as specified in the **greenery.NewBaseConfig** call.

If this part of the annotation is not present, the flag is not going to be
available via the environment.

### Precedence

The precedence of flags is command line overrides environment overrides
configuration file.

## Localization

The default language for the library is "en" but can be set in
BaseConfigOptions with the version strings, by default the library is
currently providing "en" and "it" translations for the default commands.


In order to simplify localization efforts by non-coders, besides writing the
wanted help documentation in a DocSet struct format, it is possible to write
it as a list of strings with identifiers, this should hopefully enable easy
editing by non-go-developers who can be simply told they should change only
text between double quotes.

The format is as follows

    Section Delimiter (base, help, cmd, )
    Section content
    Section Delimiter
    ...

Each section can appear only once in the list of strings in any order. Each
section has a separate set of internal identifiers discussed below.  Complete
examples are available in the unit tests as well as in doc/internal/
DefaultDocs contains all the supported languages as a map of string lists
following this format. As a library user one can decide to use this format,
and call greenery.ConvertDocs, or write a DocSet directly.

    ------ HELP ------
    help variable name,
    translated help text,
    .....
    ------ USAGE ------
    command name, short translated help text, long translated help text
    command name, short translated,long translated
    ...
    ------ COMMANDLINE ------
    config struct variable name, translated variable commandline help text
    .....
    ------ CONFIG ------
    config block, config block help
    .....
    config struct variable name, translated variable config file help text
    .....

In order to make it easy to differentiate between library-related strings and
user provided strings, constants have been defined containing all the relevant
strings mentioned above.

Please see the localized sample code for a complete example that will fully
localize any default provided flag or message.

## Logging

Logging calls with custom severities and verbosities are supported, by default
a very basic stderr logger will service any greenery logging calls, however a
higher performance and more full featured layer is provided via
*greenery/zapbackend* that will use the excellent
https://github.com/uber-go/zap/ logger.

It is possible also to use an arbitrary provided logger to service these calls
as long as it implements the provided Logger interface.

## Tracing

Tracing, controlled by the --trace / APPNAME_TRACE flags, or programmatical
Start/StopTracing calls, is a special very verbose logging level which in
addition to any user code tracing calls will be printing a lot of internal
greenery tracing information.

Trace cannot be set via config file, only via the environment or the command
line. Note that for tracing there is no override, if it's set in the command
line OR environment it will be turned on.

## Configuration files

Applications created with the greenery library will support configuration
files by default, if this is not desired the *--no-cfg / APPNAME_NOCFG* flags
can be used.

Only TOML-formatted files are supported for now, the the application
automatically will look for a configuration file named applicationname.toml in
the various XDG configuration directories, as well as in the current
directory, or in the location specified via the *--config / -c* flag.
