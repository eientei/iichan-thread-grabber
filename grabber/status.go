package grabber

type StatusMessage struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorStatus struct {
	Error string `json:"error"`
}

type ReadyStatus struct {
	Base string `json:"base"`
}

type QueueStatus struct {
	Position int `json:"position"`
}

type DownloadStatus struct {
	TotalDownload   int `json:"total_download"`
	CurrentDownload int `json:"current_download"`
}
