require 'date-ext'

class Datepicker
    constructor: ->
        @month_data = undefined
        @year_data = undefined
        @previous_data = []
        @current_data = []
        @next_data = []
        @_date = new Date()
        @reload()

    reload: ->
        @_date.setDate 1
        days = @_date.getDaysInMonth()
        @current_data = [1..days]
        console.log(@_date)
        @month_data = @_date.getMonth()
        @year_data = @_date.getFullYear()
        day = @_date.getDay()
        @_date.setMonth @_date.getMonth() - 1
        if day != 6
            previous_days = @_date.getDaysInMonth()
            @previous_data = [previous_days - 5 + day..previous_days]
        @_date.setMonth @_date.getMonth() + 2
        day = @_date.getDay()
        length = @previous_data.length + @current_data.length
        if length < 42
            @next_data = [1..42 - length]
        @_date.setMonth @_date.getMonth() - 1

    next_month: ->
        @_date.setMonth @month_data + 1
        @reload()

    previus_month: ->
        @_date.setMonth @month_data - 1
        @reload()

datepicker = new Datepicker()

module.exports =
    month: -> datepicker.month_data + 1
    year: -> datepicker.year_data
    previous: -> datepicker.previous_data
    current: -> datepicker.current_data
    next: -> datepicker.next_data
    forward: -> datepicker.next_month()
    backward: -> datepicker.previous_month()
