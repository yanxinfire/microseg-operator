package types

type SuggestionRule struct {
	Direction   string             `json:"direction"`
	SrcType     string             `json:"srcType"`
	SrcId       uint32             `json:"srcId"`
	SrcResource SuggestionResource `json:"srcResource"`
	SrcSegment  SuggestionSegment  `json:"srcSegment"`
	SrcIPBlock  string             `json:"srcIPBlock"`
	DstType     string             `json:"dstType"`
	DstId       uint32             `json:"dstId"`
	DstResource SuggestionResource `json:"dstResource"`
	DstSegment  SuggestionSegment  `json:"dstSegment"`
	DstIPBlock  string             `json:"dstIPBlock"`
	Action      string             `json:"action"`
	Protocol    string             `json:"protocol"`
	Ports       string             `json:"ports"`
	Comment     string             `json:"comment"`
}

type SuggestionResource struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
}

type SuggestionSegment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
}
