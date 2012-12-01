fs = require 'fs'
path = require 'path'
coffee = require 'coffee-script'
stylus = require 'stylus'
url = require 'url'

class Assets
    constructor: (@_path) ->
        @assets =
            js: []
            css: []
        @_get 'js', '.coffee', '.js', coffee.compile
        @_get 'css', '.styl', '.css', (str) ->
            ret = undefined
            stylus.render str, (err, css) =>
                ret = css # FIXME: stylus seems to be sync, but it's bad idea
            return ret

    _get: (asset, ext1, ext2, compiler) ->
        fpath = "#{@_path}/#{asset}"
        for file in fs.readdirSync fpath
            content = fs.readFileSync("#{fpath}/#{file}").toString()
            ext = path.extname file
            switch ext
                when ext1
                    @assets[asset][path.basename file, ext] = compiler content
                when ext2
                    @assets[asset][path.basename file, ext] = content

    middleware: (req, res, next) =>
        pathname = url.parse(req.originalUrl).pathname
        elements = pathname.split(path.sep)[-2..]
        asset_type = elements[0]
        asset_name = elements[1]
        assertion = asset_type != '' and asset_name != ''
        if assertion and @assets[asset_type][asset_name]?
            res[asset_type] null, @assets[asset_type][asset_name]
        else
            next()

module.exports = Assets
