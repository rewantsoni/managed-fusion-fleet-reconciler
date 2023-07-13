package types

type ProviderClusterMetadata struct {
}
type ProviderClusterSpec struct {
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
