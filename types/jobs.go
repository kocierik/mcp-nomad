// File: types/jobs.go
package types

// JobSummary represents a summary of a Nomad job
type JobSummary struct {
	ID          string                 `json:"ID"`
	Summary     map[string]TaskSummary `json:"Summary"`
	Children    *JobChildrenSummary    `json:"Children"`
	CreateIndex int                    `json:"CreateIndex"`
	ModifyIndex int                    `json:"ModifyIndex"`
}

// JobSummaryDetails represents detailed summary information for a job
type JobSummaryDetails struct {
	JobID       string                 `json:"JobID"`
	Namespace   string                 `json:"Namespace"`
	Summary     map[string]TaskSummary `json:"Summary"`
	Children    *JobChildrenSummary    `json:"Children"`
	CreateIndex int                    `json:"CreateIndex"`
	ModifyIndex int                    `json:"ModifyIndex"`
}

// TaskSummary represents summary information for a task
type TaskSummary struct {
	Queued   int `json:"Queued"`
	Complete int `json:"Complete"`
	Failed   int `json:"Failed"`
	Running  int `json:"Running"`
	Starting int `json:"Starting"`
	Lost     int `json:"Lost"`
}

// JobChildrenSummary represents summary information for child jobs
type JobChildrenSummary struct {
	Pending int `json:"Pending"`
	Running int `json:"Running"`
	Dead    int `json:"Dead"`
}

// Job represents a Nomad job
type Job struct {
	ID             string            `json:"ID"`
	ParentID       string            `json:"ParentID"`
	Name           string            `json:"Name"`
	Type           string            `json:"Type"`
	Priority       int               `json:"Priority"`
	Status         string            `json:"Status"`
	Datacenters    []string          `json:"Datacenters"`
	NodePool       string            `json:"NodePool"`
	TaskGroups     []TaskGroup       `json:"TaskGroups"`
	Update         *Update           `json:"Update"`
	Periodic       *Periodic         `json:"Periodic"`
	Parameterized  *Parameterized    `json:"Parameterized"`
	Meta           map[string]string `json:"Meta"`
	CreateIndex    int               `json:"CreateIndex"`
	ModifyIndex    int               `json:"ModifyIndex"`
	JobModifyIndex int               `json:"JobModifyIndex"`
}

// Update represents the update strategy for a job
type Update struct {
	Stagger          int    `json:"Stagger"`
	MaxParallel      int    `json:"MaxParallel"`
	HealthCheck      string `json:"HealthCheck"`
	MinHealthyTime   int    `json:"MinHealthyTime"`
	HealthyDeadline  int    `json:"HealthyDeadline"`
	ProgressDeadline int    `json:"ProgressDeadline"`
	AutoRevert       bool   `json:"AutoRevert"`
	Canary           int    `json:"Canary"`
}

// Periodic represents periodic job configuration
type Periodic struct {
	Enabled         bool   `json:"Enabled"`
	Spec            string `json:"Spec"`
	SpecType        string `json:"SpecType"`
	ProhibitOverlap bool   `json:"ProhibitOverlap"`
	TimeZone        string `json:"TimeZone"`
}

// Parameterized represents parameterized job configuration
type Parameterized struct {
	Payload      string   `json:"Payload"`
	MetaRequired []string `json:"MetaRequired"`
	MetaOptional []string `json:"MetaOptional"`
}

// TaskGroupVolume represents a volume configuration for a task group
type TaskGroupVolume struct {
	Name         string                       `json:"Name"`
	Type         string                       `json:"Type"`
	Source       string                       `json:"Source"`
	ReadOnly     bool                         `json:"ReadOnly"`
	MountOptions *TaskGroupVolumeMountOptions `json:"MountOptions"`
}

// TaskGroupVolumeMountOptions represents mount options for a task group volume
type TaskGroupVolumeMountOptions struct {
	FSType     string   `json:"FSType"`
	MountFlags []string `json:"MountFlags"`
}

// TaskGroup represents a task group within a job
type TaskGroup struct {
	Name             string                     `json:"Name"`
	Count            int                        `json:"Count"`
	Tasks            []Task                     `json:"Tasks"`
	Networks         []Network                  `json:"Networks"`
	Services         []Service                  `json:"Services"`
	Volumes          map[string]TaskGroupVolume `json:"Volumes"`
	RestartPolicy    *RestartPolicy             `json:"RestartPolicy"`
	ReschedulePolicy *ReschedulePolicy          `json:"ReschedulePolicy"`
	EphemeralDisk    *EphemeralDisk             `json:"EphemeralDisk"`
	Update           *Update                    `json:"Update"`
	Meta             map[string]string          `json:"Meta"`
}

// Network represents network configuration for a task group
type Network struct {
	Mode          string `json:"Mode"`
	DynamicPorts  []Port `json:"DynamicPorts"`
	ReservedPorts []Port `json:"ReservedPorts"`
}

// Port represents a port configuration
type Port struct {
	Label string `json:"Label"`
	To    int    `json:"To"`
	Value int    `json:"Value"`
}

// Service represents a service definition
type Service struct {
	Name       string            `json:"Name"`
	PortLabel  string            `json:"PortLabel"`
	Tags       []string          `json:"Tags"`
	CanaryTags []string          `json:"CanaryTags"`
	Checks     []ServiceCheck    `json:"Checks"`
	Connect    *ConsulConnect    `json:"Connect"`
	Meta       map[string]string `json:"Meta"`
}

// ServiceCheck represents a service health check
type ServiceCheck struct {
	Name          string   `json:"Name"`
	Type          string   `json:"Type"`
	Command       string   `json:"Command"`
	Args          []string `json:"Args"`
	Path          string   `json:"Path"`
	Protocol      string   `json:"Protocol"`
	PortLabel     string   `json:"PortLabel"`
	Interval      int      `json:"Interval"`
	Timeout       int      `json:"Timeout"`
	InitialStatus string   `json:"InitialStatus"`
	TLSSkipVerify bool     `json:"TLSSkipVerify"`
}

// ConsulConnect represents Consul Connect configuration
type ConsulConnect struct {
	Native bool `json:"Native"`
}

// RestartPolicy represents the restart policy for a task group
type RestartPolicy struct {
	Attempts        int    `json:"Attempts"`
	Interval        int    `json:"Interval"`
	Delay           int    `json:"Delay"`
	Mode            string `json:"Mode"`
	RenderTemplates bool   `json:"RenderTemplates"`
}

// ReschedulePolicy represents the reschedule policy for a task group
type ReschedulePolicy struct {
	Attempts      int    `json:"Attempts"`
	Interval      int    `json:"Interval"`
	Delay         int    `json:"Delay"`
	DelayFunction string `json:"DelayFunction"`
	MaxDelay      int    `json:"MaxDelay"`
	Unlimited     bool   `json:"Unlimited"`
}

// EphemeralDisk represents ephemeral disk configuration
type EphemeralDisk struct {
	Sticky  bool `json:"Sticky"`
	Migrate bool `json:"Migrate"`
	SizeMB  int  `json:"SizeMB"`
}

// Task represents a task within a task group
type Task struct {
	Name            string                 `json:"Name"`
	Driver          string                 `json:"Driver"`
	User            string                 `json:"User"`
	Config          map[string]interface{} `json:"Config"`
	Resources       Resources              `json:"Resources"`
	Services        []Service              `json:"Services"`
	Vault           *Vault                 `json:"Vault"`
	Templates       []Template             `json:"Templates"`
	DispatchPayload *DispatchPayload       `json:"DispatchPayload"`
	Lifecycle       *TaskLifecycle         `json:"Lifecycle"`
	Meta            map[string]string      `json:"Meta"`
}

// Resources represents the resources required by a task
type Resources struct {
	CPU      int               `json:"CPU"`
	MemoryMB int               `json:"MemoryMB"`
	DiskMB   int               `json:"DiskMB"`
	Networks []NetworkResource `json:"Networks"`
	Devices  []Device          `json:"Devices"`
}

// NetworkResource represents network resource requirements
type NetworkResource struct {
	Device        string `json:"Device"`
	CIDR          string `json:"CIDR"`
	IP            string `json:"IP"`
	MBits         int    `json:"MBits"`
	ReservedPorts []Port `json:"ReservedPorts"`
	DynamicPorts  []Port `json:"DynamicPorts"`
}

// Device represents a device resource requirement
type Device struct {
	Name        string       `json:"Name"`
	Count       int          `json:"Count"`
	Constraints []Constraint `json:"Constraints"`
	Affinities  []Affinity   `json:"Affinities"`
}

// Constraint represents a constraint
type Constraint struct {
	LTarget string `json:"LTarget"`
	RTarget string `json:"RTarget"`
	Operand string `json:"Operand"`
}

// Affinity represents an affinity
type Affinity struct {
	LTarget string `json:"LTarget"`
	RTarget string `json:"RTarget"`
	Operand string `json:"Operand"`
	Weight  int    `json:"Weight"`
}

// Vault represents Vault configuration for a task
type Vault struct {
	Policies     []string `json:"Policies"`
	Env          bool     `json:"Env"`
	ChangeMode   string   `json:"ChangeMode"`
	ChangeSignal string   `json:"ChangeSignal"`
}

// Template represents a template configuration
type Template struct {
	SourcePath   string      `json:"SourcePath"`
	DestPath     string      `json:"DestPath"`
	EmbeddedTmpl string      `json:"EmbeddedTmpl"`
	ChangeMode   string      `json:"ChangeMode"`
	ChangeSignal string      `json:"ChangeSignal"`
	Splay        int         `json:"Splay"`
	Perms        string      `json:"Perms"`
	LeftDelim    string      `json:"LeftDelim"`
	RightDelim   string      `json:"RightDelim"`
	Envvars      bool        `json:"Envvars"`
	VaultGrace   int         `json:"VaultGrace"`
	Wait         *WaitConfig `json:"Wait"`
}

// WaitConfig represents template wait configuration
type WaitConfig struct {
	Min string `json:"Min"`
	Max string `json:"Max"`
}

// DispatchPayload represents dispatch payload configuration
type DispatchPayload struct {
	File string `json:"File"`
}

// TaskLifecycle represents task lifecycle configuration
type TaskLifecycle struct {
	Hook    string `json:"Hook"`
	Sidecar bool   `json:"Sidecar"`
}

// Evaluation represents a Nomad evaluation
type Evaluation struct {
	ID                   string                 `json:"ID"`
	Priority             int                    `json:"Priority"`
	Type                 string                 `json:"Type"`
	TriggeredBy          string                 `json:"TriggeredBy"`
	JobID                string                 `json:"JobID"`
	JobModifyIndex       int                    `json:"JobModifyIndex"`
	NodeID               string                 `json:"NodeID"`
	NodeModifyIndex      int                    `json:"NodeModifyIndex"`
	Status               string                 `json:"Status"`
	StatusDescription    string                 `json:"StatusDescription"`
	Wait                 int                    `json:"Wait"`
	NextEvalID           string                 `json:"NextEvalID"`
	PreviousEvalID       string                 `json:"PreviousEvalID"`
	BlockedEvalID        string                 `json:"BlockedEvalID"`
	FailedTGAllocs       map[string]interface{} `json:"FailedTGAllocs"`
	ClassEligibility     map[string]bool        `json:"ClassEligibility"`
	EscapedComputedClass bool                   `json:"EscapedComputedClass"`
	AnnotatePlan         bool                   `json:"AnnotatePlan"`
	QueuedAllocations    map[string]int         `json:"QueuedAllocations"`
	SnapshotIndex        int                    `json:"SnapshotIndex"`
	CreateIndex          int                    `json:"CreateIndex"`
	ModifyIndex          int                    `json:"ModifyIndex"`
}

// JobDeployment represents a Nomad deployment
type JobDeployment struct {
	ID                 string                      `json:"ID"`
	JobID              string                      `json:"JobID"`
	JobVersion         int                         `json:"JobVersion"`
	JobModifyIndex     int                         `json:"JobModifyIndex"`
	JobSpecModifyIndex int                         `json:"JobSpecModifyIndex"`
	JobCreateIndex     int                         `json:"JobCreateIndex"`
	IsMultiregion      bool                        `json:"IsMultiregion"`
	Namespace          string                      `json:"Namespace"`
	Status             string                      `json:"Status"`
	StatusDescription  string                      `json:"StatusDescription"`
	TaskGroups         map[string]*DeploymentState `json:"TaskGroups"`
	CreateIndex        int                         `json:"CreateIndex"`
	ModifyIndex        int                         `json:"ModifyIndex"`
}

// DeploymentState represents the state of a deployment for a task group
type DeploymentState struct {
	AutoRevert        bool   `json:"AutoRevert"`
	ProgressDeadline  int    `json:"ProgressDeadline"`
	RequireProgressBy string `json:"RequireProgressBy"`
	Promoted          bool   `json:"Promoted"`
	DesiredCanaries   int    `json:"DesiredCanaries"`
	DesiredTotal      int    `json:"DesiredTotal"`
	PlacedAllocs      int    `json:"PlacedAllocs"`
	HealthyAllocs     int    `json:"HealthyAllocs"`
	UnhealthyAllocs   int    `json:"UnhealthyAllocs"`
}

// JobPlan represents a Nomad job plan
type JobPlan struct {
	JobModifyIndex     int                    `json:"JobModifyIndex"`
	CreatedEvals       []Evaluation           `json:"CreatedEvals"`
	Diff               *JobDiff               `json:"Diff"`
	Annotations        *PlanAnnotations       `json:"Annotations"`
	FailedTGAllocs     map[string]interface{} `json:"FailedTGAllocs"`
	NextPeriodicLaunch string                 `json:"NextPeriodicLaunch"`
	Warnings           string                 `json:"Warnings"`
}

// JobDiff represents the differences in a job plan
type JobDiff struct {
	Type       string          `json:"Type"`
	ID         string          `json:"ID"`
	Fields     []FieldDiff     `json:"Fields"`
	Objects    []ObjectDiff    `json:"Objects"`
	TaskGroups []TaskGroupDiff `json:"TaskGroups"`
}

// FieldDiff represents a field difference
type FieldDiff struct {
	Type        string   `json:"Type"`
	Name        string   `json:"Name"`
	Old         string   `json:"Old"`
	New         string   `json:"New"`
	Annotations []string `json:"Annotations"`
}

// ObjectDiff represents an object difference
type ObjectDiff struct {
	Type        string       `json:"Type"`
	Name        string       `json:"Name"`
	Fields      []FieldDiff  `json:"Fields"`
	Objects     []ObjectDiff `json:"Objects"`
	Annotations []string     `json:"Annotations"`
}

// TaskGroupDiff represents a task group difference
type TaskGroupDiff struct {
	Type        string         `json:"Type"`
	Name        string         `json:"Name"`
	Fields      []FieldDiff    `json:"Fields"`
	Objects     []ObjectDiff   `json:"Objects"`
	Tasks       []TaskDiff     `json:"Tasks"`
	Updates     map[string]int `json:"Updates"`
	Annotations []string       `json:"Annotations"`
}

// TaskDiff represents a task difference
type TaskDiff struct {
	Type        string       `json:"Type"`
	Name        string       `json:"Name"`
	Fields      []FieldDiff  `json:"Fields"`
	Objects     []ObjectDiff `json:"Objects"`
	Annotations []string     `json:"Annotations"`
}

// PlanAnnotations represents annotations for a plan
type PlanAnnotations struct {
	DesiredTGUpdates map[string]DesiredUpdates `json:"DesiredTGUpdates"`
}

// DesiredUpdates represents desired updates for a task group
type DesiredUpdates struct {
	Ignore            int64 `json:"Ignore"`
	Place             int64 `json:"Place"`
	Migrate           int64 `json:"Migrate"`
	Stop              int64 `json:"Stop"`
	InPlaceUpdate     int64 `json:"InPlaceUpdate"`
	DestructiveUpdate int64 `json:"DestructiveUpdate"`
	Canary            int64 `json:"Canary"`
	Preemptions       int64 `json:"Preemptions"`
}

// JobScaleStatus represents the scale status of a job
type JobScaleStatus struct {
	JobID          string                          `json:"JobID"`
	Namespace      string                          `json:"Namespace"`
	JobModifyIndex int                             `json:"JobModifyIndex"`
	TaskGroups     map[string]TaskGroupScaleStatus `json:"TaskGroups"`
	CreateIndex    int                             `json:"CreateIndex"`
	ModifyIndex    int                             `json:"ModifyIndex"`
}

// TaskGroupScaleStatus represents the scale status of a task group
type TaskGroupScaleStatus struct {
	Desired   int          `json:"Desired"`
	Placed    int          `json:"Placed"`
	Running   int          `json:"Running"`
	Healthy   int          `json:"Healthy"`
	Unhealthy int          `json:"Unhealthy"`
	Events    []ScaleEvent `json:"Events"`
}

// ScaleEvent represents a scaling event
type ScaleEvent struct {
	Time    int64  `json:"Time"`
	Count   int    `json:"Count"`
	Message string `json:"Message"`
	Error   bool   `json:"Error"`
}
