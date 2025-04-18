async function login() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    const data = await res.json();
    const result = document.getElementById("login-result");

    if (res.ok) {
        result.innerText = "Login successfully!";
        window.location.href = "/static/index.html";
    } else {
        result.innerText = "Failed to login: " + data.error;
    }
}