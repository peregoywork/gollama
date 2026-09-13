const ollamaURL = "http://localhost:11434/api/chat"

const messages = document.getElementById("chat-log")
const prompt = document.getElementById("prompt")
const output = document.getElementById("output")


async function send() {
    output.textContent = "Generating..."

    try {
        const res = await fetch(ollamaURL, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({
                model: "qwen2.5-coder:latest",
                messages: [
                    { role: "user", content: prompt.text }
                ],
                stream: false // Note: set to false if your Go backend returns a single JSON object
            })
        })

        const data = await res.json()
        output.textContent = data.message.content
    } catch (err) {
        output.textContent = "Error: " + err.message
    }
}

function appendMessage(message) {
    // is message ai or user
    // use "pre" for both
}
