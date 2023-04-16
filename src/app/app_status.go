package app

type AppStatus int

const (
	None AppStatus = iota
	Downloading
	PartialDownloaded
	CompleteDownload
	OpenInstaller
)
