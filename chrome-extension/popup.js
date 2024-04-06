document.getElementById("view-server-events").addEventListener("click", () => {
    chrome.tabs.create({ url: "index.html" })
})
