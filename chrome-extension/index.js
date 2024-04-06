const settingsForm = document.getElementById("settings-form")

settingsForm.addEventListener("submit", async function (event) {
    event.preventDefault()
    const formData = new FormData(event.target)
    const data = Object.fromEntries(formData.entries())

    console.log(data)
})




// chrome.storage.sync.set({ 'viewEvents': false })

// const toggleEventsBtn = document.getElementById("view-server-events")

// toggleEventsBtn.addEventListener("click", () => {
//     chrome.storage.sync.set({ 'viewEvents': true })
//     chrome.tabs.create({ url: "popup.html" })
// })


// chrome.storage.sync.get(['viewEvents'], function (items) {

//     const contentElem = document.getElementById("content")

//     if (!items.viewEvents) {
//         contentElem.style.display = "none"
//     } else {
//         contentElem.style.display = "block"

//         const settingsForm = document.getElementById("settings-form")

//         settingsForm.addEventListener("submit", async function (event) {
//             event.preventDefault()
//             const formData = new FormData(event.target)
//             const data = Object.fromEntries(formData.entries())

//             console.log(data)
//         })

//     }

// })


// // postUrl: null,
// // requestInterval: null,
// // apiKey: null,
