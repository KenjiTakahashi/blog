require 'date-ext'

class Datepicker
    constructor: (db) ->
        @_db = db
        @_date = new Date()

    get: (month=@_date.getMonth(), year=@_date.getFullYear()) ->
        @_date.setDate 1
        days = @_date.getDaysInMonth()
        current_data = [1..days]
        month_data = @_date.getMonth()
        year_data = @_date.getFullYear()
        day = @_date.getDay()
        @_date.setMonth @_date.getMonth() - 1
        previous_data = []
        if day != 6
            previous_days = @_date.getDaysInMonth()
            previous_data = [previous_days - 5 + day..previous_days]
        @_date.setMonth @_date.getMonth() + 2
        day = @_date.getDay()
        length = previous_data.length + current_data.length
        next_data = []
        if length < 42
            next_data = [1..42 - length]
        @_date.setMonth @_date.getMonth() - 1
        return {
            month: month
            year: year
            previous: previous_data
            current: current_data
            next: next_data
        }

module.exports = Datepicker
