package model

func (i *InstanceStruct) IsMyName(name string) bool {
	if i.InstanceAlias != "" && name == i.InstanceAlias {
		return true
	}
	if i.InstanceName == name {
		return true
	}
	return false
}
