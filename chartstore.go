package graphx

// ChartStore stores and retrieves user provided chart configuration.
type ChartStore interface {
	Get() ([]*Chart, error)
	GetByNames(chartNames []string) ([]*Chart, error)
	Store(charts []*Chart) error
	RemoveByNames(chartNames []string) error
}
