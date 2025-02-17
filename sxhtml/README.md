# SxHTML - Generate HTML from S-Expressions

HTML can be represented as a symbolic expression, often referred to
as an [s-expression](https://en.wikipedia.org/wiki/S-expression)
or simply a sexpr (for short). This is approach is similar to
[SXML](https://en.wikipedia.org/wiki/SXML), which is an attempt to represent
XML as s-expressions.

For example, consider the following simple HTML:

    <html>
      <head><title>Example</title></head>
      <body>
        <h1 id="main">Title</h1>
        <p>This is some example text.</p>
        <hr>
        <div class="small" id="footnote">Small text.</div>
      </body>
    </html>

The corresponding s-expression representation would be:

    (html
      (head (title "Example"))
      (body
        (h1 (@ (id "main)) "Title")
        (p "This is some example text.")
        (hr)
        (div (@ (class "small") (id "footnote)) "Small text.")
      )
    )

The s-expression representation offers the advantage of easier parsing compared
to HTML text. Additionally, s-expressions can be more straightforward to
analyze and potentially optimize, as they provide a structured format. For
example, a `((p) (p))` can be simplified to `((p))`. Similarly, in certain
cases, a `(li (p "text"))` can be transformed into `(li "text")`.

This library enables the generation of HTML from s-expressions created by
[Sx](https://t73f.de/r/sx).

HTML is often generated using string template libraries such as
[Mustache](https://mustache.github.io/) (available in many programming
languages), [Jinja](https://jinja.palletsprojects.com/) (Python), or
[html/template](https://pkg.go.dev/html/template) (Go).

One common challenge is _escaping_ certain characters that have special
meanings in various parts of HTML. For example, the less-than character "`<`"
marks the beginning of a tag and cannot appear literally in normal text; it
must be replaced by "`&lt;`". Similary, the ampersand character "`&`" also has
special significance and must be replaced with "`&amp;`". However, this is
only true for regular HTML content. Within HTML attributes (such as "href" in
"`<a href="...">...</a>`"), different rules apply. Additionally, if you embed
JavaScript in your HTML content, there is yet another set of escaping rules to
follow.

Most string template libraries fall short in certain scenarios. For example,
Mustache provides escape sequences only for HTML content, not for HTML
attributes. The same limitation applies to Jinja. The html/template library in
Go requires the developer to explicitly specify the correct escaping mode.

This is because string template libraries work at the string level, where the
structure of the HTML is lost.

In contrast, by using a structured representation of HTML, the HTML generator
can understand the context and automatically select the appropriate escape
mode.

## Language

SxHTML is based on [Sx](https://t73f.de/r/sx).

SxHTML is fairly lenient regarding the supported HTML language, but it is
primarily targeted for HTML5. All tag and attribute names must be lowercase
symbols. Do not use strings to specify tags or attributes. SxHTML does not
validate whether a symbol corresponds to a valid HTML tag or attribute.
Additionally, certain tag and attribute symbols have special meanings.

<https://html.spec.whatwg.org/multipage/syntax.html#void-elements> specifies
the list of _void elements_ that does not have and end tag. All other tags will
haven an end tag.

<https://html.spec.whatwg.org/multipage/indices.html#attributes-1> associates
attribute names with expected content. This will result in an additional
escaping mechanism for specific content type. Currently, only URL content is
recognized and escaped.

In addition to the list above, the are some heuristics in detecting content
type based on the attribute name.

* A prefix of "data-" is stripped. For example, `data-href` is also treated as
  an URL attribute.
* If there is no "data-" prefix, any namespace prefix is stripped. For example,
  `svg:href` is also treated as an URL attribute, but not `svg:data-href`.
* The namespace "xmlns" will always result in treating the attribute as an URL
  attribute, e.g. `xmlns:svg`.
* If the attribute name contains one of the strings "url", "uri", "src", it
  will be treated as an URL attribute.
* If the attribute name starts with "on", it will be treated in future versions
  as JavaScript.
* An attribute name "style" will treat the attribute value as CSS in the
  future.

SxHTML defines some additional symbols, all starting with "@":

* `@` specifies the attribute list of an HTML tag. If must follow immediately
  the tag symbol and contains a list of pairs, where the first component is a
  symbol and the second component is a string, the nil value, or a number.
* `@C` marks some content that should be written as `<![CDATA[...]]>`.
* `@H` specifies some HTML content that must not be escaped. For example,
  `(@H "&amp;")` is transformed to `&amp;`, but not `&amp;amp;`.
* `@L` contains elements that just just be transformed, without specifying a
  tag. It is used by generating software that wants to generate HTML for a
  sequence of elements that do not belong to a certain tag.
* `@@` specifies a HTML comment, e.g. `(@@ "comment")` is transformed to
 `<!-- comment -->`.
* `@@@` specifies a multiline HTML comment, e.g. `(@@@ "line1" "line2")` is
  transformed to `\n<!--\nline1\nline2\n-->\n`.
* `@@@@` specifies the doctype statement, e.g. `(@@@@ (html ...))` is
  transformed to `<!DOCTYPE html>\n<html>...</html>`.

## Tags

HTML defines some tags as *void elements*. A void element has no content,
they have a start tag only. End tags must not be specified, SxHTML will not
generated them. Any content except attributes are ignored. Void elements are:
`area`, `base`, `br`, `col`, `embed`, `hr`, `img`, `input`, `link`, `meta`,
`source`, `track`, and `wbr`.

## Attributes

Attributes are always in the second position of a list containing a tag
symbol. For example `(a (@ (href . "https://t73f.de/r/sxhtml")) "SxHTML)`
specifies a link to the page of this library. It will be transformed to
`<a href="https://t73f.de/r/sxhtml">SxHTML</a>`.

The syntax for attributes is as follows:

* The first element of the attribute list must be the symbol `@`.
* Remaining elements must be lists, where the first element of each list is a
  symbol, which names the attribute.
* If there is no second element in the list, the attribute is an *empty
  attribute*. For example, `(input (@ (disabled)))` will be transformed to
  `<input disabled>`,
* If there is a second element in the list, it must be an atomic value,
  preferably a string. For example, `(input (@ (disabled "yes")))` will be
  transformed to `<input disabled="yes">`.
* If the lists contains more elements, they are ignored.
* if the list is a pair, the second element of the pair must be an atomic
  value, preferably a string. For example, `(input (@ (disabled . "yes")))`
  will be transformed to `<input disabled="yes">`.

Since the attribute list is just a list, there might be duplicate symbols
as attribute names. Only the first occurrence of the symbol will create an
attribute. For example, `(input (@ (disabled "no") (disabled . "yes")))` will
be transformed to `<input disabled="no">`. This allows to extend the list
of attributes at the front, if you later want to overwrite the value of an
attribute.

If you want to prohibit the generation of some attribute while still extending
the list of attributes at the front, use the nil value *()* as the value of the
attribute. For example, `(input (@ (disabled ()) (disabled . "yes")))` will be
transformed to `<input>`.

## Content

HTML is not just about tags and attributes; they are essential for structuring
content. To specify content, it’s best to use strings, though numbers
are also allowed without the need to convert them to strings. Other
[Sx](https://t73f.de/r/sx) types, such as symbols, vectors, and undefined
values, are simply ignored.
