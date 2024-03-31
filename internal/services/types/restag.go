package types

type ResourceTag = string

const (
	Infrastructure ResourceTag = "Infrastructure"
	Gateway        ResourceTag = "Gateway"
)

type UpdateResourceTagRequest struct {
	Resources []SingleResource `json:"resources"`
}
