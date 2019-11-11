package system

func (user *User)HasPerm(permission string) bool  {
	return true
}

func (user *User)HasRole(roles string) bool {
	return true
}