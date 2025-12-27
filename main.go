package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var BASH string = "/bin/bash"
var BASH_CMD string = "-c"
var CRONTAB_CMD string = "crontab"
var CRONTAB_LIST string = "-l"
var CRONTAB_FILE string = "/var/tmp/cronweb.txt"

//go:embed dist/*
var static embed.FS

type CronJob struct {
	Cron string
	Job  string
}

func shellGetCrontab() []CronJob {
	var cronjobs []CronJob
	out, _ := exec.Command(CRONTAB_CMD, CRONTAB_LIST).Output()
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '#' {
			split := strings.Split(line, " ")
			if len(split) > 5 {
				var cronjob CronJob
				cronjob.Cron = strings.Join(split[0:5][:], " ")
				cronjob.Job = strings.Join(split[5:][:], " ")
				cronjobs = append(cronjobs, cronjob)
			}
		}
	}
	return cronjobs
}

func shellSaveCrontab() {
	exec.Command(CRONTAB_CMD, CRONTAB_FILE).Run()
}

func shellAddCronJob(cronjob CronJob) {
	exec.Command(BASH, BASH_CMD, "echo \""+cronjob.Cron+" "+cronjob.Job+"\" >> "+CRONTAB_FILE).Run()
	shellSaveCrontab()
}

func shellClearCrontab() {
	exec.Command(BASH, BASH_CMD, "echo > "+CRONTAB_FILE).Run()
	shellSaveCrontab()
}

func shellWriteAllCronjobs(cronjobs []CronJob) {
	shellClearCrontab()
	for _, cronjob := range cronjobs {
		shellAddCronJob(cronjob)
	}
}

func initContext() {
	cronjobs := shellGetCrontab()
	shellWriteAllCronjobs(cronjobs)
}

func crontabHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(shellGetCrontab())
	case "PUT":
		if r.Body != nil {
			var cronjobs []CronJob
			decoder.Decode(&cronjobs)
			shellWriteAllCronjobs(cronjobs)
			json.NewEncoder(w).Encode(cronjobs)
		} else {
			http.Error(w, "body cannot be empty", http.StatusInternalServerError)
		}
	case "POST":
		if r.Body != nil {
			var cronjob CronJob
			decoder.Decode(&cronjob)
			shellAddCronJob(cronjob)
			json.NewEncoder(w).Encode(cronjob)
		} else {
			http.Error(w, "body cannot be empty", http.StatusInternalServerError)
		}
	}
}

func main() {
	initContext()
	contentStatic, _ := fs.Sub(static, "dist")
	http.Handle("/", http.FileServer(http.FS(contentStatic)))
	http.HandleFunc("/crontab", crontabHandler)

	defaultAddr := "0.0.0.0:8816"
	if len(os.Args) == 2 {
		defaultAddr = os.Args[1]
	}

	fmt.Printf("Starting cronweb on http://%s\n", defaultAddr)
	err := http.ListenAndServe(defaultAddr, nil)
	if err != nil {
		fmt.Println("Failed to start the server.")
		fmt.Print(err)
		return
	}
}
