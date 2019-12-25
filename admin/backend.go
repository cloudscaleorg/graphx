package admin

func (a *admin) ReadBackend() ([]string, error) {
	names := a.beReg.List()
	return names, nil
}
