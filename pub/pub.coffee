#!/usr/bin/env coffee

c = require 'commander'
m = require 'mongoose'
f = require 'fs'
p = require 'path'

c
    .version('0.0.0')
    .option('-g --global', 'post to global (mongolab) db')
    .option('-u --update', 'update an existing item')
    .option('p <file>', 'publish a new post')
    .option('t <name>;<desc>;<site>;<active>', 'add a new project')
    .option('r <file>', 'add a new raw code example')
    .option('a <file>', 'add a new asset code')
    .option('i <file>', 'add a new image')
    .parse(process.argv)

if !c.p and !c.t and !c.r and !c.a and !c.i
    process.exit 0

if c.global
    db = m.connect "mongodb://blog:jakiesbl0ghaslo@ds033267.mongolab.com:33267/blog"
else
    db = m.connect "mongodb://localhost:27017/blog"

ppf = (po, mo) ->
    if mo?
        post = new mo
    else
        post = new Object
    post.short = p.basename po, p.extname po
    data = f.readFileSync(po).toString().split '\n'
    post.title = data[0]
    post.tags = data[1].split ','
    post.content = data[2..].join '\n'
    return post

if c.p
    PostSchema = new m.Schema
        title: type: String, required: yes
        short: type: String, index: true, required: yes
        content: type: String, default: ''
        tags: type: Array, default: []
        date: type: Date, default: new Date
    model = m.model 'Post', PostSchema
    if c.update
        model.find {}, {title: 1}, (err, data) ->
            if err or not data?
                console.log "could not read posts from db - aborting"
                console.log "#{err}"
            else
                titles = (d.title for d in data)
                c.choose titles, (i) ->
                    model.update {title: titles[i]}, ppf(c.p), (err) ->
                        if err
                            console.log "could not update post in db - aborting"
                            console.log "#{err}"
                        else
                            console.log "updated post \\m/"
                        process.exit 0
    else
        post = ppf c.p, model
        post.save (err) ->
            if err
                console.log "could not save post to db - aborting"
                console.log "#{err}"
            else
                console.log "added post \\m/"
            process.exit 0

if c.t
    ProjectSchema = new m.Schema
        name: type: String, required: yes
        desc: type: String, default: ''
        site: type: String, required: yes
        active: type: Boolean, default: true
    model = m.model 'Project', ProjectSchema
    if c.update
        process.exit 1
    else
        project = new model
        data = c.t.split ';'
        project.name = data[0]
        project.desc = data[1]
        project.site = data[2]
        if data[3] == 'true' or data[3] == 'yes'
            project.active = yes
        else if data[3] == 'false' or data[3] == 'no'
            project.active = no
        else
            console.log "setting project to being active"
        project.save (err) ->
            if err
                console.log "could not save project to db - aborting"
                console.log "#{err}"
            else
                console.log "added project \\m/"
            process.exit 0

if c.r
    RawSchema = new m.Schema
        name: type: String, required: yes
        content: type: String, default: ''
    model = m.model 'Raw', RawSchema
    if c.update
        process.exit 1
    else
        raw = new model
        raw.name = p.basename c.r, p.extname c.r
        raw.content = f.readFileSync(c.r).toString()
        raw.save (err) ->
            if err
                console.log "could not save raw to db - aborting"
                console.log "#{err}"
            else
                console.log "added raw \\m/"
            process.exit 0

if c.a
    AssetSchema = new m.Schema
        name: type: String, required: yes
        type: type: String, required: yes
        content: type: String, default: ''
    model = m.model 'Asset', AssetSchema
    if c.update
        process.exit 1
    else
        asset = new model
        ext = p.extname c.a
        asset.name = p.basename c.a, ext
        asset.type = ext[1..]
        asset.content = f.readFileSync(c.a).toString()
        asset.save (err) ->
            if err
                console.log "could not save asset to db - aborting"
                console.log "#{err}"
            else
                console.log "added asset \\m/"
            process.exit 0

if c.i
    ImageSchema = new m.Schema
        name: type: String, index: true, required: yes
        type: type: String, required: yes
        content: type: Buffer, required: yes
    model = m.model 'Image', ImageSchema
    if c.update
        process.exit 1
    else
        image = new model
        ext = p.extname c.i
        image.name = p.basename c.i, ext
        image.type = ext[1..]
        image.content = f.readFileSync c.i
        image.save (err) ->
            if err
                console.log "could not save image to db - aborting"
                console.log "#{err}"
            else
                console.log "added image \\m/"
            process.exit 0
