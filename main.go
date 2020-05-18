package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gitlab-metrics/envutil"
	"gitlab-metrics/gitlab/api"
	"gitlab-metrics/log"
	db2 "gitlab-metrics/metric/updater/db"
	"gitlab-metrics/metric/updater/prometheus"
	"gitlab-metrics/repository"
	"gitlab-metrics/service"
	"gitlab-metrics/service/collector"
	"gitlab-metrics/service/receiver"
)

func main() {
	logger := log.NewLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Fatalf("Failed to start: %v", r)
		}
	}()

	if err := log.SetModeString(os.Getenv(envutil.EnvLogMode)); err != nil {
		panic(err)
	}
	if err := log.SetLevelString(os.Getenv(envutil.EnvLogLevel)); err != nil {
		panic(err)
	}

	gitlabUrl := envutil.GetEnvStr(envutil.EnvGitlabUrl)
	gitlabToken := envutil.GetEnvStr(envutil.EnvGitlabToken)
	updateInterval := envutil.GetEnvUintOrDefault(envutil.EnvUpdateInterval, 30)

	db, err := repository.NewDbFromEnv()
	if err != nil {
		panic(fmt.Errorf("connect to DB: %v", err))
	}

	if err := db.Migrate(); err != nil {
		panic(fmt.Errorf("migrate DB: %v", err))
	}

	eventRepository := repository.NewEventRepository(db)
	projectRepository := repository.NewProjectRepository(db)

	gitlabClient := api.NewClient(gitlabUrl, gitlabToken)

	//mergeRequestReceiver := receiver.NewMergeRequestReceiver(gitlabClient)
	eventReceiver := receiver.NewEventReceiver(gitlabClient)

	//promMergeRequestUpdater := prometheus.NewMergeRequestUpdater()
	promProjectEventsUpdater := prometheus.NewProjectEventsUpdater()

	dbProjectEventsUpdater := db2.NewProjectEventsUpdater(eventRepository, projectRepository)

	//promMergeRequestCollector := collector.NewMergeRequestsCollector(mergeRequestReceiver, promMergeRequestUpdater)

	projectEventsCollector := collector.NewProjectEventsCollector(eventReceiver)
	projectEventsCollector.RegisterUpdater(promProjectEventsUpdater)
	projectEventsCollector.RegisterUpdater(dbProjectEventsUpdater)

	scheduler := service.NewScheduler()
	if err := scheduler.Submit(updateInterval, projectEventsCollector.FetchAndUpdateMetrics); err != nil {
		logger.Errorf("schedule job: %v", err)
	}

	go scheduler.Run()

	http.Handle("/", http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		if _, err := rsp.Write([]byte("ok")); err != nil {
			logger.Errorf("Failed to write response to request: %v", err)
		}
	}))
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
