package helpers

type ConditionChangedStatus struct {
	QueryByRoadId                string
	UserDeptId                   int
	Status                       string
	HasPermissionToEditOrApprove bool
	IsInRoad                     bool
}

type UpdateItem struct {
	TableName       string
	UserID          int
	UpdatedDate     string
	UpdateStatus    string
	ConditionStatus string
	IdParent        uint
	RejectReason    string
	QueryCondition  string
}
