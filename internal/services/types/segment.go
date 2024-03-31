package types

type CreateFakeResourceRequest struct {
	Name string `json:"name"`
}

type CreateSegmentRequest struct {
	Name string `json:"name"`
}

type UpdateResourceOfSegmentRequest struct {
	Resources []K8sResource `json:"resources"`
}

type K8sResource struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
}

type SingleResource struct {
	ID          uint32 `json:"id"`
	SegmentID   uint32 `json:"segmentId"`
	SegmentName string `json:"segmentName"`
	IsFake      bool   `json:"isFake"`
	DstPort     int    `json:"dstPort"`
	Protocol    string `json:"protocol"`
	K8sResource
}

type SingleSegment struct {
	ID        uint32           `json:"id"`
	Name      string           `json:"name"`
	Namespace string           `json:"namespace"`
	Cluster   string           `json:"cluster"`
	Resources []SingleResource `json:"resources"`
}

type SingleNamespace struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
}
