module github.com/KenjiTakahashi/blog/pub

go 1.14

require (
	github.com/KenjiTakahashi/blog/db v0.0.0-00010101000000-000000000000
	github.com/mitchellh/cli v1.1.1
	github.com/tidwall/buntdb v1.1.2
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
)

replace github.com/KenjiTakahashi/blog/db => ../db
