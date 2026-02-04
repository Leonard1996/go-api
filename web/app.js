const packSizesInput = document.getElementById("packSizes");
const saveButton = document.getElementById("saveSizes");
const amountInput = document.getElementById("amount");
const calcButton = document.getElementById("calculate");
const resultTable = document.getElementById("resultTable");
const resultBody = resultTable.querySelector("tbody");

const parseSizes = (value) => {
  const parts = value.split(/\r?\n/).map((part) => part.trim()).filter(Boolean);
  return parts.map((part) => Number(part));
};

const renderResults = (packs) => {
  resultBody.innerHTML = "";
  const entries = Object.entries(packs || {}).sort((a, b) => Number(b[0]) - Number(a[0]));
  if (entries.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = "<td colspan=\"2\">No packs needed.</td>";
    resultBody.appendChild(row);
    return;
  }
  for (const [size, count] of entries) {
    const row = document.createElement("tr");
    row.innerHTML = `<td>${size}</td><td>${count}</td>`;
    resultBody.appendChild(row);
  }
};

const loadPackSizes = async () => {
  try {
    const res = await fetch("/v1/pack-sizes");
    if (!res.ok) {
      throw new Error("failed");
    }
    const data = await res.json();
    const sizes = (data.pack_sizes || []).map((size) => String(size));
    packSizesInput.value = sizes.join("\n");
  } catch (err) {
    packSizesInput.value = "250\n500\n1000\n2000\n5000";
  }
};

saveButton.addEventListener("click", async () => {
  const sizes = parseSizes(packSizesInput.value);
  try {
    const res = await fetch("/v1/pack-sizes", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ pack_sizes: sizes }),
    });
    const data = await res.json();
    if (!res.ok) {
      alert(data.error || "Failed to update pack sizes");
      return;
    }
    packSizesInput.value = (data.pack_sizes || []).join("\n");
    alert("Pack sizes updated");
  } catch (err) {
    alert("Failed to update pack sizes");
  }
});

calcButton.addEventListener("click", async () => {
  const amount = Number(amountInput.value);
  try {
    const res = await fetch("/v1/calculate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ amount }),
    });
    const data = await res.json();
    if (!res.ok) {
      alert(data.error || "Calculation failed");
      return;
    }
    renderResults(data.packs || {});
  } catch (err) {
    alert("Calculation failed");
  }
});

loadPackSizes();
