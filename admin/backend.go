package admin

// ReadBackend returns a slice of all implemented Backend
// names.
func (a *Admin) ReadBackend() ([]string, error) {
	names := a.beReg.List()
	return names, nil
}
