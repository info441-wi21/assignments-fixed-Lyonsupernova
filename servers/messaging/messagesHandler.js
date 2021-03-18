// initial settings
const mongoose = require('mongoose')
mongoose.set('useFindAndModify', false);

// Patch Handler for /v1/messages/{messageID}
patchMessageHandler = async (req, res, {Message}) => {
  const user = JSON.parse(req.get('X-User'));
  if (!user) {
    res.status(401).send("unauthorized user");
    return;
  }

  // update the message
  const messageID = req.url.split('/')[3];
  const message = req.body['body'];
  const targetMessage = await Message.findOne({"id": messageID}).exec();
  if (targetMessage == null) {
    res.status(404).send("No message found");
    return;
  }

  // check if the current user is the creator of the message
  // if not, return 403
  if (targetMessage.creator.id != user['id']) {
    res.status(403).send("Cannot patch message: invalid creator");
    return;
  }

  // update the message to mongodb
  await Message.findOneAndUpdate({"id": messageID}, {$set:{"body": message, "editedAt": new Date()}}, function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
  }).exec();

  // retreive updated message and send back to the client side
  const updatedMessage = await Message.findOne({"id": messageID}).exec();
  res.setHeader("Content-Type", "application/json");
  res.status(200).json(updatedMessage);
}

// Delete Handler for /v1/messages/{messageID}
deleteMessageHandler = async (req, res, {Message}) => {
  const user = JSON.parse(req.get('X-User'));
  if (!user) {
    res.status(401).send("unauthorized user");
    return;
  }

  // Delete the message body
  const messageID = req.url.split('/')[3];
  const targetMessage = await Message.findOne({"id": messageID}).exec();
  if (targetMessage == null) {
    res.status(404).send("No message found");
    return;
  }

  // check if the current user is the creator of the message
  // if not, return 403
  if (targetMessage.creator.id != user['id']) {
    res.status(403).send("Cannot patch message: invalid creator");
    return;
  }

  // remove
  await Message.findOneAndRemove({"id": messageID}, function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
  }).exec();
  res.status(200).send("Deleted message sucessfully");
}

module.exports = {patchMessageHandler, deleteMessageHandler};