// de-job-killer
//
// A tool for either killing a job that's currently running or to mark a job
// as Failed.
//
// This tool works by either sending out a stop request for a job or by sending
// out a job status update message that marks the job as failed.
package main

import (
	"flag"
	"os"

	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/cyverse-de/messaging.v2"
	"gopkg.in/cyverse-de/model.v1"
)

var log = logrus.WithFields(logrus.Fields{"service": "de-job-killer"})

func doKillJob(client *messaging.Client, uuid string) error {
	var err error
	if err = client.SendStopRequest(uuid, "admin", "Sent from de-job-killer."); err != nil {
		return err
	}
	return nil
}

func doStatusMessage(client *messaging.Client, uuid string) error {
	var err error
	fauxJob := &model.Job{
		InvocationID: uuid,
	}
	update := &messaging.UpdateMessage{
		Job:     fauxJob,
		State:   messaging.FailedState,
		Message: "Marked as failed by an admin",
	}
	if err = client.PublishJobUpdate(update); err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		killJob     = flag.Bool("kill", false, "Send out a stop request. Conflicts with --send-status.")
		statusMsg   = flag.Bool("send-status", false, "Send out a job status. Conflicts with --kill.")
		showVersion = flag.Bool("version", false, "Print the version information.")
		config      = flag.String("config", "", "Path to the jobservices config. Required.")
		uuid        = flag.String("uuid", "", "The job UUID to operate against.")
	)

	flag.Parse()

	if *showVersion {
		version.AppVersion()
		os.Exit(0)
	}

	if *config == "" {
		flag.PrintDefaults()
		log.Fatal("--config must be set.")
	}

	if *uuid == "" {
		flag.PrintDefaults()
		log.Fatal("--uuid must be set.")
	}

	if *killJob && *statusMsg {
		log.Fatal("--kill and --send-status conflict.")
	}

	cfg, err := configurate.InitDefaults(*config, configurate.JobServicesDefaults)
	if err != nil {
		log.Fatal(err)
	}

	uri := cfg.GetString("amqp.uri")
	exchangeName := cfg.GetString("amqp.exchange.name")

	client, err := messaging.NewClient(uri, true)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.SetupPublishing(exchangeName)
	if err != nil {
		log.Fatal(err)
	}
	go client.Listen()

	switch {
	case *killJob:
		if err = doKillJob(client, *uuid); err != nil {
			log.Fatal(err)
		}
	case *statusMsg:
		if err = doStatusMessage(client, *uuid); err != nil {
			log.Fatal(err)
		}
	}
}
