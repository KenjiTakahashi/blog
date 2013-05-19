hljs = require 'highlight.js'
hljs.tabReplace = '    '

marked = require 'marked'
marked.setOptions gfm: true, highlight: (code, lang) ->
    return hljs.highlight(lang, code).value

module.exports = (data) ->
    data = marked data
    data = data.replace /\[notice\#.*\]/g, (match, def) ->
        e = match.split '#'
        return " <div class='notice'>#{e[1][..-2]}</div>"
    data = data.replace /\[canvas\#[a-zA-Z0-9_]*\#\d*\]/g, (match, def) ->
        e = match[1..-2].split '#'
        return  """
                <div class="notice">
                    <canvas id="#{e[1]}" height="#{e[2]}">
                        no canvas, sorry
                    </canvas>
                    <div>
                        <a href="##{e[1]}#Start" onclick=window.#{e[1]}.init();>start</a>
                        <a href="##{e[1]}#Reset" onclick=window.#{e[1]}.reset();>reset</a>
                        <span><a href="/raw/#{e[1]}" target="_blank">raw</a></span>
                    </div>
                </div>
                <script type="text/javascript" src="/js/#{e[1]}"></script>
                """
    return data
