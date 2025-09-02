# SxWebs - Utilities for Sx web applications

This is a collection of utility functions to build [Sx](https://t73f.de/r/sx)
web applications in [Go](https://go.dev/).

* [SxHTML](/dir?ci=tip&name=sxhtml): Generate HTML from S-Expressions
* [SxHTMLs](/dir?ci=tip&name=sxhtmls): Convert [Webs/htmls](https://t73f.de/r/webs/htmls) to SxHTML.
* [SxHTTP](/dir?ci=tip&name=sxhttp): Encapsulates net/http definitions as Sx objects
* [SxSite](/dir?ci=tip&name=sxsite): Sx code to work with [Webs/Site](https://t73f.de/r/webs)

## Usage instructions

To import this library into your own [Go](https://go.dev/) software, you need
to run the `go get` command. Since Go does not handle non-standard software and
platforms well, some additional steps are required.

First, install the version control system [Fossil](https://fossil-scm.org),
which is a superior alternative to Git in many use cases. Fossil is just a
single executable, nothing more. Make sure it is included in your system's
command search path.

Then, run the following Go command to retrieve a specific version of
this library:

    GOVCS=t73f.de:fossil go get t73f.de/r/sxwebs@HASH

Here, `HASH` represents the commit hash of the version you want to use.

Go currently does not seem to support software versioning for projects managed
by Fossil. This is why the hash value is required. However, this method works
reliably.
