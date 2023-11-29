chrome.runtime.onInstalled.addListener(function () {
    chrome.contextMenus.create({
      id: "saveWordContextMenu",
      title: "Save Word",
      contexts: ["selection"],
    });
});

chrome.contextMenus.onClicked.addListener(function (info, tab) {
  if (info.menuItemId === "saveWordContextMenu") {
    console.log("I found it: ", info.selectionText, tab.id)
    chrome.tabs.sendMessage(tab.id, { action: "showOverlay" , selectionText: info.selectionText});
  }
});
