package types

type TopologyNode struct {
	NodeType string
	Segment  SingleSegment
	Resource SingleResource
}

type SingleTopologyResource struct {
	SingleResource
	Ingress []SingleResource `json:"ingress"`
	Egress  []SingleResource `json:"egress"`
}
