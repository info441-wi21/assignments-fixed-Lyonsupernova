// MySQL connection
const mysql = require('mysql')

var sqlConnection = mysql.createConnection ({
  host: 'localhost',
  user: 'root',
  password: '123456',
  port: '3306',
  database: 'mysqldatabase',
  insecureAuth: true
});

channelGetHandler = async(req, res, {Channel}) => {
    if (!req.get('x-user')) {
        res.status(401).send("unauthorized user");
        return;
    }
    const userID = JSON.parse(req.headers['x-user']);
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
    if (!req.get('x-user')) {
        res.status(401).send("unauthorized user");
        return;
    }
    // parse x-user to get user id
    const userID = JSON.parse(req.headers['x-user']);
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
        // userid => username
        var qry = "SELECT username FROM user WHERE id = " + mysql.escape(userID);
        sqlConnection.query(qry, function (err, result) {
            if (err) {
              console.log('error retrieving new member info:', err.message);
              return;
            }
            // result?
            // result == username?
            console.log(result)
            userName = result[0]
        });
    } catch (e) {
        res.status.send(500).send("There was an issue getting the channels")
        return;
    }
    users = {"id":userID, "username": userName};
    createdAt = new Date();
    creator = users
    //TODO:edited at ? Creator? not found
    const channel = {
        name,
        description,
        private,
        users,
        createdAt,
        users,
        createdAt
    };
    // status code send with json created object
    const query = new Channel(channel);
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send('unable to create a new channel');
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