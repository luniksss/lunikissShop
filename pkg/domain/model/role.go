package model

type Role string

const (
	RoleAnonymous  Role = "anonymous"
	RoleUser       Role = "user"
	RoleSeller     Role = "seller"
	RoleAccountant Role = "accountant"
	RoleAdmin      Role = "admin"
)

func (r Role) HasPermission(required Role) bool {
	rolePriority := map[Role]int{
		RoleAnonymous:  0,
		RoleUser:       1,
		RoleSeller:     2,
		RoleAccountant: 3,
		RoleAdmin:      4,
	}

	return rolePriority[r] >= rolePriority[required]
}
