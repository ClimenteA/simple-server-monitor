const settingsForm = document.getElementById("settings-form")
const urlElem = document.querySelector('[name="url"]')
const requestIntervalElem = document.querySelector('[name="requestInterval"]')
const apiKeyElem = document.querySelector('[name="apiKey"]')
const nukeAllElem = document.getElementById("clear-database")
const settingsOpenElem = document.getElementById("open-settings")
const settingsCloseElem = document.getElementById("close-settings")
const settingsModalElem = document.getElementById("settings-modal")
const settingsContainer = document.getElementById("settings-container")
const eventsContainer = document.getElementById("events-container")
const clearEventsElem = document.getElementById("clear-events")
const pauseNotificationsElem = document.getElementById("pause-notifications")


async function init() {
    setAlarmState()
    await setEvents()
    createEventsTable()
    createSettingsTable()
}

document.addEventListener("DOMContentLoaded", init)


chrome.storage.onChanged.addListener(async function (changes, areaName) {
    if (changes.events && areaName == "local") {
        createEventsTable()
    }
})

function setAlarmState(newState = null) {

    chrome.storage.sync.get(['alarmsPaused'], async function (items) {

        if (newState != null) {
            items.alarmsPaused = items.alarmsPaused ? false : true
            chrome.storage.sync.set({ 'alarmsPaused': items.alarmsPaused })
        }

        if (items.alarmsPaused == true) {
            pauseNotificationsElem.innerText = "🔕 Notifications paused"
        } else if (items.alarmsPaused == false) {
            pauseNotificationsElem.innerText = "🔔 Pause notifications"
        }

    })

}

pauseNotificationsElem.addEventListener("click", async function (event) {
    event.preventDefault()
    setAlarmState(true)
})


clearEventsElem.addEventListener("click", async function (event) {
    event.preventDefault()
    chrome.storage.local.set({ 'events': null })
    location.reload()
})

nukeAllElem.addEventListener("click", async function (event) {
    event.preventDefault()

    nukeAllElem.innerText = "💥 Boom.."

    await chrome.alarms.clearAll()
    chrome.storage.sync.clear()
    chrome.storage.local.clear()

    chrome.storage.sync.get(['settings'], async function (items) {
        if (!items.settings) return

        for (const data of Object.values(items.settings)) {

            if (data.apiKey.length == 0) continue

            try {

                const response = await fetch(data.url + "/simple-server-monitor/clear-database", {
                    method: "DELETE",
                    headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
                })

                if (response.status != 200) {
                    console.log(response)
                    alert(`Could not delete data for ${data.url}`)
                }

            } catch (error) {
                alert(`Could not delete data for ${error}`)
            }
        }
    })

    nukeAllElem.innerText = "💥 Done!"

    location.reload()
})

settingsOpenElem.addEventListener("click", () => {
    document.body.classList.add("modal-is-open")
    settingsModalElem.setAttribute("open", null)
})

settingsCloseElem.addEventListener("click", () => {
    document.body.classList.remove("modal-is-open")
    settingsModalElem.removeAttribute("open")
})


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

    function getServerDownEvent(data) {

        const now = new Date()
        const timestamp = now.toISOString().replace(/[-:T]/g, '').slice(0, 14)

        const receivedEvent = {
            EventId: "server-error-" + timestamp,
            Title: "Server down",
            Message: "Failed to fetch data from url: " + data.url,
            Level: "critical",
            Origin: data,
            Timestamp: timestamp
        }

        return receivedEvent
    }



    chrome.storage.sync.get(['settings'], function (settingsItems) {
        if (!settingsItems.settings) return

        chrome.storage.local.get(['events'], async function (items) {

            if (!items.events) items.events = {}

            let events = items.events
            for (const data of Object.values(settingsItems.settings)) {

                if (data.apiKey.length == 0) continue

                let receviedEvents

                try {

                    const response = await fetch(data.url + "/simple-server-monitor/notifications", {
                        method: "GET",
                        headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
                    })

                    if (response.status == 200) {
                        receviedEvents = await response.json()

                        for (const receivedEvent of receviedEvents.data || []) {
                            receivedEvent.Origin = data
                            events[receivedEvent.EventId] = receivedEvent
                        }
                    } else {
                        const receivedEvent = getServerDownEvent(data)
                        events[receivedEvent.EventId] = receivedEvent
                    }

                } catch (error) {
                    console.info("server is not responding: ", error)

                    const receivedEvent = getServerDownEvent(data)
                    events[receivedEvent.EventId] = receivedEvent

                }

            }

            const newEvents = {}
            for (const event of Object.values(events)) {
                newEvents[event.EventId] = event
            }

            console.log("Events in index.js", newEvents)

            chrome.storage.local.set({ 'events': newEvents })

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

    await chrome.alarms.clear(data.url)

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

function convertUtcToLocaleTimeString(utcIsoFormatString) {
    const utcDate = new Date(`${utcIsoFormatString.slice(0, 4)}-${utcIsoFormatString.slice(4, 6)}-${utcIsoFormatString.slice(6, 8)}T${utcIsoFormatString.slice(8, 10)}:${utcIsoFormatString.slice(10, 12)}:${utcIsoFormatString.slice(12, 14)}.000Z`)
    const options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' }
    const localTime = utcDate.toLocaleString(undefined, options)
    return localTime
}


const levelEmojiMapper = {
    "info": "📜",
    "warning": "⚠️",
    "critical": "🔥",
}

function appendEventsRow(data) {

    if (document.querySelector(`[data-event-id="${data.EventId}"]`)) {
        return
    }

    const row = document.createElement("tr")

    const levelCell = document.createElement("td")
    levelCell.innerText = levelEmojiMapper[data.Level] || "📜"
    row.appendChild(levelCell)

    const titleCell = document.createElement("td")
    titleCell.innerText = data.Title
    titleCell.setAttribute("data-event-id", data.EventId)
    row.appendChild(titleCell)

    const messageCell = document.createElement("td")
    messageCell.innerText = data.Message
    row.appendChild(messageCell)

    const urlCell = document.createElement("td")
    urlCell.innerText = data.Origin.url
    row.appendChild(urlCell)

    const timestampCell = document.createElement("td")
    timestampCell.innerText = convertUtcToLocaleTimeString(data.Timestamp)
    row.appendChild(timestampCell)

    const deleteCell = document.createElement("td")
    const deleteCellLink = document.createElement("a")
    deleteCellLink.innerText = "Delete"
    deleteCellLink.setAttribute("href", "#")
    deleteCell.appendChild(deleteCellLink)

    deleteCell.addEventListener("click", async () => {

        chrome.storage.local.get(['events'], async function (items) {
            if (!items.events) return
            delete items.events[data.EventId]
            chrome.storage.local.set({ 'events': items.events })
        })

        row.remove()
    })

    row.appendChild(deleteCell)

    eventsContainer.appendChild(row)

}


function appendSettingsRow(data) {

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

    deleteCell.addEventListener("click", async () => {
        chrome.storage.sync.get(['settings'], function (items) {
            if (!items.settings) return
            delete items.settings[data.url]
            chrome.storage.sync.set({ 'settings': items.settings })
        })

        await chrome.alarms.clear(data.url)

        row.remove()
    })

    settingsContainer.appendChild(row)

}
