const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");

// Loads and renders image from file or url
function renderImage() {
  const imageUrlInput = document.getElementById("imageUrlInput");
  if (imageUrlInput.value) {
    return renderImageFromURL();
  } else {
    return renderImageFromFile();
  }
}

// Loads and renders image from url
function renderImageFromURL() {
  let imageUrlInput = document.getElementById("imageUrlInput");
  let img = new Image();
  img.src = imageUrlInput.value;

  return new Promise((resolve, reject) => {
    img.onload = function () {
      render(img)
      resolve();
    };
    img.onerror = reject;
  });
}

// Sets new url for image and renders it
function setImageUrl(url) {
  document.getElementById("imageUrlInput").value = url;
  renderImageFromURL();
  hideMessages();
}

// Clears the canvas
function cleanCanvas() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
}

// Loads and renders image from file
function renderImageFromFile() {
  const imageFileInput = document.getElementById("imageFileInput");
  let file = imageFileInput.files[0]; // Get loaded file
  if (!file) {
    return false; // If file is not chosen, return false
  }

  // Read file as Data URL
  let reader = new FileReader();
  reader.readAsDataURL(file);

  // Clear canvas before drawing
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  return new Promise((resolve, reject) => {
    reader.onload = function () {
      let img = new Image();
      img.src = reader.result;
      img.onload = function () {
        render(img)
        resolve();
      };
      img.onerror = reject;
    };
  });
}

function render(img) {
  // Set canvas dimensions according to image dimensions
  canvas.width = img.width;
  canvas.height = img.height;

  // Calculate aspect ratio of the image
  const aspectRatio = img.width / img.height;

  // Calculate image dimensions to fit the canvas
  let drawWidth = canvas.width;
  let drawHeight = canvas.height;

  if (aspectRatio > 1) {
    // If image is horizontal, scale width
    drawHeight = canvas.width / aspectRatio;
  } else {
    // If image is vertical, scale height
    drawWidth = canvas.height * aspectRatio;
  }

  // Draw image with preserved aspect ratio
  ctx.drawImage(img, 0, 0, drawWidth, drawHeight);
}

// Draw coordinates
async function drawCoordinates(coordinates) {
  // Render image again
  await renderImage();

  // Check if coordinates are available
  if (coordinates.length < 2) return; // At least two points are needed to draw a polygon

  // Set up for drawing a polygon
  ctx.strokeStyle = "#00ff00"; // Border color of polygon
  ctx.fillStyle = "rgba(0, 255, 0, 0.3)"; // Fill color of polygon with transparency
  ctx.lineWidth = 2; // Line width of polygon outline

  // Start drawing a polygon
  ctx.beginPath();
  ctx.moveTo(coordinates[0][0], coordinates[0][1]);

  // Connect lines between all points
  coordinates.forEach(([x, y]) => {
    ctx.lineTo(x, y);
  });

  // Draw polygon
  ctx.closePath();
  ctx.fill();
  ctx.stroke();

  // if displayValues is not checked - return
  if (!document.getElementById("displayValues").checked) return;

  ctx.font = "12px Arial"; // Font size/family
  ctx.fillStyle = "white"; // Font color
  const radius = 4; // Dots radius

  coordinates.forEach(([x, y], index) => {
    ctx.fillStyle = "#00ff00"; // Dots color

    // Draw dot
    ctx.beginPath();
    ctx.arc(x, y, radius, 0, 2 * Math.PI);
    ctx.fill();

    const text = `${index + 1}: (${x}, ${y})`;
    console.log(text);

    // Draw black background for text
    const textWidth = ctx.measureText(text).width;
    const padding = 2;

    ctx.globalAlpha = 0.7;
    ctx.fillStyle = "black"; // BG color
    ctx.fillRect(
      x + 10 - padding,
      y - 20 - padding,
      textWidth + padding * 2,
      16 + padding * 2
    );
    ctx.globalAlpha = 1.0;

    // Draw white text over the background
    ctx.fillStyle = "white";
    ctx.fillText(text, x + 10, y - 10);
  });

}

// Update inputs if "Is JSON" button is clicked
function updateRequestIsJson() {
  if (document.getElementById("requestIsJson").checked) {
    document.getElementById("form-request-group").style.display = "none";
    document.getElementById("json-request-group").style.display = "block";
  } else {
    document.getElementById("form-request-group").style.display = "block";
    document.getElementById("json-request-group").style.display = "none";
  }
}

// Update inputs on selecting predefined service
function setPredefinedService() {
  const service = document.getElementById("predefined-service");

  // If user defined own settings – skip
  if (!service.value) {
    return;
  }

  const endpointInput = document.getElementById("endpointInput");
  endpointInput.value = service.value;

  let selected = service.options[service.selectedIndex];
  var requestJSONTemplate = selected.getAttribute('data-json-template');
  var requestFormdataField = selected.getAttribute('data-formdata-field');
  var responseXYField = selected.getAttribute('data-xy-field');

  document.getElementById("requestJsonTemplate").value = requestJSONTemplate;
  document.getElementById("requestFormdataField").value = requestFormdataField;
  document.getElementById("responseXYField").value = responseXYField;
  document.getElementById("requestIsJson").checked = (requestJSONTemplate != "");
  updateRequestIsJson();
}

// Form submission
function sendForm() {
  if (
    !document.getElementById("imageFileInput").files[0]
    && !document.getElementById("imageUrlInput").value
  ) {
    return showAlert("Select image first");;
  }

  let form = document.getElementById("form");
  fetch(form.action, {
    method: "post",
    body: new FormData(form),
  })
    .then((response) => {
      console.log("Response status: " + response.status);
      if (response.status === 200) {
        return response.json();
      } else {
        response.json().then((data) => {
          // showAlert("Service error (HTTP " + response.status + "): " + data.error);
          showResponse(data, response.status);
        });
      }
    })
    .then((data) => {
      if (data) {
        drawCoordinates(data.coordinates);
        // showSuccess("Coordinates received (" + data.coordinates.length + "). Took " + data.executionTime + "ms");
        showResponse(data, 200);
      }
    });
}

function showResponse(data, status) {
  document.getElementById("response-div").style.display = "block";
  document.getElementById("response-timestamp").innerHTML = new Date().toLocaleString();

  let message = '';
  if (status === 200) {
    message = "Coordinates received. Took " + data.executionTime + "ms";
    document.getElementById("response-status-ok").innerHTML = "HTTP 200";
    document.getElementById("response-status-ok").style.display = "inline-block";
    document.getElementById("response-status-error").style.display = "none";
  } else {
    message = "Backend error: " + data.error;
    document.getElementById("response-status-error").innerHTML = "HTTP " + status;
    document.getElementById("response-status-error").style.display = "inline-block";
    document.getElementById("response-status-ok").style.display = "none";
  }
  document.getElementById("response-message").innerHTML = message

  var pretty_response = beautify(data.response, null, 2, 100);
  console.log(pretty_response);
  document.getElementById("response-json").innerHTML = pretty_response;
}

function hideResponse() {
  document.getElementById("response-div").style.display = "none";
}

// Show alert message
function showAlert(message) {
  hideMessages()
  document.getElementById("alert-msg").style.display = "block";
  document.getElementById("alert-msg").innerHTML = message;
}

// Hide all messages
function hideMessages() {
  document.getElementById("alert-msg").style.display = "none";
}

// Reset service selector to use user-defined config
function setUserDefinedConfig() {
  document.getElementById("predefined-service").value = "";
}

// When predefined service is changed
document.getElementById("predefined-service").addEventListener("change", function () {
  setPredefinedService();
  hideMessages();
  cleanCanvas();
  hideResponse();
  renderImage();
});

// When any request input is changed –> set user defined config
document.getElementById("endpointInput").addEventListener("change", function () {
  setUserDefinedConfig();
});
document.getElementById("requestJsonTemplate").addEventListener("change", function () {
  setUserDefinedConfig();
});
document.getElementById("requestFormdataField").addEventListener("change", function () {
  setUserDefinedConfig();
});
document.getElementById("responseXYField").addEventListener("change", function () {
  setUserDefinedConfig();
});

// When Is JSON checkbox change
document.getElementById("requestIsJson").addEventListener("change", function () {
  updateRequestIsJson();
  setUserDefinedConfig();
});

// When image changed
document.getElementById("imageUrlInput").addEventListener("change", function () {
  renderImageFromURL();
  hideMessages();
});

// When file uploaded
document.getElementById("imageFileInput").addEventListener("change", function () {
  renderImageFromFile();
  hideMessages();
});

// When page is loaded
document.addEventListener("DOMContentLoaded", function () {
  setPredefinedService();
});

// Mouse coordinates
const coordinatesDisplay = document.getElementById("coordinates");
canvas.addEventListener("mousemove", function (event) {
  let rect = canvas.getBoundingClientRect();
  let x = Math.floor(event.clientX - rect.left);
  let y = Math.floor(event.clientY - rect.top);
  coordinatesDisplay.textContent = `(x: ${x}, y: ${y})`;
});

canvas.addEventListener("mouseleave", function () {
  coordinatesDisplay.textContent = `Cursor coordinates`;
});
