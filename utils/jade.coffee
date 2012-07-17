jade = require 'jade'

module.exports =
    attach: (opts) ->
        opts.path = "#{__dirname}/../views/#{opts.path ? ""}"
        @render = (fn, data, cb) ->
            jade.renderFile [opts.path, "#{fn}.jade"].join('/'), data, cb
