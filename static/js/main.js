const form = document.getElementById("form");
const error = document.getElementById("error");
const progress = document.getElementById("progress");
const images = document.getElementById("images");

htmx.on("#form", "htmx:xhr:progress", function (evt) {
  progress.setAttribute("value", (evt.detail.loaded / evt.detail.total) * 100);
});

form.addEventListener("htmx:responseError", function (event) {
  error.innerText = event.detail.xhr.responseText;
});

form.addEventListener("htmx:afterOnLoad", function () {
  this.reset();
  progress.setAttribute("value", 0);
});

images.addEventListener("htmx:afterSwap", function () {
  error.innerText = "";
});
