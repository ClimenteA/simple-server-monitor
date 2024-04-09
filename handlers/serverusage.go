package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func parseFloat(floatStr string) (float64, error) {
	cleanFloat := strings.TrimSuffix(strings.TrimSpace(floatStr), "%")
	value, err := strconv.ParseFloat(cleanFloat, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getServerUsage() ServerUsage {

	serverUsage := ServerUsage{
		CPU_MAX_USAGE:        90,
		RAM_MAX_USAGE:        90,
		DISK_MAX_USAGE:       90,
		USAGE_INTERVAL_CHECK: 60,
	}

	cpuMaxUsage, err := parseFloat(os.Getenv("CPU_MAX_USAGE"))
	if err == nil {
		serverUsage.CPU_MAX_USAGE = cpuMaxUsage
	}

	ramMaxUsage, err := parseFloat(os.Getenv("RAM_MAX_USAGE"))
	if err == nil {
		serverUsage.RAM_MAX_USAGE = ramMaxUsage
	}

	diskMaxUsage, err := parseFloat(os.Getenv("DISK_MAX_USAGE"))
	if err == nil {
		serverUsage.DISK_MAX_USAGE = diskMaxUsage
	}

	usageIntervalCheck, err := strconv.Atoi(os.Getenv("USAGE_INTERVAL_CHECK"))
	if err == nil {
		serverUsage.USAGE_INTERVAL_CHECK = usageIntervalCheck
	}

	return serverUsage
}

func getCpuUsage() (float64, error) {
	usagePercent, err := exec.Command("bash", "-c", "cat /proc/stat |grep cpu |tail -1|awk '{print ($5*100)/($2+$3+$4+$5+$6+$7+$8+$9+$10)}'|awk '{print 100-$1\"%\"}'").Output()
	if err != nil {
		return 0, err
	}
	return parseFloat(string(usagePercent))
}

func getRamUsage() (float64, error) {
	usagePercent, err := exec.Command("bash", "-c", "free -h | awk '/^Mem:/ {print ($3/$2)*100\"%\"}'").Output()
	if err != nil {
		return 0, err
	}
	return parseFloat(string(usagePercent))
}

func getDiskUsage() (float64, error) {
	usagePercent, err := exec.Command("bash", "-c", "df -h | awk '$6 == \"/\" {print $5}'").Output()
	if err != nil {
		return 0, err
	}
	return parseFloat(string(usagePercent))
}

func MonitorServerUsage() {
	serverUsage := getServerUsage()

	for {

		now := time.Now().UTC()
		utcIsoNow := now.Format("20060102150405")

		cpuUsage, err := getCpuUsage()
		if err != nil {
			log.Println("Error executing the command:", err)
			Set("error-server-usage-cpu::"+utcIsoNow, "failed to get cpu usage")
			return
		}

		ramUsage, err := getRamUsage()
		if err != nil {
			log.Println("Error executing the command:", err)
			Set("error-server-usage-ram::"+utcIsoNow, "failed to get ram usage")
			return
		}

		diskUsage, err := getDiskUsage()
		if err != nil {
			log.Println("Error executing the command:", err)
			Set("error-server-usage-disk::"+utcIsoNow, "failed to get disk usage")
			return
		}

		if cpuUsage >= serverUsage.CPU_MAX_USAGE || ramUsage >= serverUsage.RAM_MAX_USAGE || diskUsage >= serverUsage.DISK_MAX_USAGE {

			eventId := uuid.NewString()

			event := ServerEvent{
				EventId:   eventId,
				Title:     "Server resources",
				Message:   fmt.Sprintf("Server resources have reached critical levels. CPU: %.3f%%, RAM: %.3f%%, DISK: %.3f%%", cpuUsage, ramUsage, diskUsage),
				Level:     "warning",
				Timestamp: utcIsoNow,
			}

			currentEventJSON, err := json.Marshal(event)
			if err != nil {
				log.Println("Error:", err)
				Set("error-event-marshal::"+utcIsoNow, "failed to convert struct to json on MonitorServerUsage")
				return
			}

			Set("event::"+eventId, string(currentEventJSON))

		}

		time.Sleep(time.Duration(serverUsage.USAGE_INTERVAL_CHECK) * time.Second)
	}

}
