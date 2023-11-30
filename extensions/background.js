let globalObj = {}

chrome.runtime.onInstalled.addListener(function () {
  chrome.contextMenus.create({
    id: "saveWordContextMenu",
    title: "Define Word",
    contexts: ["selection"],
  });
  return true;
});

// Listener for messages from content.js
chrome.runtime.onMessage.addListener(
  function (request, sender, sendResponse) {
    if (request.action === 'addWord') {
      AddWord(request.word, globalObj);
      
      console.log("globalObj->", globalObj)

      sendResponse(globalObj);
      return false;
    }
    return false;
  }
);



chrome.contextMenus.onClicked.addListener(async function (info, tab) {
  try {
    if (info.menuItemId === "saveWordContextMenu") {
      console.log("I found it: ", info.selectionText, tab.id)

      let word = info.selectionText

      let definitionEng = await getDefinition(word)

      console.log("I found definitionEng: ", definitionEng)

      let data = {
        "word": {
          "lang": "eng",
          "word": word,
          "example": "",
          "image_url": "",
          "link": ""
        },
        "definitions": [
          {
            "lang": "eng",
            "definition": definitionEng
          },
        ]
      };

      chrome.tabs.sendMessage(tab.id, { action: "showOverlay", data: data });
      return true
    }
    return true
  } catch (error) {
    throw error
  }
});
