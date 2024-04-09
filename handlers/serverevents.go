package handlers

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func SaveEvent(event ServerEvent) error {

	now := time.Now().UTC()
	utcIsoNow := now.Format("20060102150405")
	eventId := uuid.NewString()

	event.EventId = eventId
	event.Timestamp = utcIsoNow

	currentEventJSON, err := json.Marshal(event)
	if err != nil {
		Set("error-event-marshal::"+utcIsoNow, "failed to convert struct to json on SaveEvent timestamp:"+utcIsoNow)
		return err
	}

	return Set("event::"+eventId, string(currentEventJSON))
}

func DeleteEvent(eventId string) error {
	return Del("event::" + eventId)
}
