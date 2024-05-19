## Usage
The binary can be found at ```bin``` directory in this repository.

In order to run the executable you need to provide a path to a configuration file, for example:

```./file-downloader ./config.json```

## Configuration
The configuration file contais the following properties:
1. ```filesToDownload``` - A list of JSON object that describes where to download the file from, and what should be the file name on the local machine
2. ```timeoutMilliseconds``` - The maximum amount of time the program should wait for each download before aborting the download
3. ```outputDirectory``` - A file system path on the local machine to store the files in
4. ```maxConcurrency``` - The maximum amount of concurrent downloads

An example configuation file can be found [here](https://github.com/itayankri/file-downloader/blob/master/config.json)

## Build
In order to build the binary just run ```make build```
