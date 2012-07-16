jade = require 'jade'

module.exports =
    attach: (opts) ->
        opts.path ?= __dirname + '/../view'
        @render = (fn, data, cb) ->
            jade.renderFile [opts.path, fn + '.jade'].join('/'), data, cb
