[![GoDoc](https://godoc.org/github.com/Drachenfels-GmbH/go-ini?status.svg)](https://godoc.org/github.com/Drachenfels-GmbH/go-ini)

# go-ini

Simple but powerful INI file parser.
This parser helps you to access duplicate sections and keys.

## TODO

### Comment processing

* Associate preceding (multi-line) comments with lines ?
* Associate trailing comments to lines ?
* For KeyVal split up the Value part into Value and CommentValue

Map comments to line numbers and use line number of line to collect
comments (same line, continuous lines before) ?

