module github.com/KenjiTakahashi/blog

go 1.14

require (
	github.com/KenjiTakahashi/blog/db v0.0.0-00010101000000-000000000000
	github.com/gorilla/feeds v1.1.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kr/pretty v0.2.1 // indirect
	github.com/russross/blackfriday v1.5.2
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e
)

replace github.com/KenjiTakahashi/blog/db => ./db
