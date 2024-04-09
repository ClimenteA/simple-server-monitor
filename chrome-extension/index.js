const settingsForm = document.getElementById("settings-form")
const urlElem = document.querySelector('[name="url"]')
const requestIntervalElem = document.querySelector('[name="requestInterval"]')
const apiKeyElem = document.querySelector('[name="apiKey"]')
const clearDBElem = document.getElementById("clear-database")
const settingsOpenElem = document.getElementById("open-settings")
const settingsCloseElem = document.getElementById("close-settings")
const settingsModalElem = document.getElementById("settings-modal")

settingsOpenElem.addEventListener("click", () => settingsModalElem.setAttribute("open", null))
settingsCloseElem.addEventListener("click", () => settingsModalElem.removeAttribute("open"))


document.addEventListener("DOMContentLoaded", () => {
    chrome.storage.sync.get(['settings'], function (items) {
        if (!items.settings) return
        for (const data of items.settings) appendRow(data)
    })
})


clearDBElem.addEventListener("click", async function (event) {
    event.preventDefault()
    // TODO
})


settingsForm.addEventListener("submit", async function (event) {
    event.preventDefault()
    const formData = new FormData(event.target)
    const data = Object.fromEntries(formData.entries())

    console.log(data)

    chrome.storage.sync.get(['settings'], function (items) {
        if (!items.settings) items.settings = []
        items.settings = [...items.settings, data]
        chrome.storage.sync.set({ 'settings': items.settings })
    })

    clearForm()
    appendRow(data)

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

    deleteCell.addEventListener("click", () => {
        chrome.storage.sync.get(['settings'], function (items) {
            if (!items.settings) return
            items.settings = items.settings.filter(item => item.apiKey != data.apiKey)
            chrome.storage.sync.set({ 'settings': items.settings })
        })
        row.remove()
    })

    settingsContainer.appendChild(row)

}
