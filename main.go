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

	config, err := loadConfiguration(configurationFilePath)
	check(err)

	totalDownloadTime, err := downloadFiles(config)
	check(err)

	logger.Info("Total download time: %v\n", totalDownloadTime)
}

func handleCLIArguments() string {
	if len(os.Args) < 2 {
		fmt.Println("Path to configuration file must be provided")
		os.Exit(1)
	}

	return os.Args[1]
}

func check(err error) {
	if err != nil {
		fmt.Printf("Execution failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func loadConfiguration(configurationFilePath string) (*DownloaderConfiguration, error) {
	data, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file from %s: %s", configurationFilePath, err.Error())
	}

	var config *DownloaderConfiguration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration file: %s", err.Error())
	}

	exists := isDirectoryExist(config.OutputDirectory)
	if !exists {
		return nil, fmt.Errorf("output directory %s does not exist", config.OutputDirectory)
	}

	return config, nil
}

func isDirectoryExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func downloadFiles(config *DownloaderConfiguration) (time.Duration, error) {
	downloadManager := workmanager.NewWorkManager(config.MaxConcurrency)

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

	return endTime.Sub(startTime).Round(time.Millisecond), nil
}
