document.addEventListener("DOMContentLoaded", () => {
  const container = document.getElementById("user-packages");

  container.addEventListener("click", function (e) {
    if (e.target.classList.contains("redirect-btn")) {
      const card = e.target.closest(".package-item");

      const packageId = card.querySelector(".package-id").textContent.trim();
    //   const x = prompt("Enter X coordinate:");
    //   const y = prompt("Enter Y coordinate:");
    showFloatingModal(packageId)
    console.log(coords)

    }
  });
});


function submitRedirect() {
    const x = parseInt(document.getElementById("x-coord").value);
    const y = parseInt(document.getElementById("y-coord").value);
    const id = document.getElementById("modal-package-id").textContent;
  
    console.log(`Package ${id} redirect to (${x}, ${y})`);
  
    // send POST 
    if (!x || !y || isNaN(x) || isNaN(y)) {
        alert("Invalid coordinates");
        return;
      }

      const payload = {
        package_id: id,
        "coordinate":{
            x: parseInt(x),
            y: parseInt(y)
        }
      };

      fetch("/api/package/redirect", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      })
      .then(res => {
        if (!res.ok) throw new Error("Failed to send");
        init()
        return res.json();
      })
      .then(data => {
        alert("Location updated successfully!");
      })
      .catch(err => {
        console.error(err);
        alert("Failed to update location.");
      });

    // close modal
    document.querySelector("#custom-modal").remove();
    document.querySelector(".modal-backdrop").remove();
    document.body.style.overflow = "";
  }

async function showFloatingModal(packageId) {
    const backdrop = document.createElement("div");
    backdrop.className = "modal-backdrop";
    document.body.appendChild(backdrop);
  
    // load redirect.html
    const res = await fetch("redirect.html");
    const modalHTML = await res.text();
  
    // create DOM and insert
    const temp = document.createElement("div");
    temp.innerHTML = modalHTML;
    const modal = temp.firstElementChild;
    document.body.appendChild(modal);
  
    // show id
    modal.querySelector("#modal-package-id").textContent = packageId;
  
    document.body.style.overflow = "hidden";
  
    // button
    backdrop.onclick = () => {
      modal.remove();
      backdrop.remove();
      document.body.style.overflow = "";
    };
  }