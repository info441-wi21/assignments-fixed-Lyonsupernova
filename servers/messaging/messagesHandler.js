// Patch Handler for /v1/messages/{messageID}
patchMessageHandler = async (req, res, {Message}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  // update the message
  try {
    const messageID = req.url.split('/')[3];
    const message = JSON.parse(req.body);
    const newMessage = await Message.findOneAndUpdate({"id": messageID}, {$set:{"body": message}});
    res.setHeader("Content-Type", "application/json");
    res.status(200);
    res.json(newMessage);
  } catch (e) {
    res.status(500).send("Something wrong with updating the message");
    return;
  }
}

// Delete Handler for /v1/messages/{messageID}
deleteMessageHandler = async (req, res, {Message}) => {
  if (!req.get('x-user')) {
    res.status(401).send("unauthorized user");
    return;
  }

  // Delete the message body
  try {
    const messageID = req.url.split('/')[3];
    const targetMessage = Message.findOne({"id": messageID});
    if (targetMessage == null) {
      res.status(404).send("No message found");
      return;
    }
    await Message.findOneAndRemove({"id": messageID});
  } catch(e) {
    res.status(500).send("Something wrong with deleting the message");
    return;
  }
  res.status(200).send("Deleted message sucessfully");
}

module.exports = {patchMessageHandler, deleteMessageHandler};