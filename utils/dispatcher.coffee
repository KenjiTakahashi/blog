helper = (res, number, content_type, data) ->
    res.set 'Content-Type', content_type
    res.send number, data

ifer = (res, err, data, content_type) ->
    if err
        if err == 'INTERNAL'
            res._return.internal err, content_type
        else if err.code == 'ENOENT'
            res._return.not_found err.message, content_type
        else
            res._return.internal err.message, content_type
    else
        res._return.ok data, content_type

module.exports = (req, res, next) ->
    res._return =
        ok: (data, content_type) ->
            helper res, 200, content_type, data
        not_found: (data, content_type) ->
            helper res, 404, content_type, data
        internal: (data, content_type) ->
            helper res, 500, content_type, data
    res.html = (err, data) ->
        ifer res, err, data, 'text/html'
    res.js = (err, data) ->
        ifer res, err, data, 'text/javascript'
    res.css = (err, data) ->
        ifer res, err, data, 'text/css'
    next()
