package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ServerUsage struct {
	SERVER_CPU_MAX_USAGE        float64
	SERVER_RAM_MAX_USAGE        float64
	SERVER_DISK_MAX_USAGE       float64
	SERVER_USAGE_INTERVAL_CHECK int
	TIMESTAMP                   string
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
		SERVER_CPU_MAX_USAGE:        90,
		SERVER_RAM_MAX_USAGE:        90,
		SERVER_DISK_MAX_USAGE:       90,
		SERVER_USAGE_INTERVAL_CHECK: 60,
	}

	cpuMaxUsage, err := parseFloat(os.Getenv("SERVER_CPU_MAX_USAGE"))
	if err == nil {
		serverUsage.SERVER_CPU_MAX_USAGE = cpuMaxUsage
	}

	ramMaxUsage, err := parseFloat(os.Getenv("SERVER_RAM_MAX_USAGE"))
	if err == nil {
		serverUsage.SERVER_RAM_MAX_USAGE = ramMaxUsage
	}

	diskMaxUsage, err := parseFloat(os.Getenv("SERVER_DISK_MAX_USAGE"))
	if err == nil {
		serverUsage.SERVER_DISK_MAX_USAGE = diskMaxUsage
	}

	usageIntervalCheck, err := strconv.Atoi(os.Getenv("SERVER_USAGE_INTERVAL_CHECK"))
	if err == nil {
		serverUsage.SERVER_USAGE_INTERVAL_CHECK = usageIntervalCheck
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
			fmt.Println("Error executing the command:", err)
			return
		}

		ramUsage, err := getRamUsage()
		if err != nil {
			fmt.Println("Error executing the command:", err)
			return
		}

		diskUsage, err := getDiskUsage()
		if err != nil {
			fmt.Println("Error executing the command:", err)
			return
		}

		fmt.Printf("CPU usage: %.3f%%\n", cpuUsage)
		fmt.Printf("Ram usage: %.3f%%\n", ramUsage)
		fmt.Printf("Disk usage: %.3f%%\n", diskUsage)

		currentServerUsage := ServerUsage{
			SERVER_CPU_MAX_USAGE:  cpuUsage,
			SERVER_RAM_MAX_USAGE:  ramUsage,
			SERVER_DISK_MAX_USAGE: diskUsage,
			TIMESTAMP:             utcIsoNow,
		}

		currentServerUsageJSON, err := json.Marshal(currentServerUsage)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		Set("usage::"+utcIsoNow, string(currentServerUsageJSON))

		time.Sleep(time.Duration(serverUsage.SERVER_USAGE_INTERVAL_CHECK) * time.Second)
	}

}
