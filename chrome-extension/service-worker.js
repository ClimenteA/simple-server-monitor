

chrome.alarms.onAlarm.addListener(async function (alarm) {

    chrome.storage.sync.get(['alarmsPaused'], async function (items) {

        if (items.alarmsPaused == true) return

        chrome.storage.sync.get(['settings'], async function (items) {
            if (!items.settings) return

            const data = items.settings[alarm.name]

            let receviedEvents

            try {
                const response = await fetch(data.url, {
                    method: "GET",
                    headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
                })

                if (response.status != 200) return
                receviedEvents = await response.json()
                if (!receviedEvents.data) return

            } catch (error) {
                console.error(error)

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

                let events = items.events
                for (const receivedEvent of receviedEvents.data) {
                    receivedEvent.Origin = data
                    events[receivedEvent.EventId] = receivedEvent
                }

                await chrome.offscreen.createDocument({
                    url: chrome.runtime.getURL('static/audio.html'),
                    reasons: ['AUDIO_PLAYBACK'],
                    justification: 'notification',
                })

                await chrome.notifications.create(alarm.name, {
                    type: "basic",
                    title: alarm.name,
                    message: "You got some new important events.",
                    iconUrl: chrome.runtime.getURL("static/notification.jpg")
                })

                chrome.storage.local.set({ 'events': events })

            })
        })
    })

})

