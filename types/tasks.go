package types

import (
	"time"
)

// TaskState represents the state of a task within an allocation
type TaskState struct {
	State      string      `json:"State"`
	Failed     bool        `json:"Failed"`
	StartedAt  *time.Time  `json:"StartedAt"`
	FinishedAt *time.Time  `json:"FinishedAt"`
	Events     []TaskEvent `json:"Events"`
}

// TaskEvent represents an event that occurred for a task
type TaskEvent struct {
	Type             string `json:"Type"`
	Time             int64  `json:"Time"`
	FailsTask        bool   `json:"FailsTask"`
	RestartReason    string `json:"RestartReason"`
	SetupError       string `json:"SetupError"`
	DriverError      string `json:"DriverError"`
	ExitCode         int    `json:"ExitCode"`
	Signal           int    `json:"Signal"`
	Message          string `json:"Message"`
	KillReason       string `json:"KillReason"`
	KillTimeout      int64  `json:"KillTimeout"`
	KillError        string `json:"KillError"`
	StartDelay       int64  `json:"StartDelay"`
	DownloadError    string `json:"DownloadError"`
	ValidationError  string `json:"ValidationError"`
	DiskLimit        int64  `json:"DiskLimit"`
	FailedSibling    string `json:"FailedSibling"`
	VaultError       string `json:"VaultError"`
	TaskSignalReason string `json:"TaskSignalReason"`
	TaskSignal       string `json:"TaskSignal"`
	DriverMessage    string `json:"DriverMessage"`
	GenericSource    string `json:"GenericSource"`
}
