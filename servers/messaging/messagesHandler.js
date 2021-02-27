// Patch Handler for /v1/messages/{messageID}
const patchMessageHandler = async (req, res, {Message}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  // update the message
  try {
    var messageID = req.url.split('/')[2];
    var newMessage = JSON.parse(req.body);
    await Message.update({"id": messageID}, {$set:{"body": newMessage}});
  } catch (e) {
    res.status(500).send("Something wrong with updating the message");
    return;
  }

  var newMessageObject = Message.findOne({"id": messageID})
  res.setHeader("Content-Type", "application/json");
  res.status(200);
  res.json(newMessageObject);
}

// Delete Handler for /v1/messages/{messageID}
const deleteMessageHandler = async (req, res, {Message}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  // Delete the message body
  try {
    var messageID = req.url.split('/')[2];
    var targetMessage = Message.findOne({"id": messageID});
    if (targetMessage == null) {
      res.status(404).send("No message found");
      return;
    }
    await Message.remove({"id": messageID});
  } catch(e) {
    res.status(500).send("Something wrong with deleting the message");
    return;
  }
  res.status(200).send("Deleted message sucessfully");
}

module.exports = {patchMessageHandler, deleteMessageHandler};