package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	healthUrl := os.Getenv("SIMPLE_SERVER_MONITOR_HEALTH_URL")
	if strings.HasPrefix(healthUrl, "http") {
		serverUsage.HEALTH_URL = healthUrl
	} else {
		log.Fatal("please provide SIMPLE_SERVER_MONITOR_HEALTH_URL in .env file")
	}

	cpuMaxUsage, err := parseFloat(os.Getenv("SIMPLE_SERVER_MONITOR_CPU_MAX_USAGE"))
	if err == nil {
		serverUsage.CPU_MAX_USAGE = cpuMaxUsage
	} else {
		log.Fatal("SIMPLE_SERVER_MONITOR_CPU_MAX_USAGE missing")
	}

	ramMaxUsage, err := parseFloat(os.Getenv("SIMPLE_SERVER_MONITOR_RAM_MAX_USAGE"))
	if err == nil {
		serverUsage.RAM_MAX_USAGE = ramMaxUsage
	} else {
		log.Fatal("SIMPLE_SERVER_MONITOR_RAM_MAX_USAGE missing")
	}

	diskMaxUsage, err := parseFloat(os.Getenv("SIMPLE_SERVER_MONITOR_DISK_MAX_USAGE"))
	if err == nil {
		serverUsage.DISK_MAX_USAGE = diskMaxUsage
	} else {
		log.Fatal("SIMPLE_SERVER_MONITOR_DISK_MAX_USAGE missing")
	}

	usageIntervalCheck, err := strconv.Atoi(os.Getenv("SIMPLE_SERVER_MONITOR_USAGE_INTERVAL_CHECK"))
	if err == nil {
		serverUsage.USAGE_INTERVAL_CHECK = usageIntervalCheck
	} else {
		log.Fatal("SIMPLE_SERVER_MONITOR_USAGE_INTERVAL_CHECK missing")
	}

	return serverUsage
}

func getCpuUsage() (float64, error) {
	// cat /proc/stat |grep cpu |tail -1|awk '{print ($5*100)/($2+$3+$4+$5+$6+$7+$8+$9+$10)}'|awk '{print 100-$1"%"}'
	usagePercent, err := exec.Command("bash", "-c", "cat /proc/stat |grep cpu |tail -1|awk '{print ($5*100)/($2+$3+$4+$5+$6+$7+$8+$9+$10)}'|awk '{print 100-$1\"%\"}'").Output()
	if err != nil {
		return 0, err
	}
	return parseFloat(string(usagePercent))
}

func getRamUsage() (float64, error) {
	// free -h | awk '/^Mem:/ {print ($3/$2)*100"%"}'
	usagePercent, err := exec.Command("bash", "-c", "free -h | awk '/^Mem:/ {print ($3/$2)*100\"%\"}'").Output()
	if err != nil {
		return 0, err
	}
	return parseFloat(string(usagePercent))
}

func getDiskUsage() (float64, error) {
	// df -h | awk '$6 == "/" {print $5}'
	usagePercent, err := exec.Command("bash", "-c", "df -h | awk '$6 == \"/\" {print $5}'").Output()
	if err != nil {
		return 0, err
	}
	return parseFloat(string(usagePercent))
}

func getIsHealthyResponse(healthUrl string) error {

	log.Println("Checking ", healthUrl)
	resp, err := http.Get(healthUrl)
	if err != nil {
		log.Println("Got an error ", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Status code not 200!")
		return fmt.Errorf("%s responded with a %d status code instead of 200", healthUrl, resp.StatusCode)
	}

	log.Println("Looks good ", healthUrl, resp)

	return nil

}

func sleep(waitValue int) {
	time.Sleep(time.Duration(waitValue) * time.Second)
}

func MonitorServer() {
	serverUsage := getServerUsage()

	log.Println("got server usage values: ", serverUsage)

	for {

		now := time.Now().UTC()
		utcIsoNow := now.Format("20060102150405")

		err := getIsHealthyResponse(serverUsage.HEALTH_URL)
		if err != nil {
			Set("error-server-health-url::"+utcIsoNow, err.Error())
			log.Println(err.Error())
			sleep(serverUsage.USAGE_INTERVAL_CHECK)
			continue
		}

		cpuUsage, err := getCpuUsage()
		if err != nil {
			log.Println("Error executing the command:", err)
			Set("error-server-usage-cpu::"+utcIsoNow, "failed to get cpu usage")
			log.Println(err.Error())
			sleep(serverUsage.USAGE_INTERVAL_CHECK)
			continue
		}

		ramUsage, err := getRamUsage()
		if err != nil {
			log.Println("Error executing the command:", err)
			Set("error-server-usage-ram::"+utcIsoNow, "failed to get ram usage")
			log.Println(err.Error())
			sleep(serverUsage.USAGE_INTERVAL_CHECK)
			continue
		}

		diskUsage, err := getDiskUsage()
		if err != nil {
			log.Println("Error executing the command:", err)
			Set("error-server-usage-disk::"+utcIsoNow, "failed to get disk usage")
			log.Println(err.Error())
			sleep(serverUsage.USAGE_INTERVAL_CHECK)
			continue
		}

		if cpuUsage >= serverUsage.CPU_MAX_USAGE || ramUsage >= serverUsage.RAM_MAX_USAGE || diskUsage >= serverUsage.DISK_MAX_USAGE {

			eventId := uuid.NewString()
			message := fmt.Sprintf("Server resources have reached critical levels. CPU: %.3f%%, RAM: %.3f%%, DISK: %.3f%%", cpuUsage, ramUsage, diskUsage)

			log.Println(message)

			event := ServerEvent{
				EventId:   eventId,
				Title:     "Server resources",
				Message:   message,
				Level:     "warning",
				Timestamp: utcIsoNow,
			}

			currentEventJSON, err := json.Marshal(event)
			if err != nil {
				log.Println(err.Error())
				Set("error-event-marshal::"+utcIsoNow, "failed to convert struct to json on MonitorServerUsage")
				continue
			}

			Set("event::"+eventId, string(currentEventJSON))

		}

		log.Printf("server %s looks ok", serverUsage.HEALTH_URL)
		sleep(serverUsage.USAGE_INTERVAL_CHECK)
	}

}
