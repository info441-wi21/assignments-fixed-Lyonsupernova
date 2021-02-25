const Schema = require('mongoose').Schema;

const channelSchema = new Schema({
    id: {type: Schema.Type.ObjectId, required = true, unique: true, auto: true},
    name: {type: String, unique: true},
    description: {type: String},
    private: {type: Boolean},
    members: {type: [{id: Schema.Type.ObjectId, username:String}]},
    createdAt: {type: DateTime},
    creator: {type: {id: Schema.Type.ObjectId, username: String}},
    editedAt: {type: DateTime}
});

const messageSchema = new Schema({
    id: {type: Schema.Type.ObjectId, required = true, unique: true, auto: true},
    channelID: {type: Schema.Type.ObjectId, unique: true},
    body: {type: String},
    createdAt: {type: DateTime},
    creator: {type: {id: Schema.Type.ObjectId, username: String}},
    editedAt: {type: DateTime}
});

module.exports = {channelSchema, messageSchema};