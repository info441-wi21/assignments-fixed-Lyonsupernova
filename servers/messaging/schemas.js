const Schema = require('mongoose').Schema;
var ObjectId = Schema.ObjectId;
const channelSchema = new Schema({
    id: {type: ObjectId, required: true, unique: true, auto: true},
    name: {type: String, unique: true},
    description: {type: String},
    private: {type: Boolean},
    members: {type: [{id: ObjectId, username:String}]},
    createdAt: {type: Date},
    creator: {type: {id: ObjectId, username: String}},
    editedAt: {type: Date}
});

const messageSchema = new Schema({
    id: {type: ObjectId, required: true, unique: true, auto: true},
    channelID: {type: ObjectId, unique: true},
    body: {type: String},
    createdAt: {type: Date},
    creator: {type: {id: ObjectId, username: String}},
    editedAt: {type: Date}
});

module.exports = {channelSchema, messageSchema};