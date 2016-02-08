package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/theckman/uselessd/echo"
)

const appVersion = "0.0.1"

const rfc3339Msec = "2006-01-02T15:04:05.000Z07:00"

func main() {
	flags := &commandLine{}

	out, err := flags.Parse(nil)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if out != "" {
		fmt.Print(out)
		os.Exit(0)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: rfc3339Msec,
	})

	es := &uselessecho.Server{
		Host: flags.Host,
		Port: 10007,
	}

	if err := es.ListenAndServe(); err != nil {
		log.Fatalf("listen error: %s", err)
	}
}
