helper = (res, number, content_type, data) ->
    res.writeHead number, 'Content-Type': content_type
    res.end data

module.exports = (req, res) ->
    res.html =
        ok: (data) ->
            helper res, 200, 'text/html', data
        not_found: (data) ->
            helper res, 404, 'text/html', data
        internal: (data) ->
            helper res, 500, 'text/html', data
    res.emit 'next'
