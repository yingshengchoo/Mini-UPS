async function trackPackage() {
    const trackingNumber = document.getElementById('trackingInput').value;
    const resultDiv = document.getElementById('result');
  
    if (!trackingNumber) {
      resultDiv.innerHTML = 'Please enter package ID';
      return;
    }
  
    resultDiv.innerHTML = 'Tracing...';
  
    try {
      const response = await fetch(`/api/track/${trackingNumber}`);
      if (!response.ok) throw new Error('Fail to trace');
  
      const data = await response.json();
      resultDiv.innerHTML = `
        <p><strong>status:</strong>${data.status}</p>
        <p><strong>current palce:</strong>${data.location}</p>
        <p><strong>update time:</strong>${data.updated_at}</p>
      `;
    } catch (error) {
      resultDiv.innerHTML = 'Fail to trace, please check your package ID';
    }
  }