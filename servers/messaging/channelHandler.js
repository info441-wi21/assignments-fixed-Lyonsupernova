
channelGetHandler = async(req, res, {Channel}) => {
    if (!req.get('X-User')) {
        res.status(401).send("unauthorized user");
        return;
    }
    // parse x-user to get user id
    const user = JSON.parse(req.get('X-User'));
    if (!user || !user['username'] || !user['id']) {
        res.status(401).send("no user found");
        return;
    }
    try {
        res.setHeader("Content-Type", "application/json");
        channels = await Channel.find({$or: [{"private": true, "members.id": user['id']}, {"private": false}]});
        res.json(channels);
    } catch (e) {
        res.status(500).send("channels not found" + user['id'] + users[username]);
    }
};


channelPostHandler = async(req, res, {Channel}) => {
    if (!req.get('X-User')) {
        res.status(401).send("unauthorized user");
        return;
    }
    // parse x-user to get user id
    const user = JSON.parse(req.get('X-User'));
    if (!user || !user['username'] || !user['id']) {
        res.status(401).send("no user found");
        return;
    }

    const{name, description, private} = req.body;
    if (!name) {
        res.status(400).send("no username found");
        return;
    }
    const createdAt = new Date();

    const channel = {
        "name": name,
        "description": description,
        "private": private,
        "members": [user],
        "createdAt": createdAt,
        "creator": user
    };
    // status code send with json created object
    const query = new Channel(channel);
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send('unable to create a new channel' + err);
            return;
        }
        res.status(201).json(newChannel);
        if (req.get('Content-Type') != "application/json") {
            res.setHeader('Content-Type', "application/json");
        }
        return;
    });
};

module.exports = {channelGetHandler, channelPostHandler};