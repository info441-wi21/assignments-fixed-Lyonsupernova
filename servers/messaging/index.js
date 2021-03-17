const mongoose = require('mongoose')
const express = require('express')
const {channelSchema, messageSchema} = require('./schemas')
const Channel = mongoose.model("Channel", channelSchema)
const Message = mongoose.model("Message", messageSchema)
//const port = 4000
//const mongoPort = process.env.MONGOPORT;
//const mongoEndPoint = "mongodb://localhost:27017/test"
const mongoEndPoint = process.env.MONGOADDR
const port = process.env.PORT;
const app = express();
app.use(express.json());

const connect = () => {
    mongoose.connect(mongoEndPoint, {useNewUrlParser:true});
}
connect();

app.listen(port, "", () => {
    console.log(`server is listening ${port}`);
})
const {
    channelGetHandler,
    channelPostHandler
} = require('./channelHandler')

const {specificChannelGetHandler,
       specificChannelPostHandler,
       specificChannelPatchHandler,
       specificChannelDeleteHandler} = require('./specificChannelHandler');

const {postMembersHandler,
       deleteMembersHandler} = require('./membersHandler');

const {patchMessageHandler,
       deleteMessageHandler} = require('./messagesHandler');
const RequestWrapper = (handler, SchemeAndDBForwarder) => {
    return (req, res) => {
        handler(req, res, SchemeAndDBForwarder);
    }
};

// Define the method not allow function
const methodNotAllowed = (req, res, next) => res.status(405).send("Method not allowed");

app
.route('/v1/channels')
.get(RequestWrapper(channelGetHandler, { Channel }))
.post(RequestWrapper(channelPostHandler, { Channel }))
.all(methodNotAllowed);

app
.route("/v1/channels/:channelID")
.get(RequestWrapper(specificChannelGetHandler, { Channel, Message }))
.post(RequestWrapper(specificChannelPostHandler, { Channel, Message }))
.patch(RequestWrapper(specificChannelPatchHandler, { Channel }))
.delete(RequestWrapper(specificChannelDeleteHandler, { Channel, Message }))
.all(methodNotAllowed);

app
.route("/v1/channels/:channelID/members")
.post(RequestWrapper(postMembersHandler, { Channel }))
.delete(RequestWrapper(deleteMembersHandler, { Channel }))
.all(methodNotAllowed);

app
.route("/v1/messages/:messageID")
.patch(RequestWrapper(patchMessageHandler, { Message }))
.delete(RequestWrapper(deleteMessageHandler, { Message }))
.all(methodNotAllowed);