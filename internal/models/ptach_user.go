package models

type PatchUser struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func (pu PatchUser) IsEmpty() bool {
	if pu.FirstName == nil && pu.LastName == nil {
		return true
	}
	return false
}

func (pu *PatchUser) Apply(user *User) {
	if pu.FirstName != nil {
		user.FirstName = *pu.FirstName
	}
	if pu.LastName != nil {
		user.LastName = *pu.LastName
	}
}
