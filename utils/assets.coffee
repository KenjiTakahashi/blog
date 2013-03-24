fs = require 'fs'
path = require 'path'
coffee = require 'coffee-script'
uglifyjs = require 'uglify-js'
cleancss = require 'clean-css'
stylus = require 'stylus'
url = require 'url'
sync = require './sync'

class Assets
    constructor: (@_path, @_db, callback) ->
        @assets =
            js: []
            css: []
        @_get 'js'
        @_get 'css'
        @_db.all_of_type 'js', (err, data) =>
            for d in data
                @assets['js'][d.name] = @_compile d.type, d.content
            @_db.all_of_type 'css', (err, data) =>
                for d in data
                    @assets['css'][d.name] = @_compile d.type, d.content
                callback()

    _compile: (ext, data) ->
        _ = no
        if ext == 'coffee'
            data = coffee.compile data
            _ = yes
        if _ or ext == 'js'
            return uglifyjs.minify(data, fromString: yes).code
        if ext == 'styl'
            data = (sync stylus.render) data
            _ = yes
        if _ or ext == 'css'
            return cleancss.process data
        return null

    _get: (asset) ->
        fpath = "#{@_path}/#{asset}"
        try
            for file in fs.readdirSync fpath
                if file[-1..] == '~'
                    continue
                content = fs.readFileSync("#{fpath}/#{file}").toString()
                ext = path.extname file
                base = path.basename file, ext
                @assets[asset][base] = @_compile ext[1..], content
        catch ENOENT

    middleware: (req, res, next) =>
        pathname = url.parse(req.originalUrl).pathname
        elements = pathname.split(path.sep)[-2..]
        asset_type = elements[0]
        asset_name = elements[1]
        assertion = asset_type in ['js', 'css'] and asset_name != ''
        if assertion
            if @assets[asset_type][asset_name]?
                res[asset_type] null, @assets[asset_type][asset_name]
            else
                @_db.one asset_type, asset_name, (err, data) =>
                    if err or not data
                        next()
                    else
                        content = @_compile asset_type, data.content
                        if asset_type == 'coffee'
                            asset_type = 'js'
                        if asset_type == 'styl'
                            asset_type = 'css'
                        @assets[asset_type][asset_name] = content
                        res[asset_type] null, content
        else
            next()

module.exports = Assets
