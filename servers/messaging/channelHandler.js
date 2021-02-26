channelGetHandler = async(req, res, {Channel}) => {
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
};


channelPostHandler = async(req, res, {Channel}) => {
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
};

Module.exports = {channelGetHandler, channelPostHandler};