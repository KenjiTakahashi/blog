class Datepicker
    constructor: () ->
        @_months = [
            'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
            'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
        ]

    get_month_name: (number) ->
        return @_months[number]

    get_string: (date) =>
        return "#{date.getDate()}
                #{@get_month_name(date.getMonth())}
                #{date.getFullYear()}"

module.exports = new Datepicker()
