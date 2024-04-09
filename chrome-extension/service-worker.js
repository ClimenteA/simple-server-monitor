

chrome.alarms.onAlarm.addListener(async (alarm) => {

    await chrome.notifications.create(alarm.name, {
        type: "basic",
        title: alarm.name,
        message: "You got some new important events.",
        iconUrl: chrome.runtime.getURL("static/notification.jpg")
    })

    await chrome.offscreen.createDocument({
        url: chrome.runtime.getURL('static/audio.html'),
        reasons: ['AUDIO_PLAYBACK'],
        justification: 'notification',
    })

})

