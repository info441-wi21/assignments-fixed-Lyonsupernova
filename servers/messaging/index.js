const mongoose = require('mongoose')
const express = require('express')
const {channelSchema, messageSchema} = require('./schemas')
const Channel = mongoose.model("Channel", channelSchema)
const Message = mongoose.model("Message", messageSchema)
const port = 4000
// const mongoPort = process.env.MONGOPORT;
const mongoEndPoint = "mongodb://localhost:27017/test"
// const mongoEndPoint = "mongodb://" + mongoPort;
// const port = process.env.PORT;
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

const {postMembersHandler,
       deleteMembersHandler} = require('./membersHandler');

const {patchMessageHandler,
       deleteMessageHandler} = require('./messagesHandler');
const RequestWrapper = (handler, SchemeAndDBForwarder) => {
    return (req, res) => {
        handler(req, res, SchemeAndDBForwarder);
    }
};

app.get("/v1/channels", RequestWrapper(channelGetHandler, { Channel }))
app.post("/v1/channels", RequestWrapper(channelPostHandler, { Channel }))
app.get("/v1/channels/:channelID", RequestWrapper(specificChannelGetHandler, { Channel, Message }))
app.post("/v1/channels/:channelID", RequestWrapper(specificChannelPostHandler, { Channel, Message }))
app.patch("/v1/channels/:channelID", RequestWrapper(specificChannelPatchHandler, { Channel }))
app.delete("/v1/channels/:channelID", RequestWrapper(specificChannelDeleteHandler, { Channel, Message }))

app.post("/v1/channels/:channelID/members", RequestWrapper(postMembersHandler, { Channel }))
app.delete("/v1/channels/:channelID/members", RequestWrapper(deleteMembersHandler, { Channel }))
app.patch("/v1/messages/:messageID", RequestWrapper(patchMessageHandler, { Message }))
app.delete("/v1/messages/:messageID", RequestWrapper(deleteMessageHandler, { Message }))


async function main() {
    app.listen(port, "", () => {
        console.log(`server is listening ${port}`);
    })
}
