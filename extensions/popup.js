document.addEventListener('DOMContentLoaded', function () {
    var saveButton = document.getElementById('saveButton');
  
    saveButton.addEventListener('click', function () {
      // The saveButton click logic will now handle the action.
      alert('No text selected. Please select a word on the webpage.');
    });
  
    chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
      if (message.action === 'openPopup' && message.selectedText) {
        openPopup(message.selectedText);
      }
    });
  
    function openPopup(selectedText) {
      // TODO: Implement the logic to display the selected word in the popup.
      alert('Selected Word: ' + selectedText);
    }
  });