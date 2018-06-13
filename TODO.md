Overall
-------

Add codecov / circleci / travisci badges

Add godoc badge

Tag 0.9

Documentation
-------------
Finalize the toplevel README

Add missing documentation FIXME:DOC

We don't allow users to shadow our fields, "field collision" in
createBindings"

We are not allowing & | < > for cmdline names which is pretty obvious given
that they are i/o redirects

Tag format needs documentation obviously, also relevant to the above
sepKeyParts...

Only TOML is supported for now

In the config we have at most one level of indirection, don't allow a.b.c.d

The custom logger can be used to interface directly to the underlying user
provided logger via calling Custom

For custom variables the documentation string is not
automatically commented, but it is printed as-is to allow for default values
to be added to the generated configuration

Custom flags need to support flag.Value and encoding.TextMarshaler in order to
work for config.init

Stress that the default language is "en" but can be set in BaseConfigOptions
with the version strings, by default the library is currently providing "en"
and "it"

Possible enhancements
---------------------

CSV flag mapping to []string / []int
