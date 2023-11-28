chrome.runtime.onInstalled.addListener(function () {
    chrome.contextMenus.create({
      id: "saveWordContextMenu",
      title: "Save Word",
      contexts: ["selection"],
    });
  });
  
  chrome.contextMenus.onClicked.addListener(function (info, tab) {
    if (info.menuItemId === "saveWordContextMenu") {
      saveWord(info.selectionText);
    }
  });
  
  function saveWord(word) {
    console.log('Word saved:', word);

  }
  
  function sendSelectedTextToPopup(selectedText) {
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
      chrome.tabs.sendMessage(tabs[0].id, { action: 'openPopup', selectedText: selectedText });
    });
  }