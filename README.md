blogctl
=========

# Overview

`blogctl` is a streamlined system for publishing a photo blog as a static site (on s3 or similar). The primary UI is the `blogctl` commandline utility.

It's what I use to make [charczuk.com](http://www.charczuk.com).

# Installation

* Make sure you have golang installed
	- see: [Golang Installation Instructions](https://golang.org/doc/install) for more detail.
* From a commandline, run:
	```bash
	> go get github.com/wcharczuk/blogctl
	```

You should now be able to create a blog and publish it.

# Blog Structure

There are a couple main things that are required to build a blog with `blogctl`.

These are contained within the blog's `config.yml` (found at the root of the blog, a file in YAML encoding):
- `outputPath` Where the blog will be built (defaults to `./dist`)
- `postsPath` Where blogctl reads posts (defaults to `./posts`). Posts should be in their own folder and appear in the order in the blog they appear on disk.
	* A post consists of:
		- The image file (must be a `.jpg`).
		- `meta.yml` Where you can specify things like the posted date, the title, the location, commands and tags.
- `postTemplate` Where the html template for each post lives (defaults to `layout/post.html`)
- `tagTemplate` Where the html template for each tag's posts lives (defaults to `layout/tag.html`)
- `pagesPath` A path to a directory of pages to render (defaults to `layout/pages`). Typically includes `index.html`, or the root page.
- `partialsPath` A path to a directory of partials to include when rendering pages or the `post` or `tag` template.
- `s3` Options for deploying to s3 like the `bucket` and the `region`.
- `cloudfront` Options for caching with `cloudfront`, includes options like the `distribution`.

# `blogctl` usage

With `blogctl` you can initialize a new blog, create a new post, compile the blog to the output directory, and publish it to s3.

It reads from `config.yml` and then files on disk.

Main Commands:
- `blogctl init` Creates a new blog from scratch with a functioning gallery and (1) sample post, and creates a `config.yml` for you.
- `blogctl new` Creates a new post from a given file (must be run in your blog's directory).
- `blogctl build` Compiles posts found in your `postsPath`

See: `blogctl --help` for more info.

# Other stuff

`blogctl` uses golang's great `html/template` package to generate the site. More docs on `html/template` can be found [here](https://godoc.org/html/template)

When rendering templates and pages, `blogctl` also includes the view functions found within [go-sdk's template library](https://github.com/blend/go-sdk/tree/master/template/view_funcs.go).
