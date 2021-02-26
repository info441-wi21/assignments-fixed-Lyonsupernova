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

const {
    channelGetHandler,
    channelPostHandler
} = require('./channelHandler')

const {specificChannelGetHandler,
       specificChannelPostHandler,
       specificChannelPatchHandler,
       specificChannelDeleteHandler} = require('./specificChannelHandler');

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




connect();
// GET: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}



// POST: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}



// PATCH: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}



// DELETE: /v1/channels/{channelID}: refers to a specific channel identified by {channelID}


async function main() {
    app.listen(port, "", () => {
        console.log(`server listening ${port}`)
    })
}
