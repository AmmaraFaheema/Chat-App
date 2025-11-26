(function () {
  "use strict";

  console.log("app.js loader start");

  let ws = null;
  let myUsername = "";
  let currentChat = "";
  let messagesDiv = null;
  let messageInput = null;
  let sendBtn = null;

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
        if (from === currentChat || from === myUsername || from === "server") {
          addMessage(from === myUsername ? "sent" : "received", msg ? msg.trim() : "");
        }
      };

      ws.onopen = () => {
        console.log("WebSocket connected");
      };
      ws.onclose = () => {
        console.log("WebSocket disconnected");
      };
      ws.onerror = (err) => {
        console.error("WebSocket error:", err);
      };
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
          div.innerHTML = `<div class="avatar">${escapeHtml(name[0] ? name[0].toUpperCase() : "?")}</div><div><b>${escapeHtml(name)}</b><br><small>${escapeHtml(users[name])}</small></div>`;
          div.onclick = (ev) => openChat(name, ev);
          list.appendChild(div);
        });
      })
      .catch(err => {
        console.error("failed fetch users", err);
      });
  }

  function openChat(username, ev) {
    currentChat = username;
    const cw = document.getElementById("chatWith");
    if (cw) cw.textContent = username;
    if (messageInput) messageInput.disabled = false;
    if (sendBtn) sendBtn.disabled = false;
    if (messagesDiv) messagesDiv.innerHTML = "";
    document.querySelectorAll(".user-item").forEach(u => u.classList.remove("active"));
    if (ev && ev.currentTarget) ev.currentTarget.classList.add("active");
  }

  function sendMessage() {
    if (!messageInput) return;
    const msg = messageInput.value.trim();
    if (!msg || !currentChat) return;
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(`${currentChat}:${msg}`);
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
      if (!inputEl) {
        alert("Input username tidak ditemukan pada halaman.");
        return;
      }
      myUsername = inputEl.value.trim();
      if (!myUsername) {
        console.log("username kosong");
        return alert("Masukkan nama!");
      }

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
    window.addPeerToGroup = addPeerToGroup;
  }

const chatHistory = {}; 

function appendMessage(sender, text, target) {
  const key = target || currentChatTarget;
  if (!chatHistory[key]) chatHistory[key] = [];

  const div = document.createElement("div");
  div.className = sender === myName ? "message sent" : "message received";
  div.innerHTML = `<strong>${sender}:</strong> ${text}`;
  chatHistory[key].push(div.outerHTML);

  if (currentChatTarget === key) {
    document.getElementById("messages").appendChild(div);
    div.scrollIntoView();
  }
}

function startChat(target, type) {
  currentChatTarget = target;
  const displayName = type === "group" ? "Group: " + target : target;
  document.getElementById("chatWith").textContent = displayName;
  
  const messagesDiv = document.getElementById("messages");
  messagesDiv.innerHTML = "";

  if (chatHistory[target]) {
    chatHistory[target].forEach(html => {
      messagesDiv.innerHTML += html;
    });
    messagesDiv.lastElementChild?.scrollIntoView();
  }

  document.getElementById("messageInput").disabled = false;
  document.getElementById("sendBtn").disabled = false;
  document.getElementById("messageInput").focus();
}

function sendMessage() {
  const input = document.getElementById("messageInput");
  const msg = input.value.trim();
  if (!msg || !currentChatTarget) return;

  const payload = {
    to: currentChatTarget,
    message: msg,
    type: currentChatTarget === myGroup || 
           ["umum","kelompok-a","kelompok-b","kelompok-c","rahasia"].includes(currentChatTarget)
           ? "group" : "private",
    from: myName
  };

  fetch("/send", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  appendMessage(myName, msg, currentChatTarget);
  input.value = "";
}

  document.addEventListener("DOMContentLoaded", () => {
    console.log("DOM ready - initializing script.js");
    messagesDiv = document.getElementById("messages");
    messageInput = document.getElementById("messageInput");
    sendBtn = document.getElementById("sendBtn");

    const loginBtn = document.getElementById("loginBtn");
    if (loginBtn) {
      loginBtn.addEventListener("click", login);
    }

    if (sendBtn) {
      sendBtn.addEventListener("click", sendMessage);
    }

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
