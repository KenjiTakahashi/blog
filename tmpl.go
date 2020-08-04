package main

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

var sm = `
html,
body {
	margin: 0px;
	height: 100%;
}
body {
	color: #b7b7b7;
	font-family: monospace;
}
#o {
	height: 100%;
	min-height: 100%;
	display: table;
	width: 100%;
}
#i {
	display: table-cell;
	vertical-align: middle;
	padding: 10px 0px 14px;
	width: 70%;
}
.m {
	max-width: 15%;
}
#foot {
	clear: both;
	margin-top: -14px;
	position: relative;
	right: 10px;
	bottom: 6px;
	float: right;
}
#head {
	font-size: 16px;
}
#post #title,
#items {
	font-size: 13px;
	list-style: none;
	margin: 4px 0px 0px;
	padding: 0px;
}
#post #title a,
a {
	color: #454545;
	text-decoration: none;
}
#post #title a:hover,
a:hover {
	color: #1d7878;
}
.dead {
	text-decoration: line-through;
}
#post {
	font-family: sans-serif;
	font-size: 14px;
	color: #454545;
	text-align: justify;
}
#post a {
	font-weight: bold;
	text-decoration: underline;
}
#post #title {
	font-family: monospace;
	color: #b7b7b7;
}
#post #title a {
	font-weight: normal;
}
canvas {
	background: #fff;
}
.notice {
	background: #b7b7b7;
	padding: 5px;
	margin: 5px 0px;
}
.notice span {
	float: right;
}
img {
	width: 100%;
}
pre {
	overflow-x: auto;
	background: #fff;
	margin: 0px;
	padding: 0px;
	font-size: 11px;
}
.code {
	margin-left: -1.9rem;
}
.code td:first-of-type span {
	display: block;
	width: 1.24rem;
	text-align: right;
}
`

var o = `
{{define "i"}}{{end}}
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title>Karol Woźniak aka Kenji Takahashi :: place</title>
		<link rel="alternate" type="application/rss+xml" title="Karol Woźniak aka Kenji Takahashi :: rss" href="/feed/rss" />
		<link rel="alternate" type="application/atom+xml" title="Karol Woźniak aka Kenji Takahashi :: atom" href="/feed/atom" />
		<style>
			{{template "sm"}}
		</style>
	</head>
	<body>
		<div id="o">
			<div class="m"></div>
			<div id="i">
				<div id="head">
					<a href="/posts">posts</a> ::
					<a href="/projects">projects</a> ::
					<a href="/">Karol Woźniak</a> aka
					<a href="/">Kenji Takahashi</a>
				</div>
				{{template "i" .}}
			</div>
			<div class="m"></div>
		</div>
		<div id="foot">
			<a target="_blank" href="https://linkedin.com/in/wozniakkarol">linkedin</a> ::
			<a target="_blank" href="https://github.com/KenjiTakahashi">github</a> ::
			Kenji Takahashi © 2013-2014,2016,2020
		</div>
	</body>
</html>
`

var r = `
<ul id="items">
	{{range .}}
	<li>{{.CreatedAt | d}} <a href="/posts/{{.Short}}">{{.Title}}</a></li>
	{{end}}
</ul>
`

var p = `
<div id="post">
	<p id="title">{{.CreatedAt | d}} <a href="/posts/{{.Short}}">{{.Title}}</a></p>
	{{.Content | m}}
</div>
`

var t = `
<ul id="items">
	{{range .}}
	<li>
		<a href="{{.Site}}"{{if (not .Active)}} class="dead"{{end}}>{{.Name}}</a>
		{{.Description}}
	</li>
	{{end}}
</ul>
`

var c = `
<div class="notice">
	<canvas id="%[1]s" height="%[2]s">
		no canvas, sorry
	</canvas>
	<div>
		<a href="#%[1]s#Start" onclick=window.%[1]s.init();>start</a>
		<a href="#%[1]s#Reset" onclick=window.%[1]s.reset();>reset</a>
		<span><a href="/assets/raw/%[1]s" target="_blank">raw</a></span>
	</div>
</div>
<script type="text/javascript" src="/assets/script/%[1]s"></script>
`

var rn = regexp.MustCompile(`\[notice\#.*\]`)
var rc = regexp.MustCompile(`\[canvas\#[a-zA-Z0-9_]*\#\d*\]`)
var ri = regexp.MustCompile(`\[image\#[a-zA-Z0-9_]*\]`)

var md = goldmark.New(
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("lovelace"),
			highlighting.WithFormatOptions(
				html.WithLineNumbers(true),
				html.LineNumbersInTable(true),
			),
			highlighting.WithWrapperRenderer(func(w util.BufWriter, context highlighting.CodeBlockContext, entering bool) {
				if entering {
					w.Write([]byte(`<div class="code">`))
				} else {
					w.Write([]byte(`</div>`))
				}
			}),
		),
	),
	goldmark.WithRendererOptions(
		gmhtml.WithUnsafe(),
	),
)

var tmplFuncs = template.FuncMap{
	"d": func(arg interface{}) string {
		return arg.(time.Time).Format("02 Jan 2006")
	},
	"m": func(arg interface{}) template.HTML {
		var bfb bytes.Buffer
		md.Convert([]byte(arg.(string)), &bfb)
		bf := bfb.Bytes()
		bf = rn.ReplaceAllFunc(bf, func(m []byte) []byte {
			return []byte(fmt.Sprintf(`<div class="notice">%s</div>`, m[8:len(m)-1]))
		})
		bf = rc.ReplaceAllFunc(bf, func(m []byte) []byte {
			ms := bytes.Split(m[8:len(m)-1], []byte("#"))
			return []byte(fmt.Sprintf(c, ms[0], ms[1]))
		})
		bf = ri.ReplaceAllFunc(bf, func(m []byte) []byte {
			m = m[7 : len(m)-1]
			return []byte(fmt.Sprintf(`<img src="/assets/image/%[1]s" alt="%[1]s">`, m))
		})

		return template.HTML(bf)
	},
}

var tmpl = template.Must(
	template.Must(
		template.New("sm").Parse(sm),
	).New("o").Funcs(tmplFuncs).Parse(o),
)

var e = `
<html>
<head>
	<title>Karol Woźniak aka Kenji Takahashi :: error</title>
	<style>
	body {
		position: fixed;
		top: 46%;
		left: 36%;
		font-family: monospace;
		font-size: 43px;
		color: #454545;
	}
	a {
		color: #454545;
		font-size: 10px;
		text-decoration: none;
		margin-left: 2px;
	}
	a:hover {
		color: #1d7878;
	}
	span {
		color: #b7b7b7;
	}
	</style>
</head>
<body>
	code <span>::</span> {{.}}
	<br/>
	<a href="/">go back to normal</a>
</body>
</html>
`

var etmpl = template.Must(template.New("e").Parse(e))
