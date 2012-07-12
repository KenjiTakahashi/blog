helper = (res, number, content_type, data) ->
    res.writeHead number, 'Content-Type': content_type
    res.end data

module.exports =
    attach: (opts) ->
        @html =
            ok: (res, data) ->
                helper res, 200, 'text/html', data
            not_found: (res, data) ->
                helper res, 404, 'text/html', data
            internal: (res, data) ->
                helper res, 500, 'text/html', data
