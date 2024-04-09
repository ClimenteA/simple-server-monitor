

chrome.alarms.onAlarm.addListener(async (alarm) => {

    await chrome.notifications.create(alarm.name, {
        type: "basic",
        title: alarm.name,
        message: "You got some new important events. Checkout extension.",
        iconUrl: chrome.runtime.getURL("icons/notification.jpg")
    })

})

