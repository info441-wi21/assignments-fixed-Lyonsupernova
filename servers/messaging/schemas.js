const Schema = require('mongoose').Schema;
var ObjectId = Schema.ObjectId;
const channelSchema = new Schema({
    id: {type: ObjectId, required: true, unique: true, auto: true},
    name: {type: String, unique: true, required: true},
    description: String,
    private: Boolean,
    members: {type: [{id: Number, username:String}]},
    createdAt: {type: Date, required: true},
    creator: {type: {id: ObjectId, username: String}},
    editedAt: Date,
});

const messageSchema = new Schema({
    id: {type: ObjectId, required: true, unique: true, auto: true},
    channelID: {type: ObjectId, unique: true},
    body: {type: String, required: true},
    createdAt: {type: Date, required: true},
    creator: {type: {id: ObjectId, username: String}},
    editedAt: Date,
});

module.exports = {channelSchema, messageSchema};