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
