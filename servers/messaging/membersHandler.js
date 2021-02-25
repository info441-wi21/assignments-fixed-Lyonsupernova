/*
Parts that should be add into the index.js:

const {
  postMembersHandler,
  deleteMembersHandler
} = require('./membersHandler');

const handlerWrapper = (handler, forwarder) => {
  return (req, res) => {
    handler(req, res, forwarder);
  }
}

app.post("/v1/channels/{channelID}/members", handlerWrapper(postMembersHandler, Channel));
app.delete("/v1/channels/{channelID}/members", handlerWrappers(postMembersHandler, Channel));
*/
const mysql = require('mysql')

var sqlConnection = mysql.createConnection ({
  host: 'localhost',
  user: 'root',
  password: '123456',
  port: '3306',
  database: 'mysqldatabase',
  insecureAuth: true
});

// post handler for /v1/channels/{channelID}/members
const postMembersHandler = async (req, res, {Channel}) => {
  if (!('X-User' in req.header)) {
    res.status(401).send("unauthorized user");
    return;
  }
  try {
    const channelID = req.url.split('/')[2];
    const userID = JSON.parse(req.headers['x-user']);
    //check creator
    var currDocument = Channel.findOne({id: channelID}, {creator: 1, members: 1});
    var parsed = JSON.parse(currDocument);
    if (parsed.creator != userID) {
      res.status(403).send("Incorrect creator")
    }
  } catch (e) {
    res.status(500).send("Something wrong with validating creator");
    return;
  }

  // Find the member profile in mysql
  const memberID = req.body;
  var newMembers;
  var qry = 'SELECT * from user WHERE id = ' + mysql.escape(memberID);
  sqlConnection.query(qry, function (err, result) {
    if (err) {
      console.log('error retrieving new member info:', err.message);
      return;
    }
    newMembers = parsed.members.push(result[0]);
  });

  await Channel.update({id: channelID},{$set:{members: newMembers}});
  res.status(201).send('Member has been added');
}

// delete handler for /v1/channels/{channelID}/members
const deleteMembersHandler = async (req, res, {Channel}) => {
  if (!('X-User' in req.header)) {
    res.status(401).send("unauthorized user");
    return;
  }

  //check creator
  try {
    const channelID = req.url.split('/')[2];
    const userID = JSON.parse(req.headers['x-user']);
    var currDocument = Channel.findOne({id: channelID}, {creator: 1, members: 1});
    var parsed = JSON.parse(currDocument);
    if (parsed.creator != userID) {
      res.status(403).send("Incorrect creator")
    }
  } catch (e) {
    res.status(500).send("Something wrong with validating creator");
    return;
  }

  // update new members lists
  try {
    const memberID = req.body;
    var members = parsed.members;
    members = members.fileter(member => member.id != memberID)
    await Channel.update({id: channelID},{$set:{members: members}});
  } catch (e) {
    res.status(500).send("Something wrong with updating new member list");
    return;
  }
  res.status(200).send('Member has been deleted');
}

module.exports = {postMembersHandler, deleteMembersHandler};