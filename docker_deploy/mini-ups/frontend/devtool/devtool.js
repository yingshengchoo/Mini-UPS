

function addItem(name = "", quantity = "") {
  const container = document.getElementById("itemsContainer");
  const row = document.createElement("div");
  row.className = "item-row";

  row.innerHTML = `
    <input placeholder="Item Name" value="${name}" />
    <input placeholder="Quantity" type="number" value="${quantity}" />
    <button type="button" class="remove-btn" onclick="this.parentElement.remove()">✖</button>
  `;

  container.appendChild(row);
}

function toggleCard(id) {
    const card = document.getElementById(id);
    const icon = card.querySelector('.toggle-icon');
    card.classList.toggle('expanded');
    const expanded = card.classList.contains('expanded');
    icon.textContent = expanded ? '▼ Collapse' : '▶ Expand';

    const body = card.querySelector('.card-body');
    if (body.style.display === 'none' || body.style.display === '') {
        body.style.display = 'flex';
    } else {
        body.style.display = 'none';
    }
}

const form = document.getElementById('packageForm');
const messageEl = document.getElementById('message');



form.addEventListener('submit', async (e) => {
  e.preventDefault();

  // get all items
  const itemRows = document.querySelectorAll("#itemsContainer .item-row");
  const items = Array.from(itemRows).map(row => {
    const inputs = row.querySelectorAll("input");
    return {
      name: inputs[0].value,
      quantity: parseInt(inputs[1].value)
    };
  });
  const data = {
    package_id: document.getElementById('package_id').value,
    username: document.getElementById('username').value,
    items: items,
    destination_x: parseInt(document.getElementById('destination_x').value),
    destination_y: parseInt(document.getElementById('destination_y').value),
    warehouse_id: parseInt(document.getElementById('warehouse_id').value),
  };
  // console.log(data)

  try {
    const res = await fetch('/api/package/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data),
    });

    const result = await res.json();
    if (res.ok) {
      messageEl.textContent = result.message;
      messageEl.className = 'message success';
    } else {
      messageEl.textContent = result.error || 'Something went wrong.';
      messageEl.className = 'message error';
    }
  } catch (err) {
    messageEl.textContent = 'Network error.';
    messageEl.className = 'message error';
  }
});