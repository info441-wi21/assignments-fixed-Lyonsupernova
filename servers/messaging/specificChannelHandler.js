specificChannelGetHandler = async function(req, res, {Channel, Message}) {
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
};





specificChannelPostHandler = async function(req, res, {Channel, Message}) {
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
};




specificChannelPatchHandler = async function(req, res, {Channel}) {
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
    const filter = {'id' : channelID};
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
};

specificChannelDeleteHandler = async function(req, res, {Channel, Message}) {
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
        const condition = {'channelID' : channelID};
        Channel.findOneAndDelete(condition);
        Message.deleteMany(condition)
    } catch(e) {
        res.status(500).send("Something wrong with deleting the channel");
        return;
    }
    res.status(200).send("Deleted message sucessfully");
};


Module.exports = {specificChannelGetHandler, specificChannelPostHandler, specificChannelPatchHandler, specificChannelDeleteHandler};