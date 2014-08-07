# viewmaster

viewmaster is a simple view engine for Go. Uses html/template and text/template.
Supports layouts/master pages.

[View API documentation on GoDoc](http://godoc.org/github.com/jpoehls/viewmaster)

## Using

Install the package with:

    go get gopkg.in/jpoehls/viewmaster.v0
  
Import it with:

    import "gopkg.in/jpoehls/viewmaster.v0"
  
and use *viewmaster* as the package name inside the code.

## Credits

<a href='http://www.babygopher.org'><img src='https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png' ></a>

## Wish list

- add support for `optional templates` via a builtin `optional_template` function
- support for nested layouts that use the same template name would be cool
- move to 'github.com/jpoehls/viewengine' package
- add docs
- add some kind of opt-in logging
- future, add caching for production