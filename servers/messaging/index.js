const mongoose = require('mongoose')
const express = require('express')
const mysql = require('mysql')
const {channelSchema, messageSchema} = require('./schemas')
const Channel = mongoose.model("Channel", channelSchema)
const Message = mongoose.model("Message", messageSchema)
// const port = ;
const mongoEndPoint = "mongodb://localhost:27017/";
const app = express();
app.use(express.json());

const connect = () => {
    mongoose.connect(mongoEndPoint);
}

// MySQL connection
const sqlPool = mysql.createPool({
    host: mysqlEndPoint[0],
    user: mysqlEndPoint[1],
    password: mysqlEndPoint[2],
    database: mysqlEndPoint[3],
    insecureAuth: true
});
// query MySQL
async function querySQL(query) {
    return new Promise(function (resolve, reject) {
        sqlPool.query(query, async function(err, result, field) {
            if (err) {
                return reject(err);
            }
            resolve(result);
        })
    })
}

// GET : /v1/channels
app.get('/v1/channels', async(req, res) => {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    try {
        res.setHeader("Content-Type", "application/json");
        //TODO: what await find() will return, a list of obj?
        channels = await Channel.find();
        res.json(channels);
    } catch (e) {
        res.status(500).send("channels not found");
    }
})


// POST : /v1/channels
app.post('/v1/channels', async(req, res) => {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
    }
    // parse x-user to get user id

    
    const {userID} = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    //TODO: since the description is optional, what if the lack of description will lead to it will store value of private
    //TODO: do we need to pass in memebr id or member name in request body
    const{name, description, private} = req.body;
    if (!name) {
        res.status(400).send("no username found");
        return;
    }
    // get user profile from MySQL db
    try {
        userName = "";
        var userNameFromDB = await querySQL("SELECT username FROM user WHERE id = " + mysql.escape(userID));
        if (userNameFromDB) {
            userName = userNameFromDB[0].userName;
        } else {
            res.status(500).send("username not found");
        }
    } catch (e) {
        res.status.send(500).send("There was an issue getting the channels")
        return;
    }
    users = {"id":userID, "username": userName};
    createdAt = new Date();
    creator = users
    const channel = {
        "name": name,
        "description": description,
        "private": private,
        "members": users,
        "createdAt": createdAt,
        "creator": users
    };
    // status code send with json created object
    const query = new Channel(channel);
    //TODO: is here we stored in the DB?
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send('unable to create a new channel');
            return;
        }
        res.status(201).json(newChannel);
        res.setHeader("Content-Type", "application/json");
        return;
    });
})

connect();
// GET: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}
app.get('/v1/channel/{channelID}', async function(req, res) {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
    }
    // parse x-user to get user id
    const {userID} = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const {channelID} = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(401).send("Channel not exist");
        return;
    }
    // If this is a private channel and the current user is not a member, respond with a 
    // 403 (Forbidden) status code.
    if (!channel['members'].some(el => el.id == userID) && channel.private) {
       res.status(403).send("member not private");
       return; 
    } 
    // Otherwise, respond with the most recent 100 messages posted to the specified channel, 
    try {
        //TODO: messaged most recent 100
        const messages = await Message.find({"channelID":channelID}).sort({"id":-1}).limit(100);
        // encoded as a JSON array of message model objects. 
        res.json(messages);
        // Include a Content-Type header set to application/json so that your client knows what sort of data 
        // is in the response body.
        res.setHeader("Content-Type", "application/json");
    } catch (e) {
        console.log("message not found most recent 100");
    }
})





// POST: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}
app.post('/v1/channel/{channelID}', async function(req, res) {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
    }
    // parse x-user to get user id
    const {userID} = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const {channelID} = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(401).send("Channel not exist");
        return;
    }
    // If this is a private channel and the current user is not a member, respond with a 
    // 403 (Forbidden) status code.
    if (!channel['members'].some(el => el.id == userID) && channel.private) {
       res.status(403).send("member not private");
       return; 
    } 
    // Otherwise, create a new message in this channel using the JSON in the request body.
    // The only message property you should read from the request is body. Set the others based on context.
    const{body} = req.body;
    // get user profile from MySQL db
    try {
        userName = "";
        var userNameFromDB = await querySQL("SELECT username FROM user WHERE id = " + mysql.escape(userID));
        if (userNameFromDB) {
            userName = userNameFromDB[0].userName;
        } else {
            res.status(500).send("username not found");
        }
    } catch (e) {
        res.status.send(500).send("There was an issue getting the channels")
        return;
    }
    users = {"id":userID, "username": userName};
    createdAt = new Date();
    creator = users
    const message = {
        "channelID": channelID,
        "body": body,
        "createdAt": createdAt,
        "creator": users,
        "editedAt": createdAt
    };
    const query = new Message(message);
    query.save((err, newMessage) => {
        if (err) {
            res.status(400).send('unable to create a new channel');
            return;
        }
        res.status(201).json(newMessage);
        res.setHeader("Content-Type", "application/json");
        return;
    });
})


// PATCH: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}
app.patch('/v1/channel/{channelID}', async function(req, res) {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
    }
    // parse x-user to get user id
    const {userID} = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const {channelID} = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(401).send("Channel not exist");
        return;
    }
    // If the current user isn't the creator of this channel,
    // respond with the status code 403 (Forbidden).
    if (!channel['creator'].id != userID) {
       res.status(403).send("creator not found");
       return; 
    } 
    // Otherwise, update only the name and/or description using the JSON in the request body and respond with
    // a copy of the newly-updated channel, encoded as a JSON object. 
    const{name, description} = req.body;
    //TODO: how to respond a new copy of updated channel? the id would be auto-incremented and be unique
    //TODO: findOneAndUpdate?
    const filter = {id : channelID};
    const update = {'name' : name, 'description' : description};
    const updateChannel = Channel.findOneAndUpdate(filter, update, {
        new: true
    });
    // status code send with json created object
    const query = new Channel(updateChannel);
    //TODO: is here we stored in the DB?
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send('unable to create a new channel');
            return;
        }
        res.json(newChannel);
        res.setHeader("Content-Type", "application/json");
        return;
    });
})


// DELETE: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}
app.patch('/v1/channel/{channelID}', async function(req, res) {
    if (!('X-User' in req.header)) {
        res.status(401).send("unauthorized user");
    }
    // parse x-user to get user id
    const {userID} = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const {channelID} = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(401).send("Channel not exist");
        return;
    }
    // If the current user isn't the creator of this channel,
    // respond with the status code 403 (Forbidden).
    if (!channel['creator'].id != userID) {
       res.status(403).send("creator not found");
       return; 
    } 
    // delete the channel and all messages related to it. 
    // Respond with a plain text message indicating that the delete was successful.
    try {
        const channel = await Channel.findOne({"id" : channelID});
        const condition = {channelID : channelID};
        Message.findByIdAndDelete(condition)
        Channel.findOneAndDelete(condition);
        const query = new Channel(updateChannel);
    } catch(e) {
        res.status(500).send("Something wrong with deleting the channel");
        return;
    }
    res.status(200).send("Deleted channel sucessfully");
});


async function main() {
    app.listen(port, "", () => {
        console.log(`server listening ${port}`)
    })
}
