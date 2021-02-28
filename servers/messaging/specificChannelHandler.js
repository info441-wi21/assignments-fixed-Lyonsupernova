// MySQL connection
const mysql = require('mysql')
const mongoose = require('mongoose')
mongoose.set('useFindAndModify', false);
var sqlConnection = mysql.createConnection ({
  host: 'localhost',
  user: 'root',
  password: '123456',
  port: '3306',
  database: 'mysqldatabase',
  insecureAuth: true
});

specificChannelGetHandler = async function(req, res, {Channel, Message}) {
    if (!req.get('x-user')) {
        res.status(401).send("unauthorized user");
        return;
    }
    const userID = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const channelID = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(400).send("Channel not exist channelID");
        return;
    }
    // If this is a private channel and the current user is not a member, respond with a
    // 403 (Forbidden) status code.
    if (!channel['members'].some(el => el["id"] == userID) && channel.private) {
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
};





specificChannelPostHandler = async function(req, res, {Channel, Message}) {
    if (!req.get('x-user')) {
        res.status(401).send("unauthorized user");
        return;
    }
    const userID = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const channelID = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(400).send("Channel not exist channelID");
        return;
    }
    // If this is a private channel and the current user is not a member, respond with a
    // 403 (Forbidden) status code.
    if (!channel['members'].some(el => el["id"] == userID) && channel.private) {
       res.status(403).send("member not private");
       return;
    }
    // Otherwise, create a new message in this channel using the JSON in the request body.
    // The only message property you should read from the request is body. Set the others based on context.
    const{body} = req.body;
    if (!body) {
        res.status(400).send("no body found");
        return;
    }
    // get user profile from MySQL db
    try {
        userName = "";
        var qry = "SELECT username FROM user WHERE id = " + mysql.escape(userID);
        sqlConnection.query(qry, function (err, result) {
            if (err) {
              console.log('error retrieving new user info:', err.message);
              return;
            }
            userName = result[0]
        });
    } catch (e) {
        res.status.send(500).send("There was an issue getting the user")
        return;
    }
    users = {"id":userID, "username": userName};
    createdAt = new Date();
    const message = {
        "channelID": channelID,
        "body": body,
        "createdAt": createdAt,
        "creator": users
    };
    const query = new Message(message);
    query.save((err, newMessage) => {
        if (err) {
            res.status(400).send('unable to create a new message' + err);
            return;
        }
        if (req.get('Content-Type') != "application/json") {
            res.setHeader('Content-Type', "application/json");
        }
        res.status(201).json(newMessage);
        return;
    });
};




specificChannelPatchHandler = async function(req, res, {Channel}) {
    if (!req.get('x-user')) {
        res.status(401).send("unauthorized user");
        return;
    }
    const userID = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const channelID = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID}).exec();
    if (!channel) {
        res.status(400).send("Channel not exist channelID");
        return;
    }
    // If the current user isn't the creator of this channel,
    // respond with the status code 403 (Forbidden).
    if (channel['creator'].id != userID) {
       res.status(403).send("creator not found");
       return;
    }
    // Otherwise, update only the name and/or description using the JSON in the request body and respond with
    // a copy of the newly-updated channel, encoded as a JSON object.
    const{name, description} = req.body;
    Channel.findOneAndUpdate({"id" : channelID}, { $set: {"name" : name, "description" : description}}, { new: true }, function(err, data) {
        if (err) {
            console.log(err);
        } else if (!data) {
            console.log("data not found" + data);
        } else if (data) {
            console.log(data);
            if (req.get('Content-Type') != "application/json") {
                res.setHeader('Content-Type', "application/json");
            }
            res.json(data);
        }
    });

};

specificChannelDeleteHandler = async function(req, res, {Channel, Message}) {
    if (!req.get('x-user')) {
        res.status(401).send("unauthorized user");
        return;
    }
    const userID = JSON.parse(req.headers['x-user']);
    if (!userID) {
        res.status(401).send("no id found");
        return;
    }
    const channelID = req.params.channelID;
    const channel = await Channel.findOne({"id" : channelID});
    if (!channel) {
        res.status(400).send("Channel not exist channelID");
        return;
    }
    // If the current user isn't the creator of this channel,
    // respond with the status code 403 (Forbidden).
    if (channel['creator'].id != userID) {
       res.status(403).send("the user is not the creator of the channel");
       return;
    }
    // delete the channel and all messages related to it.
    // Respond with a plain text message indicating that the delete was successful.
    Message.deleteMany({'channelID' : channelID}, function(err, data) {
        if (err) {
            res.status(400).send("message: " + data + " delete error: " + err);
            return;
        }
    });
    Channel.findOneAndDelete({'id' : channelID}, function(err, data) {
        if (err) {
            res.status(400).send("channel not find");
            return;
        }
        res.status(200).send("Deleted channel sucessfully: " + data);
    });
};


module.exports = {specificChannelGetHandler, specificChannelPostHandler, specificChannelPatchHandler, specificChannelDeleteHandler};