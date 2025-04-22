window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    var lastSentMessage;
  
    function appendLog(item) {
      var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
      log.appendChild(item);
      if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
      }
    }
  
    function formatTime() {
      var now = new Date();
      return now.getHours().toString().padStart(2, '0') + ":" + now.getMinutes().toString().padStart(2, '0');
    }
  
    function addMessage(username, message) {
      var messageElement = document.createElement("div");
      messageElement.classList.add("message");
  
      
      if (username === "user1") {
        messageElement.classList.add("user1");
      } else {
        messageElement.classList.add("user2");
      }
  
      var usernameElement = document.createElement("div");
      usernameElement.classList.add("username");
      usernameElement.innerText = username;
  
      var timeElement = document.createElement("div");
      timeElement.classList.add("time");
      timeElement.innerText = formatTime();
  
      var textElement = document.createElement("div");
      textElement.classList.add("text");
      textElement.innerText = message;
  
      messageElement.appendChild(usernameElement);
      messageElement.appendChild(timeElement);
      messageElement.appendChild(textElement);
  
      appendLog(messageElement);
    }
  
    document.getElementById("form").onsubmit = function () {
      if (!conn) {
        return false;
      }
      if (!msg.value) {
        return false;
      }
      
      const message = msg.value;
      conn.send(message);
  
      addMessage("user1", message);
  
      lastSentMessage = message;
  
      msg.value = "";
      return false;
    };
  
    if (window["WebSocket"]) {
      conn = new WebSocket("ws://" + document.location.host + "/ws");
      conn.onclose = function (evt) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        appendLog(item);
      };
      conn.onmessage = function (evt) {
        const messages = evt.data.split('\n');
        for (let i = 0; i < messages.length; i++) {
          if (messages[i] !== lastSentMessage) {
            addMessage("user2", messages[i]);
          }
        }
      };
    } else {
      var item = document.createElement("div");
      item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
      appendLog(item);
    }
  };

  // переключатель
  const toggle = document.getElementById('theme-switch');

  toggle.addEventListener('change', () => {
    document.body.classList.toggle('light-theme');
  });
  