require 'date-ext'

class Datepicker
    constructor: (db) ->
        @_db = db

    get: (string) ->
        if not string?
            _date = new Date()
        else
            _date = new Date string
        _date.setDate 1
        days = _date.getDaysInMonth()
        current_data = [1..days]
        month_data = _date.getMonth()
        year_data = _date.getFullYear()
        day = _date.getDay()
        _date.setMonth _date.getMonth() - 1
        previous_link = "#{_date.getFullYear()}-#{_date.getMonth() + 1}"
        previous_data = []
        if day != 6
            previous_days = _date.getDaysInMonth()
            previous_data = [previous_days - 5 + day..previous_days]
        _date.setMonth _date.getMonth() + 2
        next_link = "#{_date.getFullYear()}-#{_date.getMonth() + 1}"
        day = _date.getDay()
        length = previous_data.length + current_data.length
        next_data = []
        if length < 42
            next_data = [1..42 - length]
        _date.setMonth _date.getMonth() - 1
        return {
            month: month_data + 1
            year: year_data
            previous: previous_data
            current: current_data
            next: next_data
            plink: previous_link
            nlink: next_link
        }

module.exports = Datepicker
