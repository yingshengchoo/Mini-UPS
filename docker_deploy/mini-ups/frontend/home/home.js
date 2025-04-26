// track button
async function track(packageID) {

  var trackingNumber = document.getElementById('trackingNumber').value;
  if (packageID != null){
    trackingNumber = packageID
  }
  const resultEl = document.getElementById('result');
  resultEl.style.display = 'block';
  resultEl.innerText = 'tracking...';

  try {
    const res = await fetch(`/api/package/info/${trackingNumber}`);
    if (!res.ok) {
      document.getElementById('track-package').style.display = "none";
      throw new Error('Fail to track');
    }

    const data = await res.json();

    console.log("tracking response:", data);


    const container = document.getElementById('track-package');
    const template = document.getElementById('package-template');
    const title = document.getElementById('package-title');

    // show title
    title.style.display = "block";
    container.style.display = "block";

    // refresh data
    container.querySelectorAll('.package-item:not(#package-template)').forEach(e => e.remove());

    const clone = template.cloneNode(true);
    const btn = clone.querySelector(".redirect-btn");
    const btn2 = clone.querySelector(".prioritize-btn");
    if (btn){
      btn.remove()
      btn2.remove()
    }
    clone.id = "";
    clone.style.display = "block";

    pkg = data
    clone.querySelector('.package-id').textContent = pkg.package_id;
    clone.querySelector('.package-contents').textContent = `${formatItems(pkg.items)}`;
    clone.querySelector('.package-address').textContent = `(${pkg.coord.x}, ${pkg.coord.y})`;
    clone.querySelector('.package-status').textContent = pkg.status;
    clone.querySelector('.package-warehouse').textContent = pkg.warehouse_id; //maybe not display?
    clone.querySelector('.package-updatedAt').textContent = `${formatDate(pkg.updated_at)}`;
    clone.querySelector('.package-priority').textContent = pkg.is_prioritized;


    const progressBar = clone.querySelector(".fancy-progress-bar");
    highlightProgressBar(progressBar, pkg.status);
    progressBar.style.display = 'block'

    // show this ele
    container.appendChild(clone);

    // hide message bar
    resultEl.style.display = 'none';
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
window.addEventListener('DOMContentLoaded', init);

async function init() {
  await getUserInfo();
  await getPackageInfo();
}

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
      sessionStorage.setItem("username", null);
      const btn_login = document.getElementById('btn-login');
      const btn_register = document.getElementById('btn-register');
      btn_login.style.display = "block"
      btn_register.style.display = "block"

      // hide login status
      document.getElementById("login-status").style.display = "none"
    }else{
      sessionStorage.setItem("username", data.username);

      // logined, show info
      const status = document.getElementById("login-status")
      status.innerText = "Logged in as: "+data.username;
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
    const userID = sessionStorage.getItem("username")
    //console.log("User info data:", data);
    const response = await fetch(`/api/package/user/${userID}`,{  // 有沒有可以拿現在User ID 的 function?
      method:"GET",
      credentials:"include"
    })

    if (!response.ok) {
      throw new Error(`Server error: ${response.status}`);
    }

    const packages = await response.json();
    sessionStorage.setItem("packages", JSON.stringify(packages));


    // // use fake data now
    //sessionStorage.setItem("packages", packages);

    const container = document.getElementById('user-packages');
    const template = document.getElementById('package-template');
    const title = document.getElementById('package-title');

    // show title
    title.style.display = "block";
    container.style.display = "block";

    // refresh data
    container.querySelectorAll('.package-item:not(#package-template)').forEach(e => e.remove());
    //packages = fakeData

    packages.forEach(pkg => {
      const clone = template.cloneNode(true);

      clone.id = "";
      clone.style.display = "block";
      const btn = clone.querySelector('.redirect-btn')
      const btn2 = clone.querySelector(".prioritize-btn");
      if (pkg.status === 'out_for_delivery' || pkg.status === 'delivered') {
        btn.disabled = true;
        btn2.disabled = true;
      }

      if (pkg.is_prioritized) {
        btn2.remove();
      }
      clone.querySelector('.package-id').textContent = pkg.package_id;
      clone.querySelector('.package-contents').textContent = `${formatItems(pkg.items)}`;
      clone.querySelector('.package-address').textContent = `(${pkg.coord.x}, ${pkg.coord.y})`;
      clone.querySelector('.package-status').textContent = pkg.status;
      clone.querySelector('.package-warehouse').textContent = pkg.warehouse_id; //maybe not display?
      clone.querySelector('.package-updatedAt').textContent = `${formatDate(pkg.updated_at)}`;
      clone.querySelector('.package-priority').textContent = pkg.is_prioritized;


      const progressBar = clone.querySelector(".fancy-progress-bar");
      highlightProgressBar(progressBar, pkg.status);
      progressBar.style.display = 'block'

      container.appendChild(clone);
    });
  }
  catch(e){
    console.error("failed to request", e);
  }
}

const dbStatusOrder = [
  "created", "packed", "pickup_complete", "loaded", "out_for_delivery", "delivered"
];


function highlightProgressBar(container, currentStatus) {
  const steps = container.querySelectorAll('.step');
  const currentIndex = dbStatusOrder.indexOf(currentStatus.toLowerCase());

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
  },
  {
    id: 2,
    name: "Package 3",
    details: "Funiture",
    content: "Paperwork",
    address: "456 Avenue, City, Country",
    status: "packed",
    location: "Warehouse B",
    updatedAt: "2025-04-19T14:00:00"
  }
];

//helper function for displaying items
const formatItems = (items) => {
  if (!Array.isArray(items)) {
    try {
      items = JSON.parse(items);
    } catch (e) {
      return String(items);
    }
  }
  return items.map(item => `${item.quantity} x ${item.name}`).join(', ');
};

const formatDate = (dateStr) => {
  const date = new Date(dateStr);
  return date.toLocaleString(); // e.g., "4/20/2025, 2:23:43 PM"
};

// sleep for test
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function copyLink(button) {
  // 获取当前按钮所在的 package-item
  const card = button.closest('.package-item');
  // 从 span 中获取 packageID
  const packageID = card.querySelector('.package-id').innerText.trim();

  if (!packageID) {
    // alert("未找到包裹 ID！");
    return;
  }

  // generate share link
  const response = await fetch(`/share/upshost`,{
    method:"GET",
    credentials:"include"
  })
  var data = await response.json()
  const shareUrl = `http://${data.upshost}:8080/share/${packageID}`;

  // copy link
  const textArea = document.createElement("textarea");
  textArea.value = shareUrl;
  document.body.appendChild(textArea);
  textArea.select();

  try {
    const successful = document.execCommand('copy');
    if (successful) {
      const status = card.querySelector('.copy-status');
      status.style.display = 'inline';
      setTimeout(() => {
        status.style.display = 'none';
      }, 2000);
    } else {
      alert("fail to copy");
    }
  } catch (err) {
    alert("fail to copy: " + err);
  }

  document.body.removeChild(textArea);
}


async function prioritizePackage(button) {
  // 获取当前按钮所在的 package-item
  const card = button.closest('.package-item');
  // 从 span 中获取 packageID
  const packageID = card.querySelector('.package-id').innerText.trim();
  if (!packageID) {
    // alert("未找到包裹 ID！");
    return;
  }

  const res = await fetch(`/api/package/prioritize/${packageID}`, {
    method: "POST",
    credentials: "include",
  });

  console.log("Removing button:", button);
  button.remove();
  if (!res.ok) throw new Error('Fail to prioritize package');
  location.reload();
}