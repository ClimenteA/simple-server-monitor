

chrome.alarms.onAlarm.addListener(async function (alarm) {

    chrome.storage.sync.get(['alarmsPaused'], async function (items) {

        if (items.alarmsPaused == true) return

        chrome.storage.sync.get(['settings'], async function (items) {
            if (!items.settings) return

            const data = items.settings[alarm.name]

            const response = await fetch(data.url, {
                method: "GET",
                headers: { "Content-Type": "application/json", "ApiKey": data.apiKey }
            })

            if (response.status != 200) return
            const receviedEvents = await response.json()
            if (!receviedEvents.data) return

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

