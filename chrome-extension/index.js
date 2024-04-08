const settingsForm = document.getElementById("settings-form")
const urlElem = document.querySelector('[name="url"]')
const requestIntervalElem = document.querySelector('[name="requestInterval"]')
const apiKeyElem = document.querySelector('[name="apiKey"]')
const clearDBElem = document.getElementById("clear-database")


clearDBElem.addEventListener("click", async function (event) {
    event.preventDefault()
    // TODO
})


settingsForm.addEventListener("submit", async function (event) {
    event.preventDefault()
    const formData = new FormData(event.target)
    const data = Object.fromEntries(formData.entries())

    console.log(data)

    const response = await fetch(data.url, {
        method: "POST",
        body: JSON.stringify(data)
    })

    if (response.status == 200) {

        const jsonResponse = await response.json()
        console.log(jsonResponse)

        clearForm()
        appendRow(data)

    } else {
        console.error(response)
    }

})


function clearForm() {
    urlElem.value = null
    requestIntervalElem.value = null
    apiKeyElem.value = null
}


function appendRow(data) {

    const settingsContainer = document.getElementById("settings-container")

    let count = 0
    for (const key in data) if (data[key]) count += 1
    if (Object.values(data).length != count) return

    const row = document.createElement("tr")
    for (const key in data) {
        const cell = document.createElement("td")
        cell.innerText = data[key]
        row.appendChild(cell)
    }

    const editCell = document.createElement("td")
    const deleteCell = document.createElement("td")
    const editCellLink = document.createElement("a")
    const deleteCellLink = document.createElement("a")

    editCellLink.innerText = "Edit"
    editCellLink.setAttribute("href", "#")
    editCellLink.setAttribute("data-edit", data.apiKey)

    deleteCellLink.innerText = "Delete"
    deleteCellLink.setAttribute("href", "#")
    deleteCellLink.setAttribute("data-delete", data.apiKey)

    editCell.appendChild(editCellLink)
    deleteCell.appendChild(deleteCellLink)

    row.setAttribute("id", data.apiKey)

    row.appendChild(editCell)
    row.appendChild(deleteCell)

    editCell.addEventListener("click", () => {
        urlElem.value = data.url
        requestIntervalElem.value = data.requestInterval
        apiKeyElem.value = data.apiKey
        row.remove()
    })

    deleteCell.addEventListener("click", () => row.remove())

    settingsContainer.appendChild(row)

}



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


// // url: null,
// // requestInterval: null,
// // apiKey: null,
