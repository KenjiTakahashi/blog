require 'date-ext'

class Datepicker
    constructor: (db) ->
        @_db = db

    get: (string, callback) ->
        if not string?
            _date = new Date
        else
            _date = new Date string
        _date.setDate 1
        month_data = _date.getMonth() + 1
        year_data = _date.getFullYear()
        day = _date.getDay()
        start = new Date _date.getTime()
        start.addDays -(6 - day)
        end = new Date start.getTime()
        end.addDays 41
        previous_link = "#{start.getFullYear()}-#{start.getMonth() + (day == 6 and 0 or 1)}"
        next_link = "#{end.getFullYear()}-#{end.getMonth() + 1}"
        @_db.dates start, end, (err, data) ->
            if err or not data
                callback err, null
            else
                previous_data = []
                poffset = start.getDate()
                current_data = []
                next_data = []
                end.setDate 1
                for d in data
                    d_ = d.date.getDate()
                    if d.date < _date
                        previous_data[d_ - poffset] = d_
                    else if d.date < end
                        current_data[d_ - 1] = d_
                    else
                        next_data[d_ - 1] = d_
                poffset = start.getDaysInMonth() - poffset
                if not previous_data[poffset]?
                    previous_data[poffset] = undefined
                coffset = start.getDaysInMonth()
                if not current_data[coffset]?
                    current_data[coffset] = undefined
                callback err, {
                    month: month_data
                    year: year_data
                    previous: previous_data
                    current: current_data
                    next: next_data
                    plink: previous_link
                    nlink: next_link
                }

module.exports = Datepicker
