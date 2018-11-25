photoblog
=========

# Overview

Photoblog is a streamlined system for publishing a photo blog as a static site (on s3 or similar). The primary UI is the `blogctl` commandline utility.

It's what I use to make [charczuk.com](http://www.charczuk.com).

# Installation

* Make sure you have golang installed
	- see: [Golang Installation Instructions](https://golang.org/doc/install) for more detail.
* From a commandline, run:
	```bash
	> go install github.com/wcharczuk/photoblog/blogctl
	```

You should now be able to create a blog and publish it.

# `blogctl` usage

With `blogctl` you can initialize a new blog, create a new post, compile the blog to the output directory, and publish it to s3.

See: `blogctl --help` for more info.

# Other stuff

Photoblog uses golang's great `html/template` package to generate the site. More docs on `html/template` can be found [here](https://godoc.org/html/template)
