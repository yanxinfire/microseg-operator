package types

import "time"

type CreatePolicyRequest struct {
	Status        int    `json:"status"`
	Rules         []Rule `json:"rules"`
	Revision      int    `json:"revision"`
	Operator      string `json:"operator"`
	AllowExternal bool   `json:"allowExternal"`
}

type UpdatePolicyRequest = CreatePolicyRequest

type PolicyResponse struct {
	Rules         []Rule    `json:"rules"`
	Creator       string    `json:"creator"`
	Updater       string    `json:"updater"`
	UpdateTime    time.Time `json:"updateTime"`
	Revision      int       `json:"revision"`
	Total         int       `json:"total"`
	AllowExternal bool      `json:"allowExternal"`
}

type Rule struct {
	Direction   string             `json:"direction"`
	SrcType     string             `json:"srcType"`
	SrcId       uint32             `json:"srcId"`
	SrcIPBlock  string             `json:"srcIPBlock"`
	SrcResource SuggestionResource `json:"srcResource"`
	SrcSegment  SuggestionSegment  `json:"srcSegment"`
	DstType     string             `json:"dstType"`
	DstId       uint32             `json:"dstId"`
	DstIPBlock  string             `json:"dstIPBlock"`
	DstResource SuggestionResource `json:"dstResource"`
	DstSegment  SuggestionSegment  `json:"dstSegment"`
	Action      string             `json:"action"`
	Protocol    string             `json:"protocol"`
	Ports       string             `json:"ports"`
	Comment     string             `json:"comment"`
	IsSegRule   bool               `json:"isSegRule"`
}

type PolicyEnabling struct {
	Enable   int `json:"enable"`
	Revision int `json:"revision"`
}
