package handlers

import (
	"encoding/json"
	"log"
	"strings"
)

func ParseResults(results []map[string]string) []ServerEvent {

	var serverEventResults []ServerEvent

	for _, result := range results {
		for key, value := range result {

			if strings.HasPrefix(key, "error-server-usage") {

				serverEvent := ServerEvent{
					Title:     "Watcher server error",
					Message:   value,
					Level:     "critical",
					Timestamp: strings.Split(key, "::")[1],
				}

				serverEventResults = append(serverEventResults, serverEvent)

			} else if strings.HasPrefix(key, "event") {

				var serverEvent ServerEvent
				err := json.Unmarshal([]byte(value), &serverEvent)
				if err != nil {
					log.Println("cannot unmarshal server event string")
					continue
				}

				serverEventResults = append(serverEventResults, serverEvent)

			}
		}
	}

	return serverEventResults
}
