let globalObj = {}
function highlightWordsInPageWithTooltip() {
  // Get the text content of the entire page
  const pageText = document.body.innerText;

  if (globalObj == null || globalObj == undefined || globalObj == {}) {
    return
  }
  // Iterate through each key in the globalObj root
  for (const key in globalObj) {
    if (globalObj.hasOwnProperty(key)) {
      const wordToHighlight = key;
      const regex = new RegExp(`\\b${wordToHighlight}\\b`, 'gi');

      // Replace each occurrence of the word with a span and add a tooltip
      document.body.innerHTML = document.body.innerHTML.replace(regex, match => {
        const tooltipContent = getTooltipContent(globalObj[key]);
        return `<span class="highlighted-word" title="${tooltipContent}" style="background-color: yellow; cursor: pointer;">${match}</span>`;
      });
    }
  }
}

function getTooltipContent(data) {
  // Customize this function based on the structure of your JSON data
  // Here, we are assuming a simple structure for demonstration purposes
  if (data && data.eng && data.eng.example) {
    return `Example: ${data.eng.example}`;
  } else {
    return "No additional information available";
  }
}

// Add event listener to show tooltip on hover
document.body.addEventListener('mouseover', function (event) {
  const target = event.target;

  if (target.classList.contains('highlighted-word')) {
    const tooltipContent = target.getAttribute('title');
    showTooltip(event.clientX, event.clientY, tooltipContent);
  }
});

document.body.addEventListener('mouseout', function (event) {
  const target = event.target;

  if (target.classList.contains('highlighted-word')) {
    hideTooltip();
  }
});

function showTooltip(x, y, content) {
  const tooltip = document.createElement('div');
  tooltip.className = 'tooltip';
  tooltip.innerHTML = content;
  tooltip.style.position = 'absolute';
  tooltip.style.top = y + 'px';
  tooltip.style.left = x + 'px';
  document.body.appendChild(tooltip);
}

function hideTooltip() {
  const tooltips = document.querySelectorAll('.tooltip');
  tooltips.forEach(tooltip => tooltip.remove());
}


window.addEventListener("load", function () {
  highlightWordsInPageWithTooltip();
});

window.addEventListener("scroll", function () {
  highlightWordsInPageWithTooltip();
});

document.addEventListener("click", function (event) {
  console.log("User clicked at coordinates:", {
    x: event.clientX,
    y: event.clientY,
  });
  hidePopup();
});

// Attach the handleAddButtonClick function to the click event of the "popup-add" button
document.addEventListener('click', function (event) {
  if (event.target.id === 'popup-add') {
    handleAddButtonClick(event);
  }
});

chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
  if (message.action === "showOverlay") {
    if (message.data) {
      console.log("togglePopup", message.data);

      let selection = window.getSelection();
      let range = selection.getRangeAt(0);
      let rect = range.getBoundingClientRect();

      let top = window.scrollY + rect.top + rect.height + "px";
      let left = window.scrollX + rect.left + rect.width + "px";

      togglePopup(top, left, message.data.word.word, "definition of " + message.data.definitions[0].definition + " ");
    }
  }
});

function ShowDefinitionOverlay(top, left, definition) {
  const overlayContainerId = "definition-overlay-container";
  // Check if the popup is already open
  const existingPopup = document.getElementById(overlayContainerId);
  if (existingPopup) {
    // If open, close the popup and return
    document.body.removeChild(existingPopup);
    return;
  }

  const overlayTemplate = `
  <div id="${overlayContainerId}" style="${getDefinitionOverlayStyles(
    top,
    left
  )}">
    <p>${definition}</p>
  </div>
`;

  // Append the popup container to the body
  document.body.insertAdjacentHTML("beforeend", overlayTemplate);

  // Show the popup
  const popupContainer = document.getElementById(overlayContainerId);
  popupContainer.style.display = "block";
}

function HideDefinitionOverlay() {
  const popupContainer = document.getElementById(
    "definition-overlay-container"
  );
  if (popupContainer) {
    document.body.removeChild(popupContainer);
  }
}

function addWord(word) {
  chrome.runtime.sendMessage(
    { action: 'addWord', word: word },
    function (res) {
      console.log('globalObj Response from background.js:', res);
      console.log(JSON.stringify(res))
      Object.assign(globalObj, res)
      highlightWordsInPageWithTooltip()
    }
  );
}

function handleAddButtonClick(event) {
  var word = event.target.value;
  console.log('Add button clicked for word:', word);

  addWord(word)
}

function togglePopup(top, left, word, definition) {
  const popupContainerId = "popup-container";

  // Check if the popup is already open
  const existingPopup = document.getElementById(popupContainerId);
  if (existingPopup) {
    // If open, close the popup and return
    document.body.removeChild(existingPopup);
    return;
  }

  // Create the popup container template
  const popupTemplate = `
    <div id="${popupContainerId}" style="${getPopupStyles(top, left)}">
      <h1>${word}</p>
      <p style="${getPopupContentStyles()}">${definition}</p>
      <button id="popup-add" value="${word}" style="${getAddButtonStyles()}" style="${getAddButtonStyles()}">Add</button>
    </div>
  `;

  // Append the popup container to the body
  document.body.insertAdjacentHTML("beforeend", popupTemplate);

  // Show the popup
  const popupContainer = document.getElementById(popupContainerId);
  popupContainer.style.display = "block";

  // Attach the close event to the close button
  const closeButton = document.getElementById("popup-add");
  closeButton.onclick = function () {
    hidePopup();
  };
}

function hidePopup() {
  const popupContainer = document.getElementById("popup-container");
  if (popupContainer) {
    document.body.removeChild(popupContainer);
  }
}

function getPopupStyles(top, left) {
  return `
    display: none;
    position: absolute;
    top: ${top};
    left: ${left};
    max-width: 400px;
    width: 100%;
    padding: 20px;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    z-index: 1000;
    text-align: center;
  `;
}

function getDefinitionOverlayStyles(top, left) {
  return `
    display: none;
    position: absolute;
    top: ${top};
    left: ${left};
    max-width: 400px;
    width: 100%;
    padding: 20px;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    z-index: 1000;
    text-align: center;
  `;
}

function getPopupContentStyles() {
  // Return the styles as a template string
  return "margin-bottom: 15px;";
}

function getAddButtonStyles() {
  // Return the styles as a template string
  return `
    background-color: #4CAF50;
    color: white;
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  `;
}
