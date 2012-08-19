class Datepicker
    constructor: (db) ->
        @_db = db

    get: (string, callback) ->
        if not string?
            _date = new Date
        else
            _date = new Date string
        year_data = _date.getFullYear()
        start = new Date _date.getFullYear(), 0, 1
        end = new Date _date.getFullYear(), 11, 31
        @_db.dates start, end, (err, data) ->
            if err or not data
                callback err, null
            else
                current_data = []
                for d in data
                    d_ = d.date.getMonth()
                    current_data[d_] = d_ + 1
                if not current_data[11]?
                    current_data[11] = undefined
                callback err, {
                    year: year_data
                    data: current_data
                }

module.exports = Datepicker
