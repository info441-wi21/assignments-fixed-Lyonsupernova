// initial settings
const mongoose = require('mongoose')
mongoose.set('useFindAndModify', false);

// Patch Handler for /v1/messages/{messageID}
patchMessageHandler = async (req, res, {Message}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  // update the message
  const messageID = req.url.split('/')[3];
  const message = req.body['body'];
  const targetMessage = await Message.findOne({"id": messageID});
  if (targetMessage == null) {
    res.status(404).send("No message found");
    return;
  }

  await Message.findOneAndUpdate({"id": messageID}, {$set:{"body": message}}, function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
  });
  res.setHeader("Content-Type", "application/json");
  res.status(200);
    //res.json(newMessage);
}

// Delete Handler for /v1/messages/{messageID}
deleteMessageHandler = async (req, res, {Message}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  // Delete the message body
  const messageID = req.url.split('/')[3];
  const targetMessage = await Message.findOne({"id": messageID});
  if (targetMessage == null) {
    res.status(404).send("No message found");
    return;
  }

  await Message.findOneAndRemove({"id": messageID}, function(err, data) {
    if (err) {
        res.status(400).send("message: " + data + " delete error: " + err);
        return;
    }
  });
  res.status(200).send("Deleted message sucessfully");
}

module.exports = {patchMessageHandler, deleteMessageHandler};