const mysql = require('mysql')
const mongoose = require('mongoose')
mongoose.set('useFindAndModify', false);

// if mysql connection is needed
/*
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
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  const channelID = req.url.split('/')[3];
  const userID = JSON.parse(req.headers['x-user']);
  if (!userID) {
  res.status(401).send("no id found");
    return;
  }
  //check creator
  const currDocument = await Channel.findOne({"id": channelID}, {"creator": 1, "members": 1});
  if (currDocument['creator']['id'] != userID) {
    res.status(403).send("Incorrect creator")
  }

  // TODO: add another for loop to check if the user is already in the channel
  const member = req.body;
  currDocument['members'] = currDocument['members'].push(member);
  await Channel.findOneAndUpdate({"id": channelID},{$set:{"members": currDocument['members']}});
  res.status(201).send('Member has been added');
}

// delete handler for /v1/channels/{channelID}/members
deleteMembersHandler = async (req, res, {Channel}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  //check creator
  const channelID = req.url.split('/')[3];
  const userID = JSON.parse(req.headers['x-user']);
  const currDocument = await Channel.findOne({"id": channelID}, {"creator": 1, "members": 1});
  if (currDocument['creator']['id'] != userID) {
    res.status(403).send("Incorrect creator");
  }

  // update new members lists
  const memberID = JSON.parse(req.body['id']);
  //res.json(req.body['id']);
  currDocument['members'] = currDocument['members'].filter(el => el['id'] != memberID);
  const newChannel = await Channel.findOneAndUpdate(
    {"id": channelID}, {$set:{"members": currDocument['members']}},
    function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
  });
  res.json(newChannel);
  res.status(200).send('Member has been deleted');
}

module.exports = {postMembersHandler, deleteMembersHandler};