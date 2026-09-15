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
    const userMsg = prompt.value
    if (!userMsg) return;

    appendMessage("user", userMsg)
    prompt.value = ""

    fetch(ollamaURL, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
            model: modelName,
            messages: messages.map(m => ({ role: m.role, content: m.content })),
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
    let tag = "";
    messages.push(msg);

    switch (msg.role) {
        case "user":
            tag = 'div';
            break;
        case "assistant":
            tag = 'pre';
            break;
        default:
            return;
    }
    const html = `<${tag} class="message ${role}">${content}</${tag}>`;
    chatLog.insertAdjacentHTML('beforeend', html);
}
