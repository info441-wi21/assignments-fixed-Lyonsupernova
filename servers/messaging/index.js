const mongoose = require('mongoose')
const express = require('express')
const mysql = require('mysql')
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
    // get user profile from mysql
    const query = new Channel(channel);
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send('unable to create a new channel');
            return;
        }
        res.status(201).send(newChannel);
        return;
    });


connect();
async function main() {
    app.listen(port, "", () => {
        console.log(`server listening ${port}`)
    })
}
