package main

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
