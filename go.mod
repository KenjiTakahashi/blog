module github.com/KenjiTakahashi/blog

go 1.14

require (
	github.com/KenjiTakahashi/blog/db v0.0.0-00010101000000-000000000000
	github.com/alecthomas/chroma v0.7.2-0.20200305040604-4f3623dce67a
	github.com/gorilla/feeds v1.1.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kr/pretty v0.2.1 // indirect
	github.com/yuin/goldmark v1.2.1
	github.com/yuin/goldmark-highlighting v0.0.0-20200307114337-60d527fdb691
)

replace github.com/KenjiTakahashi/blog/db => ./db
