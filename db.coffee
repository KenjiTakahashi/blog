mongo = require 'mongoose'

mongo.connect process.env.MONGO_URI || "mongodb://localhost:27017/blog"

PostSchema = new mongo.Schema
    title: type: String, required: yes
    content: type: String, default: ''
    tags: type: Array, default: []
    date: type: Date, default: new Date

ProjectSchema = new mongo.Schema
    name: type: String, required: yes
    desc: type: String, default: ''
    active: type: Boolean, default: true

HelperSchema = new mongo.Schema
    latest: [PostSchema]

class Post
    constructor: ->
        @_model = mongo.model 'Post', PostSchema
        @_helper = mongo.model 'Helper', HelperSchema

    latest: (callback) ->
        @_helper.findOne {}, (err, data) ->
            if not data
                callback err, null
            else
                callback err, data.latest[0]

class Project
    constructor: ->
        @_model = mongo.model 'Project', ProjectSchema

    all: (callback) ->
        @_model.find {}, callback

module.exports =
    posts: new Post
    projects: new Project
