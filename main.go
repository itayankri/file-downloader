package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	configurationFilePath := handleCLIArguments()
	data, err := os.ReadFile(configurationFilePath)
	check(err)

	var config *DownloaderConfiguration
	err = json.Unmarshal(data, &config)
	check(err)

	downloadManager := workmanager.WorkManager{
		MaxWorkers: config.MaxConcurrency,
	}

	var fileDownloaders []workmanager.Worker
	for _, file := range config.FilesToDownload {
		fileDownloaders = append(
			fileDownloaders,
			workmanager.NewHTTPDownloader(
				file.URL,
				path.Join(config.OutputDirectory, file.FileName),
				config.TimeoutMilliseconds,
			),
		)
	}

	startTime := time.Now()
	downloadManager.Run(fileDownloaders)
	endTime := time.Now()

	totalDownloadTime := endTime.Sub(startTime).Round(time.Millisecond)
	InfoLog.Printf("Total download time: %v\n", totalDownloadTime)
}

func handleCLIArguments() string {
	if len(os.Args) < 2 {
		fmt.Println("Path to configuration file must be provided")
		os.Exit(1)
	}

	return os.Args[1]
}
