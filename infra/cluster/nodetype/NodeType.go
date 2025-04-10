package nodetype

type NodeType int

const (
	AI_SERVING NodeType = iota
	GATEWAY
	SERVICE
	MOCK
)

type NodeTypeInfo struct {
	ID          string
	DisplayName string
}

var nodeTypeInfos = map[NodeType]NodeTypeInfo{
	AI_SERVING: {ID: "gurms-ai-serving", DisplayName: "Gurms AI Serving"},
	GATEWAY:    {ID: "gurms-gateway", DisplayName: "Gurms Gateway"},
	SERVICE:    {ID: "gurms-service", DisplayName: "Gurms Service"},
	MOCK:       {ID: "gurms-mock", DisplayName: "Gurms Mock"},
}

func (n NodeType) GetId() string {
	if _, ok := nodeTypeInfos[n]; ok {
		return nodeTypeInfos[n].ID
	}
	return "unknown-id"
}

func (n NodeType) GetDisplayName() string {
	if _, ok := nodeTypeInfos[n]; ok {
		return nodeTypeInfos[n].DisplayName
	}
	return "Unknown Display"
}
