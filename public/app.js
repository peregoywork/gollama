const ollamaURL = "http://localhost:11434/api/chat"
const modelName = "qwen2.5-coder:latest"

class ChatMessage {
    static VALID_ROLES = Object.freeze(["user", "assistant"])

    constructor(role, content) {
        this.role = role 
        this.content = content
    }
}

const messages = [];

const chatLog = document.getElementById("chat-log")
const prompt = document.getElementById("prompt")
const output = document.getElementById("output")


async function send() {

    console.log(prompt)

    const userMsg = prompt.value
    if (!userMsg) return;

    console.log("submitting prompt")

    appendMessage("user", userMsg)
    prompt.text = ""

    fetch(ollamaURL, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            model: modelName,
            messages: messages.map(m => ({ role: m.messageType, content: m.content })),
            stream: false // Note: set to false if your Go backend returns a single JSON object
        })
    })
    .then(res => {
        if (!res.ok) { 
            throw new Error("Server error: ${response.status}");
        }
        return res.json();
    })
    .then(data => {
        console.log(data.message)
        appendMessage(data.message.role, data.message.content);
    })
    .catch(error => {
        console.error("Fetch error:", error);
        appendMessage("assistant", "error generating response")
    })
}


function appendMessage(role, content) {
    const msg = new ChatMessage(role, content);
    messages.push(msg);

    const html = `<div class="message ${role}">${content}</div>`;
    chatLog.insertAdjacentHTML('beforeend', html);
}
