// initial settings
const mongoose = require('mongoose')
mongoose.set('useFindAndModify', false);

// if mysql connection is needed
/*
const mysql = require('mysql')
var sqlConnection = mysql.createConnection ({
  host: 'localhost',
  user: 'root',
  password: '123456',
  port: '3306',
  database: 'mysqldatabase',
  insecureAuth: true
});
*/

// post handler for /v1/channels/{channelID}/members
postMembersHandler = async(req, res, {Channel}) => {
  if (!req.get('X-User')) {
    res.status(401).send("unauthorized user");
    return;
  }

  const channelID = req.url.split('/')[3];
  const user = JSON.parse(req.headers['X-User']);
  if (!user || !user['username'] || !user['id']) {
      res.status(401).send("no user found");
      return;
  }
  //check creator
  const currDocument = await Channel.findOne({"id": channelID}, {"creator": 1, "members": 1});
  if (currDocument['creator']['id'] != user['id']) {
    res.status(403).send("Incorrect creator")
  }

  // TODO: add another for loop to check if the user is already in the channel
  const member = req.body;
  currDocument['members'] = currDocument['members'].push(member);
  Channel.findOneAndUpdate({"id": channelID},{$set:{"members": currDocument['members']}}, { new: true }, function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
    //res.status(201).send('Member has been added');
    res.json(data);
  });
}

// delete handler for /v1/channels/{channelID}/members
deleteMembersHandler = async (req, res, {Channel}) => {
  if (!req.get('X-User')) {
    res.status(401).send("unauthorized user");
    return;
  }

  //check creator
  const channelID = req.url.split('/')[3];
  const user = JSON.parse(req.headers['X-User']);
  if (!user || !user['username'] || !user['id']) {
      res.status(401).send("no user found");
      return;
  }
  const currDocument = await Channel.findOne({"id": channelID}, {"creator": 1, "members": 1});
  if (currDocument['creator']['id'] != user['id']) {
    res.status(403).send("Incorrect creator");
  }

  // update new members lists
  const memberID = JSON.parse(req.body['id']);
  //res.json(req.body['id']);
  currDocument['members'] = currDocument['members'].filter(el => el['id'] != memberID);
  Channel.findOneAndUpdate(
    {"id": channelID}, {$set:{"members": currDocument['members']}},
    { new: true },
    function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
    //res.status(200).send('Member has been deleted');
    res.json(data);
  });
}

module.exports = {postMembersHandler, deleteMembersHandler};