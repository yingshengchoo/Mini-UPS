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

    result.style.display = 'block'
    if (res.ok) {
        if (data.login == true){

        
            result.innerText = "Login successfully!\nRedirecting to home...";
            setTimeout(() => {
                window.location.href = "/home/home.html";
            },1000);
        }else{
            result.innerText = "wrong username or password";
        }
    } else {
        result.innerText = "Failed to login: " + data.error;
    }
}