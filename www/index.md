# SxWebs - Utilities for Sx web applications

This is a collection of utility functions to build [Sx](https://t73f.de/r/sx)
web applications in [Go](https://go.dev/).

* [sxhtml](/dir?name=sxhtml): Generate HTML from S-Expressions
* [sxhttp](/dir?name=sxhttp): Encapsulates net/http definitions as Sx objects
* [sxforms](/dir?name=sxforms): HTML form rendering and validation, similar to [WTForms](https://wtforms.readthedocs.io/)
* [sxsite](/dir?name=sxsite): Sx code to work with [Webs/site](https://t73f.de/r/webs)

## Use instructions

If you want to import this library into your own [Go](https://go.dev/)
software, you must execute a `go get` command. Since Go treats non-standard
software and non-standard platforms quite badly, you must use some non-standard
commands.

First, you must install the version control system
[Fossil](https://fossil-scm.org), which is a superior solution compared to Git,
in too many use cases. It is just a single executable, nothing more. Make sure,
it is in your search path for commands.

How you can execute the following Go command to retrieve a given version of
this library:

    GOVCS=t73f.de:fossil go get t73f.de/r/sxwebs@HASH

where `HASH` is the hash value of the commit you want to use.

Go currently seems not to support software versions when the software is
managed by Fossil. This explains the need for the hash value. However, this
methods works.
