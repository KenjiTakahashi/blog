require 'date-ext'

date = new Date()
date.setDate 1
days = date.getDaysInMonth()
current_data = [1..days]
month_data = date.getMonth()
year_data = date.getFullYear()
day = date.getDay()
previous_data = []
if day != 6
    date.addDays -1
    previous_days = date.getDaysInMonth()
    previous_data = [previous_days - 5 + day..previous_days]
date.addDays days
day = date.getDay()
next_data = []
if day != 0
    next_data = [1..7 - day]

module.exports =
    month: month_data
    year: year_data
    previous: previous_data
    current: current_data
    next: next_data
