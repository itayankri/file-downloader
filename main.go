package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/itayankri/file-downloader/httpdownloader"
	"github.com/itayankri/file-downloader/logger"
	"github.com/itayankri/file-downloader/workmanager"
)

type FileToDownload struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
}

type DownloaderConfiguration struct {
	FilesToDownload     []FileToDownload `json:"filesToDownload"`
	TimeoutMilliseconds int64            `json:"timeoutMilliseconds"`
	OutputDirectory     string           `json:"outputDirectory"`
	MaxConcurrency      int              `json:"maxConcurrency"`
}

func main() {
	configurationFilePath := handleCLIArguments()
	data, err := os.ReadFile(configurationFilePath)
	if err != nil {
		fmt.Printf("Failed to read configuration file from %s: %s\n", configurationFilePath, err.Error())
		os.Exit(1)
	}

	var config *DownloaderConfiguration
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Invalid configuration file: %s\n", err.Error())
		os.Exit(1)
	}

	downloadManager := workmanager.WorkManager{
		MaxWorkers: config.MaxConcurrency,
	}

	var fileDownloaders []workmanager.Worker
	for _, file := range config.FilesToDownload {
		fileDownloaders = append(
			fileDownloaders,
			httpdownloader.NewHTTPDownloader(
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
	logger.Info("Total download time: %v\n", totalDownloadTime)
}

func handleCLIArguments() string {
	if len(os.Args) < 2 {
		fmt.Println("Path to configuration file must be provided")
		os.Exit(1)
	}

	return os.Args[1]
}
