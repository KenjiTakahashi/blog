module.exports = (f) ->
    return (args...) ->
        result = null
        done = no
        f.apply null, args.concat [(err, data) ->
            result = data
            done = yes
        ]
        (wait = ->
            if !done
                process.nextTick wait
        )()
        return result
