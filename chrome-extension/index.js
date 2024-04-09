const settingsForm = document.getElementById("settings-form")
const urlElem = document.querySelector('[name="url"]')
const requestIntervalElem = document.querySelector('[name="requestInterval"]')
const apiKeyElem = document.querySelector('[name="apiKey"]')
const clearDBElem = document.getElementById("clear-database")
const settingsOpenElem = document.getElementById("open-settings")
const settingsCloseElem = document.getElementById("close-settings")
const settingsModalElem = document.getElementById("settings-modal")
const settingsContainer = document.getElementById("settings-container")
const eventsContainer = document.getElementById("events-container")


async function init() {
    await setEvents()
    createEventsTable()
    createSettingsTable()
}

document.addEventListener("DOMContentLoaded", init)


clearDBElem.addEventListener("click", async function (event) {
    event.preventDefault()

    chrome.storage.sync.set({ 'settings': null })
    chrome.storage.local.set({ 'events': null })

    chrome.storage.sync.get(['settings'], async function (items) {
        if (!items.settings) return

        for (const data of Object.values(items.settings)) {

            if (data.apiKey.length == 0) continue

            const response = await fetch(data.url + "clear-database", {
                method: "DELETE",
                headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
            })

            if (response.status != 200) {
                console.log(response)
                alert(`Could not delete data for ${data.url}`)
            }
        }
    })

    location.reload()
})

settingsOpenElem.addEventListener("click", () => settingsModalElem.setAttribute("open", null))
settingsCloseElem.addEventListener("click", () => settingsModalElem.removeAttribute("open"))


function createSettingsTable() {
    chrome.storage.sync.get(['settings'], function (items) {
        if (!items.settings) return
        for (const data of Object.values(items.settings)) appendSettingsRow(data)
    })
}


function createEventsTable() {
    chrome.storage.local.get(['events'], function (items) {
        if (!items.events) return
        for (const data of Object.values(items.events)) appendEventsRow(data)
    })
}

async function setEvents() {

    chrome.storage.sync.get(['settings'], function (settingsItems) {
        if (!settingsItems.settings) return

        chrome.storage.local.get(['events'], async function (items) {

            if (!items.events) items.events = {}

            let events = items.events
            for (const data of Object.values(settingsItems.settings)) {

                if (data.apiKey.length == 0) continue

                console.log("Get events with settings", data)

                const response = await fetch(data.url, {
                    method: "GET",
                    headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
                })

                if (response.status == 200) {

                    const receviedEvents = await response.json()

                    for (const receivedEvent of receviedEvents.data || []) {
                        receivedEvent.Origin = data
                        events[receivedEvent.EventId] = receivedEvent
                    }
                } else {
                    console.log(data)
                    alert("Could not get events")
                }

            }

            chrome.storage.local.set({ 'events': events })

        })

    })

}


settingsForm.addEventListener("submit", async function (event) {
    event.preventDefault()
    const formData = new FormData(event.target)
    const data = Object.fromEntries(formData.entries())


    if (!(data.url.startsWith("http://") || data.url.startsWith("https://"))) {
        alert("Not a valid url")
        return
    }

    if (data.apiKey.length == 0) {
        alert("Not a valid apiKey")
        return
    }

    chrome.storage.sync.get(['settings'], function (items) {
        if (!items.settings) items.settings = {}
        items.settings[data.url] = data
        chrome.storage.sync.set({ 'settings': items.settings })
    })

    clearForm()
    appendSettingsRow(data)
    await setAlarm(data)

})


async function setAlarm(data) {

    await chrome.alarms.create(data.url, {
        delayInMinutes: 1,
        periodInMinutes: Number(data.requestInterval)
    })

}

function clearForm() {
    urlElem.value = null
    requestIntervalElem.value = null
    apiKeyElem.value = null
}


function appendEventsRow(data) {

    const row = document.createElement("tr")

    const levelCell = document.createElement("td")
    levelCell.innerText = data.Level
    row.appendChild(levelCell)

    const titleCell = document.createElement("td")
    titleCell.innerText = data.Title
    row.appendChild(titleCell)

    const messageCell = document.createElement("td")
    messageCell.innerText = data.Message
    row.appendChild(messageCell)

    const timestampCell = document.createElement("td")
    timestampCell.innerText = data.Timestamp
    row.appendChild(timestampCell)

    const urlCell = document.createElement("td")
    urlCell.innerText = data.Origin.url
    row.appendChild(urlCell)

    const deleteCell = document.createElement("td")
    const deleteCellLink = document.createElement("a")
    deleteCellLink.innerText = "Delete"
    deleteCellLink.setAttribute("href", "#")
    deleteCell.appendChild(deleteCellLink)

    deleteCell.addEventListener("click", async () => {

        console.log("Deleting:", data)

        const response = await fetch(data.Origin.url + `delete/${data.EventId}`, {
            method: "DELETE",
            headers: { "Content-Type": "application/json", "ApiKey": data.Origin.apiKey }
        })

        if (response.status == 200) {

            chrome.storage.local.get(['events'], async function (items) {
                if (!items.events) return
                delete items.events[data.EventId]
                chrome.storage.local.set({ 'events': items.events })
            })

            row.remove()

        } else {
            console.log(data)
            alert("Could not delete row")
        }

    })

    row.appendChild(deleteCell)

    eventsContainer.appendChild(row)

}


function appendSettingsRow(data) {

    console.log("Settings row:", data)

    let count = 0
    for (const key in data) if (data[key]) count += 1
    if (Object.values(data).length != count) return

    const row = document.createElement("tr")

    const urlCell = document.createElement("td")
    urlCell.innerText = data.url
    row.appendChild(urlCell)

    const reqInterCell = document.createElement("td")
    reqInterCell.innerText = data.requestInterval
    row.appendChild(reqInterCell)

    const apiKeyCell = document.createElement("td")
    apiKeyCell.innerText = data.apiKey
    row.appendChild(apiKeyCell)


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
            delete items.settings[data.url]
            chrome.storage.sync.set({ 'settings': items.settings })
        })
        row.remove()
    })

    settingsContainer.appendChild(row)

}
