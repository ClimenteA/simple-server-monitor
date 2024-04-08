package handlers

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ServerUsage struct {
	CPU_MAX_USAGE        float64
	RAM_MAX_USAGE        float64
	DISK_MAX_USAGE       float64
	USAGE_INTERVAL_CHECK int
	TIMESTAMP            string
	ERROR                string
}

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

		log.Printf("CPU usage: %.3f%%\n", cpuUsage)
		log.Printf("Ram usage: %.3f%%\n", ramUsage)
		log.Printf("Disk usage: %.3f%%\n", diskUsage)

		if cpuUsage >= serverUsage.CPU_MAX_USAGE || ramUsage >= serverUsage.RAM_MAX_USAGE || diskUsage >= serverUsage.DISK_MAX_USAGE {

			currentServerUsage := ServerUsage{
				CPU_MAX_USAGE:  cpuUsage,
				RAM_MAX_USAGE:  ramUsage,
				DISK_MAX_USAGE: diskUsage,
				TIMESTAMP:      utcIsoNow,
			}

			currentServerUsageJSON, err := json.Marshal(currentServerUsage)
			if err != nil {
				log.Println("Error:", err)
				Set("error-server-usage-marshal::"+utcIsoNow, "failed to convert struct to json")
				return
			}

			Set("server-usage::"+utcIsoNow, string(currentServerUsageJSON))

		}

		time.Sleep(time.Duration(serverUsage.USAGE_INTERVAL_CHECK) * time.Second)
	}

}

func ParseServerUsageResults(results []map[string]string) []ServerUsage {

	usageIntervalCheck, err := strconv.Atoi(os.Getenv("USAGE_INTERVAL_CHECK"))
	if err != nil {
		panic("cannot get USAGE_INTERVAL_CHECK from envs")
	}

	var serverUsageResults []ServerUsage
	for _, result := range results {
		for key, value := range result {

			if strings.HasPrefix(key, "server-usage") {

				var serverUsage ServerUsage
				err := json.Unmarshal([]byte(value), &serverUsage)
				if err != nil {
					log.Println("cannot unmarshal server usage string")
					continue
				}
				serverUsage.USAGE_INTERVAL_CHECK = usageIntervalCheck
				serverUsageResults = append(serverUsageResults, serverUsage)

			} else if strings.HasPrefix(key, "error-server-usage") {

				serverUsage := ServerUsage{
					USAGE_INTERVAL_CHECK: usageIntervalCheck,
					TIMESTAMP:            strings.Split(key, "::")[0],
					ERROR:                value,
				}

				serverUsageResults = append(serverUsageResults, serverUsage)

			}

		}
	}

	return serverUsageResults
}
