// Package greenery is a localizable struct annotation based command-line
// application framework.
/*
Greenery can be used to create localized CLI applications supporting
command-line, environment and configuration-file options.

It is an opinionated porcelain built on top of Cobra
(https://github.com/spf13/cobra) and Viper
(https://github.com/spf13/viper). In greenery, rather than via code, the
configuration variables are defined in a single configuration structure, which
is mapped to the user-visible variables via golang struct annotations.

User commands are invoked with a configuration structure set to the currently
effective configuration, taking into account the command line, the environment
and any specified configuration files.

A localizeable documentation functionality is also provided, together with
some predefined commands that can be used to generate configuration files that
can be used by the application itself.

See https://github.com/woodensquares/greenery for additional usage information.
*/
package greenery
