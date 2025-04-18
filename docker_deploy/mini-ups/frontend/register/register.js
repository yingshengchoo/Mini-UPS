async function register() {
    const username = document.getElementById("reg-username").value;
    const password = document.getElementById("reg-password").value;

    const res = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    const data = await res.json();
    const result = document.getElementById("register-result");

    if (res.ok) {
        result.style.display = 'block';
        if (data.error != null) {
            result.innerText = data.error;
        }else{
            result.innerText = data.message+"\nredirecting to login...";
            setTimeout(() => {
                window.location.href = "/login/login.html";
            }, 1000);
        }
    } else {
        result.innerText = "Failed to register: " + data.error;
    }
}