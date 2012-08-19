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
    latest: type: mongo.Schema.ObjectId, ref: 'Post'

class Post
    constructor: ->
        @_model = mongo.model 'Post', PostSchema
        @_helper = mongo.model 'Helper', HelperSchema

    one: (id, callback) ->
        @_model.findOne _id: id, (err, data) ->
            if err or not data
                callback err, null
            else
                callback err, data

    latest: (callback) ->
        @_helper.findOne({}, {latest: 1}).populate('latest').exec (err, data) ->
            if err or not data
                callback err, null
            else
                callback err, data.latest

    dates: (start, end, callback) ->
        start.setHours 0, 0, 0, 0
        end.setDate end.getDate() + 1
        end.setHours 0, 0, 0, 0
        query = @_model.find {date: {$gte: start, $lt: end}}, {date: 1}
        query.sort(date: 1).exec (err, data) ->
            if err or not data
                callback err, null
            else
                callback err, data

    tags: (o, callback) ->
        if not o?
            o = 0
        query = @_model.find({}, {tags: 1}).sort date: -1
        query.exec (err, data) ->
            if err or not data
                callback err, null
            else
                tmp = []
                out = []
                for post in data
                    for tag in post.tags
                        if tag in tmp
                            index = (tmp.indexOf tag) - o
                            if out[index]?
                                out[index][1] += 1
                        else if tmp.length < o
                            tmp.push tag
                        else
                            out.push [tag, 1]
                            tmp.push tag
                has_prev = true
                if o == 0
                    has_prev = false
                has_next = true
                if tmp.length - o <= 12
                    has_next = false
                callback err, out.slice(0, 12), has_prev, has_next

    titles: (d, t, o, callback) ->
        if not o?
            o = 0
        if d?
            start = new Date d
            end = new Date d
            end.setMonth end.getMonth() + 1
            query = @_model.find {date: {$gte: start, $lt: end}}, {title: 1}
        else if t?
            query = @_model.find {tags: t}, {title: 1}
        else
            query = @_model.find {}, {title: 1}
        query = query.sort(date: -1).skip(o).limit 10
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
                    callback err, data, has_prev, has_next

class Project
    constructor: ->
        @_model = mongo.model 'Project', ProjectSchema

    all: (callback) ->
        @_model.find {}, callback

module.exports =
    posts: new Post
    projects: new Project
