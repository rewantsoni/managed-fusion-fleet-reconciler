package types

type CloudProvider string

const (
	AWSCloudProvider CloudProvider = "aws"
)

type ProviderClusterAWSSpec struct {
	PrimaryCIDR   string
	SecondaryCIDR string
	InstanceCount string
	InstanceType  string
	ImageID       string
}

type ProviderClusterAWSMetadata struct {
	ControlPlaneStackID string
	StackID             string
	ChangeSetID         string
}

type ProviderClusterMetadata struct {
	ProviderClusterAWSMetadata
	LastKnownSpecHash string
}

type ProviderClusterSpec struct {
	Kind CloudProvider
	ProviderClusterAWSSpec
	ROKSVersion       string
	DeletionTimeStamp string
}

type ProviderClusterStatus struct {
}

type ProviderCluster struct {
	ClusterID   string
	AccountID   string
	SatelliteID string
	MetaData    ProviderClusterMetadata
	Spec        ProviderClusterSpec
	Status      ProviderClusterStatus
}
