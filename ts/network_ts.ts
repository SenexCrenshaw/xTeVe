class Server {
  protocol: string
  cmd: string

  constructor(cmd: string) {
    this.cmd = cmd
  }

  request(data: Object): any {

    if (this.cmd != "updateLog") {
      showElement("loading", true)
      UNDO = new Object()
    }

    switch (window.location.protocol) {
      case "http:":
        this.protocol = "ws://"
        break
      case "https:":
        this.protocol = "wss://"
        break
    }

    var url = this.protocol + window.location.hostname + ":" + window.location.port + "/data/" + "?Token=" + getCookie("Token")

    data["cmd"] = this.cmd
    var ws = new WebSocket(url)
    ws.onopen = function () {

      WS_AVAILABLE = true
      this.send(JSON.stringify(data));

    }

    ws.onerror = function (wsErrEvt) {

      console.log("No websocket connection to xTeVe could be established. Check your network configuration.")

      if (WS_AVAILABLE == false) {
        alert("No websocket connection to xTeVe could be established. Check your network configuration.")
      }

    }


    ws.onmessage = function (wsMessageEvt) {

      showElement("loading", false)

      const response: Object = JSON.parse(wsMessageEvt.data);

      if (response.hasOwnProperty("token")) {
        document.cookie = "Token=" + response["token"];
      }

      if (response["status"] == false) {
        alert(response["err"]);
        return;
      }

      if (response.hasOwnProperty('openLink')) {
        window.location = response['openLink'];
      }

      if (response.hasOwnProperty("reload")) {
        window.location.reload();
      }

      if (response.hasOwnProperty("alert")) {
        alert(response["alert"]);
      }

      if (response.hasOwnProperty("logoURL")) {
        var div = (document.getElementById("channel-icon") as HTMLInputElement)
        div.value = response["logoURL"]
        div.className = "changed"
        return
      }

      if (response.hasOwnProperty("m3uURL")) {
        var div = (document.getElementById("m3u") as HTMLInputElement)
        div.value = response["m3uURL"]
        return
      }

      if (response.hasOwnProperty("xmlURL")) {
        var div = (document.getElementById("xmltv") as HTMLInputElement)
        div.value = response["xmlURL"]
        return
      }

      switch (data["cmd"]) {
        case "updateLog":
          SERVER["log"] = response["log"]
          if (document.getElementById("content_log")) {
            showLogs(false)
          }
          return

        case "saveSettingsByName":
          return
        case "showInfo":
          return
        case "showDebug":
          return
        case "showError":
          return
        case "showWarning":
          return

        default:
          SERVER = new Object()
          SERVER = response
          break;
      }

      if (response.hasOwnProperty("openMenu")) {
        var menu = document.getElementById(response["openMenu"])
        menu.click()
        showElement("popup", false)
      }

      if (response.hasOwnProperty("reload")) {
        location.reload()
      }

      if (response.hasOwnProperty("wizard")) {
        createLayout()
        configurationWizard[response["wizard"]].createWizard()
        return
      }

      createLayout()

    }

  }

}

function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) {
    return parts.pop().split(";").shift();
  }
}
