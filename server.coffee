flatiron = require 'flatiron'
director = require 'director'

app = flatiron.app
app.use flatiron.plugins.http
app.use require './utils/jade'
app.use require './utils/dispatcher'

routes =
    '/':
        get: -> app.html.ok @res, 'Hello World!'

app.router = new director.http.Router(routes).configure async: true

app.start process.env.app_port || 8080
