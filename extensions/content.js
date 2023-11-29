window.addEventListener("load", function () {
  highlight();
});

window.addEventListener("scroll", function () {
  highlight();
});

function highlight() {
  const wordsToHighlight = ["Petersen", "represent", "word"]; // Add your list of words here

  function highlightWords(words) {
    const regex = new RegExp(`\\b(${words.join("|")})\\b`, "gi");

    const allElements = document.querySelectorAll(
      "*:not(script):not(style):not(textarea)"
    );

    allElements.forEach((element) => {
      if (
        element.childNodes.length === 1 &&
        element.childNodes[0].nodeType === 3 &&
        element.className != "highlighted-word"
      ) {
        element.innerHTML = element.innerHTML.replace(
          regex,
          (match) =>
            `<span class="highlighted-word" onmouseover="{}" style="background-color: yellow;" class="highlight">${match}</span>`
        );

        const highlightedSpans = element.querySelectorAll(".highlighted-word");
        highlightedSpans.forEach(function (span) {
          span.addEventListener("mouseover", function () {
            // Add your mouse hover logic here
            console.log(`Hovered over: ${span.innerText}`);
            rect = span.getBoundingClientRect();
            let top = window.scrollY + rect.top + rect.height + "px";
            let left = window.scrollX + rect.left + rect.width + "px";
            ShowDefinitionOverlay(top, left, "test difinition");
          });

          span.addEventListener("mouseout", function () {
            // Add your mouse out logic here
            console.log(`Mouse out: ${span.innerText}`);
            HideDefinitionOverlay();
          });
        });
      }
    });
  }

  highlightWords(wordsToHighlight);
}

document.addEventListener("click", function (event) {
  console.log("User clicked at coordinates:", {
    x: event.clientX,
    y: event.clientY,
  });
  hidePopup();
});

chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
  if (message.action === "showOverlay") {
    if (message.selectionText) {
      console.log("togglePopup", message.selectionText);

      let selection = window.getSelection();
      let range = selection.getRangeAt(0);
      let rect = range.getBoundingClientRect();

      let top = window.scrollY + rect.top + rect.height + "px";
      let left = window.scrollX + rect.left + rect.width + "px";

      word = message.selectionText;
      togglePopup(top, left, word, "definition of " + word + " ");
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
      <button id="popup-add" style="${getAddButtonStyles()}">Add</button>
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
