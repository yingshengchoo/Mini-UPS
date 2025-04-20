// track button
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


// logout button
async function logout() {

  try {

    const res = await fetch(`/api/user/logout`,{
      method:"POST"
    });
    if (!res.ok) throw new Error('Fail to logout');
    window.location.href = "/home/home.html";
  } catch (e) {
    resultEl.innerText = 'Failed to logout';
  }
}

// load user info when load this page
window.addEventListener('DOMContentLoaded', getUserInfo)
window.addEventListener('DOMContentLoaded', getPackageInfo)


// get user info
async function getUserInfo() {

  try{
    // get user info
    const response = await fetch('/api/user/info',{
      method:"GET",
      credentials:"include"
    })

    const data = await response.json()

    // check login status
    if (data.userlogined === false) {
      // show buttons
      const btn_login = document.getElementById('btn-login');
      const btn_register = document.getElementById('btn-register');
      btn_login.style.display = "block"
      btn_register.style.display = "block"

      // hide login status
      document.getElementById("login-status").style.display = "none"
    }else{
      // logined, show info
      const status = document.getElementById("login-status")
      status.innerText = "Logined as: "+data.username;
      status.style.display = 'block'
      document.getElementById("btn-logout").style.display = "block"

      // hide buttons
      document.getElementById('btn-login').style.display = "none";
      document.getElementById('btn-register').style.display = "none";
    }

  }catch(error) {
    console.error("failed to request", error);
  }
}

// get package info
async function getPackageInfo() {
  try{
    // TODO get package info
    const response = await fetch('/api/package/user/' + 1,{  // 有沒有可以拿現在User ID 的 function?
      method:"GET",
      credentials:"include"
    })

    if (!response.ok) {
      throw new Error(`Server error: ${response.status}`);
    }

    const packages = await response.json();
    sessionStorage.setItem("packages", JSON.stringify(packages));


    // // use fake data now
    // sessionStorage.setItem("packages", fakeData);

    const container = document.getElementById('user-packages');
    const template = document.getElementById('package-template');
    const title = document.getElementById('pacakge-title');

    // show title
    title.style.display = "block";
    container.style.display = "block";

    // refresh data
    container.querySelectorAll('.package-item:not(#package-template)').forEach(e => e.remove());
    packages = fakeData
    packages.forEach(pkg => {
      const clone = template.cloneNode(true);
      clone.id = "";
      clone.style.display = "block";

      clone.querySelector('.package-id').textContent = pkg.id;
      clone.querySelector('.package-contents').textContent = pkg.content;
      clone.querySelector('.package-address').textContent = pkg.address;
      clone.querySelector('.package-status').textContent = pkg.status;
      clone.querySelector('.package-location').textContent = pkg.location;
      clone.querySelector('.package-updatedAt').textContent = new Date(pkg.updatedAt).toLocaleString();

      const progressBar = clone.querySelector(".fancy-progress-bar");
      highlightProgressBar(progressBar, pkg.status);
      progressBar.style.display = 'block'

      // clone.appendChild(progressBar)
      container.appendChild(clone);
    });
  }
  catch(e){
    console.error("failed to request", e);
  }
}


const statusOrder = [
  "created", "packed", "picked", "loaded", "delivering", "delivered"
];

function highlightProgressBar(container, currentStatus) {
  const steps = container.querySelectorAll('.step');
  const currentIndex = statusOrder.indexOf(currentStatus.toLowerCase());

  steps.forEach((step, index) => {
    if (index <= currentIndex) {
      step.classList.add("completed");
    } else {
      step.classList.remove("completed");
    }
  });
}
const fakeData = [
  {
    id: 1,
    name: "Package 1",
    details: "Fragile items",
    content: "Electronics",
    address: "123 Street, City, Country",
    status: "delivering",
    location: "Warehouse A",
    updatedAt: "2025-04-19T12:30:00"
  },
  {
    id: 2,
    name: "Package 2",
    details: "Documents",
    content: "Paperwork",
    address: "456 Avenue, City, Country",
    status: "packed",
    location: "Warehouse B",
    updatedAt: "2025-04-19T14:00:00"
  }
];

// sleep for test
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}