mongo = require 'mongoose'

mongo.connect process.env.MONGO_URI || "mongodb://localhost:27017/blog"

PostSchema = new mongo.Schema
    title: type: String, required: yes
    short: type: String, index: true, required: yes
    content: type: String, default: ''
    tags: type: Array, default: []
    date: type: Date, default: new Date

ImageSchema = new mongo.Schema
    name: type: String, index: true, required: yes
    type: type: String, required: yes
    content: type: Buffer, required: yes

ProjectSchema = new mongo.Schema
    name: type: String, required: yes
    desc: type: String, default: ''
    site: type: String, required: yes
    active: type: Boolean, default: true

RawSchema = new mongo.Schema
    name: type: String, required: yes
    content: type: String, default: ''

AssetSchema = new mongo.Schema
    name: type: String, required: yes
    type: type: String, required: yes
    content: type: String, default: ''

class Post
    constructor: ->
        @_model = mongo.model 'Post', PostSchema

    one: (short, callback) ->
        @_model.findOne short: short, callback

    titles: (callback) ->
        query = @_model.find {}, {title: 1, date: 1, short: 1}
        query.sort(date: -1).exec callback

    last20: (callback) ->
        query = @_model.find {}, {title: 1, date: 1}
        query.sort(date: -1).limit(20).exec callback

class Image
    constructor: ->
        @_model = mongo.model 'Image', ImageSchema

    one: (name, callback) ->
        @_model.findOne name: name, callback

class Project
    constructor: ->
        @_model = mongo.model 'Project', ProjectSchema

    all: (callback) ->
        @_model.find {}, {}, {sort: {_id: -1}}, callback

class Raw
    constructor: ->
        @_model = mongo.model 'Raw', RawSchema

    one: (id, callback) ->
        @_model.findOne name: id, callback

class Asset
    constructor: ->
        @_model = mongo.model 'Asset', AssetSchema

    _type: (type) ->
        if type == 'js' or type == 'coffee'
            return ['js', 'coffee']
        if type == 'css' or type == 'styl'
            return ['css', 'styl']

    all_of_type: (type, callback) ->
        @_model.find {type: {$in: @_type type}}, callback

    one: (type, name, callback) ->
        @_model.findOne {type: {$in: @_type type}, name: name}, callback

module.exports =
    posts: new Post
    images: new Image
    projects: new Project
    raws: new Raw
    assets: new Asset
