package types

import (
	"time"
)

// Allocation represents a Nomad allocation
type Allocation struct {
	ID                 string                 `json:"ID"`
	EvalID             string                 `json:"EvalID"`
	Name               string                 `json:"Name"`
	NodeID             string                 `json:"NodeID"`
	JobID              string                 `json:"JobID"`
	TaskGroup          string                 `json:"TaskGroup"`
	DesiredStatus      string                 `json:"DesiredStatus"`
	DesiredDescription string                 `json:"DesiredDescription"`
	ClientStatus       string                 `json:"ClientStatus"`
	ClientDescription  string                 `json:"ClientDescription"`
	TaskStates         map[string]TaskState   `json:"TaskStates"`
	DeploymentID       string                 `json:"DeploymentID"`
	DeploymentStatus   *AllocDeploymentStatus `json:"DeploymentStatus"`
	FollowupEvalID     string                 `json:"FollowupEvalID"`
	RescheduleTracker  *RescheduleTracker     `json:"RescheduleTracker"`
	NextAllocation     string                 `json:"NextAllocation"`
	CreateIndex        uint64                 `json:"CreateIndex"`
	ModifyIndex        uint64                 `json:"ModifyIndex"`
	CreateTime         int64                  `json:"CreateTime"`
	ModifyTime         int64                  `json:"ModifyTime"`
}

// TaskState represents the state of a task within an allocation
type TaskState struct {
	State      string      `json:"State"`
	Failed     bool        `json:"Failed"`
	StartedAt  time.Time   `json:"StartedAt"`
	FinishedAt time.Time   `json:"FinishedAt"`
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

// AllocDeploymentStatus represents the deployment status of an allocation
type AllocDeploymentStatus struct {
	Healthy     bool      `json:"Healthy"`
	Timestamp   time.Time `json:"Timestamp"`
	Canary      bool      `json:"Canary"`
	ModifyIndex uint64    `json:"ModifyIndex"`
}

// RescheduleTracker represents the reschedule tracking information for an allocation
type RescheduleTracker struct {
	Events []RescheduleEvent `json:"Events"`
}

// RescheduleEvent represents a reschedule event
type RescheduleEvent struct {
	RescheduleTime time.Time `json:"RescheduleTime"`
	PrevAllocID    string    `json:"PrevAllocID"`
	PrevNodeID     string    `json:"PrevNodeID"`
}
