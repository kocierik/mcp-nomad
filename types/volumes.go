package types

// Volume represents a volume in Nomad
type Volume struct {
	//ID                    string             `json:"ID"`
	Name                  string             `json:"Name"`
	Namespace             string             `json:"Namespace"`
	ExternalID            string             `json:"ExternalID"`
	Topologies            []VolumeTopology   `json:"Topologies"`
	AccessMode            string             `json:"AccessMode"`
	AttachmentMode        string             `json:"AttachmentMode"`
	MountOptions          *MountOptions      `json:"MountOptions,omitempty"`
	Secrets               map[string]string  `json:"Secrets,omitempty"`
	RequestedCapabilities []VolumeCapability `json:"RequestedCapabilities,omitempty"`
	CreateIndex           int                `json:"CreateIndex"`
	ModifyIndex           int                `json:"ModifyIndex"`
}

// VolumeTopology represents the topology of a volume
type VolumeTopology struct {
	Segments map[string]string `json:"Segments"`
}

// MountOptions represents mount options for a volume
type MountOptions struct {
	FSType     string   `json:"FSType,omitempty"`
	MountFlags []string `json:"MountFlags,omitempty"`
}

// VolumeCapability represents a volume capability
type VolumeCapability struct {
	AccessMode     string `json:"AccessMode"`
	AttachmentMode string `json:"AttachmentMode"`
}

// VolumeList represents a list of volumes
type VolumeList struct {
	Volumes []Volume `json:"volumes"`
}

// VolumeClaim represents a volume claim in Nomad
type VolumeClaim struct {
	AllocID       string `json:"AllocID"`
	CreateIndex   int    `json:"CreateIndex"`
	ID            string `json:"ID"`
	JobID         string `json:"JobID"`
	ModifyIndex   int    `json:"ModifyIndex"`
	Namespace     string `json:"Namespace"`
	TaskGroupName string `json:"TaskGroupName"`
	VolumeID      string `json:"VolumeID"`
	VolumeName    string `json:"VolumeName"`
}
