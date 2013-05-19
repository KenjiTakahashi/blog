$.spin = (id, self) ->
    data = $(self).data()
    if not data.spinner?
        opts = {
            lines: 10
            length: 0
            corners: 0
            color: "#151515"
            speed: 0.5
            trail: 10
        }
        spinner = new Spinner(opts)
        data.spinner = spinner.spin document.getElementById id

$.unspin = ->
    data = $(@).data()
    if data.spinner?
        data.spinner.stop()
        delete data.spinner
