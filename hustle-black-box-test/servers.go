package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

var (
	servers        = map[string]*exec.Cmd{}
	serversTestEnv = map[string]string{
		"HUSTLE_HTTPADDR":     ":18661",
		"HUSTLE_HTTPPUBADDR":  "localhost:18661",
		"HUSTLE_HUBADDR":      ":16379",
		"HUSTLE_WSADDR":       ":18663",
		"HUSTLE_WSPUBADDR":    "localhost:18663",
		"HUSTLE_STATSADDR":    ":18665",
		"HUSTLE_STATSPUBADDR": "localhost:18665",
	}
	serversLock = &sync.Mutex{}
)

func setupServers() {
	var (
		job *exec.Cmd
		err error
	)

	setupServersTestEnv()

	serversLock.Lock()
	defer serversLock.Unlock()

	job = exec.Command("redis-server", "--port", "16379")
	err = job.Start()
	if err != nil {
		log.Panicln("failed to start redis-server job")
	}

	servers["redis-server"] = job

	job = exec.Command("hustle-server")
	err = job.Start()
	if err != nil {
		log.Panicln("failed to start hustle-server job")
	}

	servers["hustle-server"] = job
}

func setupServersTestEnv() {
	for key, value := range serversTestEnv {
		err := os.Setenv(key, value)
		if err != nil {
			log.Panicf("failed to set env var: %s\n", err)
		}
	}
}

func teardownServers() {
	serversLock.Lock()
	defer serversLock.Unlock()

	for name, job := range servers {
		if job == nil || job.Process == nil {
			delete(servers, name)
			continue
		}

		job.Process.Kill()
		delete(servers, name)
	}
}
