path = require 'path'
coffee = require 'coffee-script'
uglifyjs = require 'uglify-js'
url = require 'url'
crypto = require 'crypto'

class Assets
    constructor: (@_db) ->
        @assets = js: [], css: []

    _compile: (ext, data) ->
        if ext == 'coffee'
            data = coffee.compile data
        return uglifyjs.minify(data, fromString: yes).code

    middleware: (req, res, next) =>
        pathname = url.parse(req.originalUrl).pathname
        [atype, aname] = pathname.split(path.sep)[-2..]
        if atype == 'js' and aname != ''
            @_db.one atype, aname, (err, data) =>
                if err or not data
                    next()
                else
                    md5 = crypto.createHash 'md5'
                    newsum = md5.update(data.content).digest 'hex'
                    cond = @assets[atype][aname]?
                    if cond and @assets[atype][aname].sum == newsum
                        res[atype] null, @assets[atype][aname].content
                    else
                        content = @_compile data.type, data.content
                        if atype == 'coffee'
                            atype = 'js'
                        @assets[atype][aname] =
                            content: content
                            sum: newsum
                        res[atype] null, content
        else
            next()

module.exports = Assets
