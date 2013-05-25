hljs = require 'highlight.js'
hljs.tabReplace = '    '

marked = require 'marked'
marked.setOptions gfm: true, highlight: (code, lang) ->
    return hljs.highlight(lang, code).value

module.exports = (data) ->
    data = marked data
    data = data.replace /\[notice\#.*\]/g, (match, def) ->
        e = match[8..-2]
        return " <div class='notice'>#{e}</div>"
    data = data.replace /\[canvas\#[a-zA-Z0-9_]*\#\d*\]/g, (match, def) ->
        e = match[8..-2].split '#'
        return  """
                <div class="notice">
                    <canvas id="#{e[0]}" height="#{e[1]}">
                        no canvas, sorry
                    </canvas>
                    <div>
                        <a href="##{e[0]}#Start" onclick=window.#{e[0]}.init();>start</a>
                        <a href="##{e[0]}#Reset" onclick=window.#{e[0]}.reset();>reset</a>
                        <span><a href="/raw/#{e[0]}" target="_blank">raw</a></span>
                    </div>
                </div>
                <script type="text/javascript" src="/js/#{e[0]}"></script>
                """
    data = data.replace /\[image\#[a-zA-Z0-9_]*\]/g, (match, def) ->
        e = match[7..-2]
        return  """<img src="/images/#{e}" alt="#{e}">"""
    return data
