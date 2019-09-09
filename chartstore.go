package graphx

// ChartStore stores and retrieves user provided chart configuration.
type ChartStore interface {
	Get() ([]*Chart, error)
	GetByNames(names []string) ([]*Chart, []string, error)
	Store(charts []*Chart) error
	RemoveByNames(names []string) error
}
