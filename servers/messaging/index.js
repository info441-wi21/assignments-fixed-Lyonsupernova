const mongoose = require('mongoose')
const express = require('express')
const {channelSchema, messageSchema} = require('./schemas')
const Channel = mongoose.model("Channel", channelSchema)
const Message = mongoose.model("Message", messageSchema)
const port = process.env.PORT;
const mongoPort = process.env.MONGOPORT;
const mongoEndPoint = "mongodb://" + mongoPort;
const app = express();
app.use(express.json());

const connect = () => {
    mongoose.connect(mongoEndPoint);
}

connect();

mongoose.connection.on('error', console.error)
        .on('disconnected', connect)
        .once('open', main);
const {
    channelGetHandler,
    channelPostHandler
} = require('./channelHandler')

const {specificChannelGetHandler,
       specificChannelPostHandler,
       specificChannelPatchHandler,
       specificChannelDeleteHandler} = require('./specificChannelHandler');

const {memberPostHandler,
       memberDeleteHandler} = require('./membersHandler');

const {messagePatchHandler,
       messageDeleteHandler} = require('./messagesHandler');      
const RequestWrapper = (handler, SchemeAndDBForwarder) => {
    return (req, res) => {
        handler(req, res, SchemeAndDBForwarder)
    }
};

app.get("/v1/channels", RequestWrapper(channelGetHandler, { Channel }))
app.post("/v1/channels", RequestWrapper(channelPostHandler, { Channel }))
app.get("/v1/channels/{channelID}", RequestWrapper(specificChannelGetHandler, { Channel, Message }))
app.post("/v1/channels/{channelID}", RequestWrapper(specificChannelPostHandler, { Channel, Message }))
app.patch("/v1/channels/{channelID}", RequestWrapper(specificChannelPatchHandler, { Channel }))
app.delete("/v1/channels/{channelID}", RequestWrapper(specificChannelDeleteHandler, { Channel, Message }))

app.post("/v1/channels/{channelID}/members", RequestWrapper(memberPostHandler, { Channel }))
app.delete("/v1/channels/{channelID}/members", RequestWrapper(memberDeleteHandler, { Channel }))
app.patch("/v1/messages/{messageID}", RequestWrapper(messagePatchHandler, { Message }))
app.delete("/v1/messages/{messageID}", RequestWrapper(messageDeleteHandler, { Message }))


async function main() {
    app.listen(port, "", () => {
        console.log(`server listening ${port}`);
    })
}
