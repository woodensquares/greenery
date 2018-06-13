
# NOTE

Greenery has not been officially released yet, it is undergoing final
polishing before the release, see the TODO.md file for more information on
what's still missing. This means that any interfaces are still subject to
change.

# Documentation

Greenery is a library that can be used to create localized CLI applications
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

Full documentation will be available as part of the 1.0 release. Some of the
library has already been documented and can be examined in godoc, which also
contains some end-to-end examples that show its usage.
