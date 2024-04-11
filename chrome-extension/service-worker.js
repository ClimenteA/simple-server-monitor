
async function showNotification(alarmName) {

    try {
        await chrome.offscreen.createDocument({
            url: chrome.runtime.getURL('static/audio.html'),
            reasons: ['AUDIO_PLAYBACK'],
            justification: 'notification',
        })
    } catch (error) {
        console.info("cannot play audio ", error)
    }

    let iconPath = "static/notification.jpg"
    let message = "You got some new server notifications!"
    if (alarmName == "Server down") {
        iconPath = "static/fire.jpg"
        message = "🔥🔥🔥 Server not responding! 🔥🔥🔥"
    }

    await chrome.notifications.create(alarmName, {
        type: "basic",
        title: alarmName,
        message: message,
        iconUrl: chrome.runtime.getURL(iconPath)
    })

}


chrome.storage.onChanged.addListener(async function (changes, areaName) {
    if (changes.events && areaName == "local") {
        const event = Object.values(changes.events?.newValue)[0]
        if (!event) return
        if (event.EventId.startsWith("server-error")) {
            await showNotification("Server down")
        } else {
            await showNotification(event.Origin.url)
        }
    }
})


chrome.alarms.onAlarm.addListener(async function (alarm) {

    chrome.storage.sync.get(['alarmsPaused'], async function (items) {

        if (items.alarmsPaused == true) return

        chrome.storage.sync.get(['settings'], async function (items) {
            if (!items.settings) return

            const data = items.settings[alarm.name]

            let receviedEvents

            try {
                const response = await fetch(data.url + "/notifications", {
                    method: "GET",
                    headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
                })

                if (response.status != 200) return
                receviedEvents = await response.json()
                if (!receviedEvents.data) return

            } catch (error) {
                console.info("cannot fetch notifications", error)

                const now = new Date()
                const timestamp = now.toISOString().replace(/[-:T]/g, '').slice(0, 14)

                receviedEvents = {
                    data: [{
                        EventId: "server-error-" + timestamp,
                        Title: "Server down",
                        Message: "Failed to fetch data from url: " + data.url,
                        Level: "critical",
                        Timestamp: timestamp
                    }]
                }
            }

            chrome.storage.local.get(['events'], async function (items) {

                if (!items.events) items.events = {}

                for (const receivedEvent of receviedEvents.data) {
                    receivedEvent.Origin = data
                }

                let events = [...Object.values(items.events), ...receviedEvents.data]

                const newEvents = {}
                for (const event of events) {
                    newEvents[event.EventId] = event
                }

                console.log("Events in service worker:", newEvents)

                chrome.storage.local.set({ 'events': newEvents })

            })
        })
    })

})

