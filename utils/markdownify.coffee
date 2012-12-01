hljs = require 'highlight.js'
hljs.tabReplace = '    '

marked = require 'marked'
marked.setOptions gfm: true, highlight: (code, lang) ->
    ol = ['<ol>']
    get_class = (i) ->
        return i < 9 and 'ten' or (i < 99 and 'hun' or 'tou')
    for line, i in hljs.highlight(lang, code).value.split '\n'
        ol.push "<li class='#{get_class i}'>#{line}</li>"
    ol.push '</ol>'
    return ol.join ''

module.exports = (data) ->
    data = marked data
    data = data.replace /\[canvas\#[a-zA-Z0-9_]*\#\d*\]/g, (match, def) ->
        e = match[1..-2].split '#'
        return  """
                <div class="notice">
                    <canvas id="#{e[1]}" height="#{e[2]}">
                        no canvas, sorry
                    </canvas>
                    <div>
                        <button id="#{e[1]}Start" onclick=window.#{e[1]}.init();>start</button>
                        <button id="#{e[1]}Reset" onclick=window.#{e[1]}.reset();>reset</button>
                        <span><a href="/raw/#{e[1]}" target="_blank">raw</a></span>
                    </div>
                </div>
                <script type="text/javascript" src="/js/#{e[1]}.js"></script>
                """
    return data
