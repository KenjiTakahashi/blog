mongo = require 'mongoose'

mongo.connect process.env.MONGO_URI || "mongodb://localhost:27017/test"

PostSchema = new mongo.Schema
    title: type: String, required: yes
    content: type: String
    tags: type: Array

class Post
    initialize: ->
        @_model = mongo.model 'Post', PostSchema

module.exports =
    posts: new Post
