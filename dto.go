package main

import "time"

type (
	CustomResponseFolder struct {
		Path       string
		RoutesData map[string]CustomResponseData
	}
	CustomResponse struct {
		Path        string `json:"path"`
		SourceFile  string `json:"source_file"`
		Method      string `json:"method"`
		StatusCode  uint   `json:"status_code"`
		ContentType string `json:"content_type"`
	}
	CustomResponseData struct {
		Path        string
		Content     []byte
		Method      string
		StatusCode  uint
		ContentType string
	}

	Response struct {
		Time         time.Time `json:"time"`
		IP           string    `json:"ip"`
		StartTime    time.Time `json:"startTime"`
		RunningTime  string    `json:"runningTime"`
		RequestCount uint64    `json:"requestCount"`
		Tag          string    `json:"tag,omitempty"`
	}
)
