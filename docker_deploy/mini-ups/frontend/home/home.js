async function track() {
    const trackingNumber = document.getElementById('trackingNumber').value;
    const resultEl = document.getElementById('result');
    resultEl.style.display = 'block';
    resultEl.innerText = 'tracking...';
  
    try {
      const res = await fetch(`/api/track/${trackingNumber}`);
      if (!res.ok) throw new Error('Fail to track');
  
      const data = await res.json();
      resultEl.innerText = JSON.stringify(data, null, 2);
    } catch (e) {
      resultEl.innerText = 'Failed to track';
    }
  }

window.addEventListener('DOMContentLoaded', () => {
  fetch('/api/user/info',{
    method:"GET",
    credentials:"include"
  })
    .then(response => {
      if (response.userlogined === false) {
        // show buttons
        const btn_login = document.getElementById('btn-login');
        const btn_register = document.getElementById('btn-register');
        btn_login.style.display = "block"
        btn_register.style.display = "block"

        // hide login status
        document.getElementById("login-status").style.display = "none"
      }
      return response.json();
    })
    .then(data => {
      if (data.userlogined) {
        // logined, show info
        const status = document.getElementById("login-status")
        status.innerText = "Logined as: "+data.username;
        status.style.display = 'block'
        document.getElementById("btn-logout").style.display = "block"

        // hide buttons
        document.getElementById('btn-login').style.display = "none";
        document.getElementById('btn-register').style.display = "none";
      }
    })
    .catch(error => {
      console.error("failed to request", error);
    });
});