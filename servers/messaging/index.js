const mongoose = require('mongoose')
const express = require('express')
const {channelSchema, messageSchema} = require('./schemas')
const Channel = mongoose.model("Channel", channelSchema)
const Message = mongoose.model("Message", messageSchema)
// const port = ;
const app = express();
app.use(express.json());

const connect = () => {
    mongoose.connect(mongoEndPoint);
}

// MySQL connection


// GET : /v1/channels
app.get('/v1/channels', async(req, res) => {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
        return;
    }
    try {
        res.setHeader("Content-Type", "application/json");
        channels = await Channel.find();
        res.json(channels);
    } catch (e) {
        res.status.send(500).send("There was an issue getting the channels")
        return;
    }
})


// POST : /v1/channels
app.post('/v1/channels', async(req, res) => {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
    }
    // parse x-user to get user id
    const{name, description, private} = req.body;
    createdAt = new Date();
    creator = 
    const channel = {
        "name": name,
        "private": private,
        "members"
    };
    // get user profile from mysql
    const query = new Channel(channel);
    query.save((err, newChannel) => {
        if (err) {
            res.status(500).send('unable to create a new channel');
        }
        res.status(201).send(newChannel);
    });


connect();
async function main() {
    app.listen(port, "", () => {
        console.log(`server listening ${port}`)
    })
}
