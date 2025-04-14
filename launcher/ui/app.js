const launchButton = document.querySelector("#launchButton");

var ui = {
  version: "",
  setVersion(version) {
    this.version = version;
    document.querySelector("#version").textContent = version;
  },
};

document.addEventListener("DOMContentLoaded", () => {
  launchButton.addEventListener("click", () => window.launchMinecraft());
  window.initLauncher().then((res) => {
    ui.setVersion(res);
  });
});
