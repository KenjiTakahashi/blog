package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"regexp"
	"time"

	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
)

var sc = `
.pln { color: #000 }
@media screen {
  .str { color: #080 }
  .kwd { color: #008 }
  .com { color: #800 }
  .typ { color: #606 }
  .lit { color: #066 }
  .pun, .opn, .clo { color: #660 }
  .tag { color: #008 }
  .atn { color: #606 }
  .atv { color: #080 }
  .dec, .var { color: #606 }
  .fun { color: red }
}
@media print, projection {
  .str { color: #060 }
  .kwd { color: #006; font-weight: bold }
  .com { color: #600; font-style: italic }
  .typ { color: #404; font-weight: bold }
  .lit { color: #044 }
  .pun, .opn, .clo { color: #440 }
  .tag { color: #006; font-weight: bold }
  .atn { color: #404 }
  .atv { color: #060 }
}
pre.prettyprint { padding: 2px; border: 1px solid #888 }
ol.linenums { margin-top: 0; margin-bottom: 0 }
li.L0,
li.L1,
li.L2,
li.L3,
li.L5,
li.L6,
li.L7,
li.L8 { list-style-type: none }
li.L1,
li.L3,
li.L5,
li.L7,
li.L9 { background: #eee }
`

var sm = `
html,
body {
	margin: 0px;
	height: 100%;
}
body {
	color: #c8c8c8;
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
#items {
	font-size: 13px;
	list-style: none;
	margin: 4px 0px 0px;
	padding: 0px;
}
a {
	color: #555;
	text-decoration: none;
}
a:hover {
	color: #1c7272;
}
.dead {
	text-decoration: line-through;
}
#post {
	font-family: sans-serif;
	font-size: 14px;
	color: #555;
	text-align: justify;
}
#post a,
#post #title {
	font-weight: 700;
}
#post a {
	text-decoration: underline;
}
pre {
	overflow-x: auto;
}
pre code {
	background: #fff;
	margin: 0px;
	padding: 0px;
	font-size: 11px;
}
canvas {
	background: #fff;
}
.notice {
	background: #c8c8c8;
	padding: 5px;
	margin: 5px 0px;
}
.notice span {
	float: right;
}
img {
	width: 100%;
}
`

var o = `
{{define "i"}}{{end}}
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title>Karol Woźniak aka Kenji Takahashi :: place</title>
		<style>
			{{template "sc"}}
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
			<a target="_blank" href="https://github.com/KenjiTakahashi">github</a> ::
			<a href="/feed">rss</a> ::
			Kenji Takahashi © 2013-2014
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
	<p id="title">{{.CreatedAt | d}} {{.Title}}</p>
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
var rs = regexp.MustCompile(`<code class="language-[a-z]*">(?s).*?<\/code>`)

var tmplFuncs = template.FuncMap{
	"d": func(arg interface{}) string {
		return arg.(time.Time).Format("02 Jan 2006")
	},
	"m": func(arg interface{}) template.HTML {
		bf := blackfriday.MarkdownCommon([]byte(arg.(string)))
		bf = rn.ReplaceAllFunc(bf, func(m []byte) []byte {
			return []byte(fmt.Sprintf(`<div class="notice">%s</div>`, m[8:len(m)-1]))
		})
		bf = rc.ReplaceAllFunc(bf, func(m []byte) []byte {
			ms := bytes.Split(m[8:len(m)-1], []byte("#"))
			return []byte(fmt.Sprintf(c, ms[0], ms[1]))
		})
		bf = ri.ReplaceAllFunc(bf, func(m []byte) []byte {
			m = m[7:len(m)-1]
			return []byte(fmt.Sprintf(`<img src="/assets/image/%[1]s" alt="%[1]s">`, m))
		})
		bf = rs.ReplaceAllFunc(bf, func(m []byte) []byte {
			ms := bytes.SplitN(m, []byte(">"), 2)
			code := []byte(html.UnescapeString(string(ms[1][0:len(ms[1])-7])))
			hl, err := syntaxhighlight.AsHTML(code)
			if err != nil {
				return m
			}
			return []byte(fmt.Sprintf("%s>%s</code>", ms[0], hl))
		})

		return template.HTML(bf)
	},
}

var tmpl = template.Must(
	template.Must(
		template.Must(
			template.New("sc").Parse(sc),
		).New("sm").Parse(sm),
	).New("o").Parse(o),
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
		color: #555;
	}
	a {
		color: #555;
		font-size: 10px;
		text-decoration: none;
		margin-left: 2px;
	}
	a:hover {
		color: #1c7272;
	}
	span {
		color: #c8c8c8;
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
