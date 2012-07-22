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
    site: type: String, requored: yes
    active: type: Boolean, default: true

HelperSchema = new mongo.Schema
    latest: [PostSchema]

class Post
    constructor: ->
        @_model = mongo.model 'Post', PostSchema
        @_helper = mongo.model 'Helper', HelperSchema

    latest: (callback) ->
        @_helper.findOne (err, data) ->
            if err or not data
                callback err, null
            else
                callback err, data.latest[0]

    tags: (o..., callback) ->
        o = o[0]
        query = @_model.find({}, {tags: 1}).sort('date', -1)
        query.exec (err, data) ->
            if err or not data
                callback err, null
            else
                tmp = []
                out = []
                for post in data
                    for tag in post.tags
                        if tag in tmp
                            out[tmp.indexOf tag][1] += 1
                        else
                            out.push [tag, 1]
                            tmp.push tag
                has_prev = true
                if not o? or o == 0
                    has_prev = false
                has_next = true
                #if count <= 12
                    #has_next = false
                callback err, out, has_prev, has_next

    titles: (o..., callback) ->
        o = o[0]
        if not o?
            o = 0
        query = @_model.find({}, {title: 1}).sort('date', -1).skip(o).limit 10
        query.exec (err, data) ->
            query.count (err, count) ->
                count -= o
                if err or not data
                    callback err, null
                else
                    has_prev = true
                    if o == 0
                        has_prev = false
                    has_next = true
                    if count <= 10
                        has_next = false
                    callback err, (d.title for d in data), has_prev, has_next

class Project
    constructor: ->
        @_model = mongo.model 'Project', ProjectSchema

    all: (callback) ->
        @_model.find {}, callback

module.exports =
    posts: new Post
    projects: new Project
