document.addEventListener("DOMContentLoaded", () => {
  const container = document.getElementById("user-packages");

  container.addEventListener("click", function (e) {
    if (e.target.classList.contains("redirect-btn")) {
      const card = e.target.closest(".package-item");

      const packageId = card.querySelector(".package-id").textContent.trim();
      const x = prompt("Enter X coordinate:");
      const y = prompt("Enter Y coordinate:");

      if (!x || !y || isNaN(x) || isNaN(y)) {
        alert("Invalid coordinates");
        return;
      }

      const payload = {
        package_id: packageId,
        x: parseInt(x),
        y: parseInt(y)
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
        return res.json();
      })
      .then(data => {
        alert("Location updated successfully!");
      })
      .catch(err => {
        console.error(err);
        alert("Failed to update location.");
      });
    }
  });
});