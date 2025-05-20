package elog

import "log"

func Info(data string) {
	log.Printf("[INFO] %s", data)
}

func Fatal(data error) {
	log.Fatalf("[ERROR] %s", data)
}
