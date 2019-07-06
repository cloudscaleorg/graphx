package graphx

// ChartStore stores and retrieves user provided chart configuration.
type ChartStore interface {
	// retrieves all configured Charts keyed by their name
	Get() (map[string]Chart, error)
	// retrieves the Charts specified by the list of Chart names
	GetByNames(chartNames []string) (map[string]Chart, error)
	// stores one or more Chart
	Store(charts []Chart) error
	// removes one or more Chart specified by the Chart
	RemoveByNames(chartNames []string) error
}
