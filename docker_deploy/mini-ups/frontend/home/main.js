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

