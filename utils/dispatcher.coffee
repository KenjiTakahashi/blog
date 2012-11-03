helper = (res, number, content_type, data) ->
    res.set 'Content-Type', content_type
    res.send number, data

module.exports = (req, res, next) ->
    res._html =
        ok: (data) ->
            helper res, 200, 'text/html', data
        not_found: (data) ->
            helper res, 404, 'text/html', data
        internal: (data) ->
            helper res, 500, 'text/html', data
    res.html = (err, data) ->
        if err
            if err.code == 'ENOENT'
                res._html.not_found err.message
            else
                res._html.internal err.message
        else
            res._html.ok data
    next()
