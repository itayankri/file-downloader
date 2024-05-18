package httpdownloader

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/itayankri/file-downloader/logger"
	"github.com/itayankri/file-downloader/workmanager"
)

var _ workmanager.Worker = (*HTTPDownloader)(nil)

type HTTPDownloader struct {
	url        string
	outputPath string
	client     http.Client
}

func (d *HTTPDownloader) Work(ctx context.Context) error {
	contextId := ctx.Value(workmanager.WORKER_ID_KEY)

	out, err := os.Create(d.outputPath)
	if err != nil {
		logger.Error("[%s] Failed to create a file at %s: %s\n", contextId, d.outputPath, err.Error())
		return err
	}
	defer out.Close()

	logger.Info("[%s] Start downloading file from %s\n", contextId, d.url)
	startTime := time.Now()

	response, err := d.client.Get(d.url)
	if err != nil {
		logger.Error("[%s] Download file from %s failed: %s\n", contextId, d.url, err.Error())
		os.Remove(d.outputPath)
		return err
	}
	defer response.Body.Close()

	endTime := time.Now()
	downloadTime := endTime.Sub(startTime).Round(time.Millisecond)
	logger.Info("[%s] File download from %s completed in %v, %s\n", contextId, d.url, downloadTime, d.outputPath)

	_, err = io.Copy(out, response.Body)
	if err != nil {
		logger.Error("[%s] Failed to create a file at %s: %s\n", contextId, d.outputPath, err.Error())
		return err
	}

	return nil
}

func NewHTTPDownloader(url string, outputPath string, timeout int64) *HTTPDownloader {
	return &HTTPDownloader{
		url:        url,
		outputPath: outputPath,
		client: http.Client{
			Timeout: time.Duration(timeout * int64(time.Millisecond)),
		},
	}
}
