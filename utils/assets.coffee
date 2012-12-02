fs = require 'fs'
path = require 'path'
coffee = require 'coffee-script'
uglifyjs = require 'uglify-js'
cleancss = require 'clean-css'
stylus = require 'stylus'
url = require 'url'
sync = require './sync'

class Assets
    constructor: (@_path) ->
        @assets =
            js: []
            css: []
        @_get 'js', '.coffee', '.js', (str) ->
            return uglifyjs.minify((coffee.compile str), fromString: yes).code
        @_get 'css', '.styl', '.css', (str) ->
            return cleancss.process (sync stylus.render) str

    _get: (asset, ext1, ext2, compiler) ->
        fpath = "#{@_path}/#{asset}"
        try
            for file in fs.readdirSync fpath
                content = fs.readFileSync("#{fpath}/#{file}").toString()
                ext = path.extname file
                base = path.basename file, ext
                switch ext
                    when ext1
                        @assets[asset][base] = compiler content
                    when ext2
                        @assets[asset][base] = content
        catch ENOENT

    middleware: (req, res, next) =>
        pathname = url.parse(req.originalUrl).pathname
        elements = pathname.split(path.sep)[-2..]
        asset_type = elements[0]
        asset_name = elements[1]
        assertion = asset_type in ['js', 'css'] and asset_name != ''
        if assertion and @assets[asset_type][asset_name]?
            res[asset_type] null, @assets[asset_type][asset_name]
        else
            next()

module.exports = Assets
