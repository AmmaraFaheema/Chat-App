(function () {
  "use strict";

  console.log("app.js loader start");

  let ws = null;
  let myUsername = "";
  let currentChat = "";
  let messagesDiv = null;
  let messageInput = null;
  let sendBtn = null;
  let chatHistory = {};

  function saveToHistory(user, type, text) {
    if (!chatHistory[user]) chatHistory[user] = [];
    chatHistory[user].push({ type, text });
  }

  function loadHistory(user) {
    if (!messagesDiv) return;
    messagesDiv.innerHTML = "";
    if (!chatHistory[user]) return;

    chatHistory[user].forEach(msg => {
      addMessage(msg.type, msg.text);
    });
  }

  function addMessage(type, text) {
    if (!messagesDiv) return;
    const div = document.createElement("div");
    div.className = `message ${type}`;
    div.innerHTML = `${escapeHtml(text)}<div class="time">${new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</div>`;
    messagesDiv.appendChild(div);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }

  function escapeHtml(unsafe) {
    return unsafe
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  function connectWebSocket() {
    try {
      const nameParam = encodeURIComponent(myUsername || "");
      const proto = (location.protocol === "https:" ? "wss://" : "ws://");
      const url = proto + location.host + "/ws?name=" + nameParam;
      console.log("Connecting WS ->", url);
      ws = new WebSocket(url);

      ws.onmessage = (e) => {
        const parts = e.data.split(":", 2);
        const from = parts[0] || "";
        const msg = parts[1] || "";
        if (from && msg) {
          saveToHistory(from, "received", msg.trim());
        }

        if (from === currentChat) {
          addMessage("received", msg.trim());
        }
      };

      ws.onopen = () => console.log("WebSocket connected");
      ws.onclose = () => console.log("WebSocket disconnected");
      ws.onerror = (err) => console.error("WebSocket error:", err);
    } catch (err) {
      console.error("connectWebSocket exception:", err);
    }
  }

  function refreshUsers() {
    fetch("/users")
      .then(r => {
        if (!r.ok) throw new Error("HTTP " + r.status);
        return r.json();
      })
      .then(users => {
        const list = document.getElementById("userList");
        if (!list) return;
        list.innerHTML = "";
        Object.keys(users).forEach(name => {
          if (name === myUsername) return;
          const div = document.createElement("div");
          div.className = "user-item";
          div.innerHTML = `<div class="avatar">${escapeHtml(name[0] ? name[0].toUpperCase() : "?")}</div>
                           <div><b>${escapeHtml(name)}</b><br><small>${escapeHtml(users[name])}</small></div>`;
          div.onclick = (ev) => openChat(name, ev);
          list.appendChild(div);
        });
      })
      .catch(err => console.error("failed fetch users", err));
  }

  function openChat(username, ev) {
    currentChat = username;
    const cw = document.getElementById("chatWith");
    if (cw) cw.textContent = username;
    if (messageInput) messageInput.disabled = false;
    if (sendBtn) sendBtn.disabled = false;
    loadHistory(username);
    document.querySelectorAll(".user-item").forEach(u => u.classList.remove("active"));
    if (ev && ev.currentTarget) ev.currentTarget.classList.add("active");
  }

  function sendMessage() {
    if (!messageInput) return;
    const msg = messageInput.value.trim();
    if (!msg || !currentChat) return;

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(`${currentChat}:${msg}`);
      saveToHistory(currentChat, "sent", msg);

      addMessage("sent", msg);
      messageInput.value = "";
    } else {
      alert("WebSocket belum terhubung.");
    }
  }

  function login() {
    try {
      console.log("login() dipanggil");
      const inputEl = document.getElementById("usernameInput");
      if (!inputEl) return alert("Input username tidak ditemukan.");

      myUsername = inputEl.value.trim();
      if (!myUsername) return alert("Masukkan nama!");

      fetch("/register", { method: "POST", body: myUsername })
        .then(r => {
          if (!r.ok) throw new Error("register failed: " + r.status);

          const loginScreen = document.getElementById("loginScreen");
          const app = document.getElementById("app");
          if (loginScreen) loginScreen.classList.add("hidden");
          if (app) app.classList.remove("hidden");

          const myUserSpan = document.getElementById("myUsername");
          if (myUserSpan) myUserSpan.textContent = myUsername;

          connectWebSocket();
          setInterval(refreshUsers, 3000);
          refreshUsers();
        })
        .catch(err => {
          console.error("register error:", err);
          alert("Gagal register: " + err.message);
        });
    } catch (err) {
      console.error("login() exception:", err);
    }
  }

  function exposeGlobals() {
    window.login = login;
    window.sendMessage = sendMessage;
  }

  document.addEventListener("DOMContentLoaded", () => {
    console.log("DOM ready - initializing script.js");
    messagesDiv = document.getElementById("messages");
    messageInput = document.getElementById("messageInput");
    sendBtn = document.getElementById("sendBtn");

    const loginBtn = document.getElementById("loginBtn");
    if (loginBtn) loginBtn.addEventListener("click", login);

    if (sendBtn) sendBtn.addEventListener("click", sendMessage);

    if (messageInput) {
      messageInput.addEventListener("keypress", (e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          sendMessage();
        }
      });
    }

    exposeGlobals();
  });

})();
