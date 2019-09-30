package graphx

// DataSourceStore stores and retrieves user provided DataSource configuration
type DataSourceStore interface {
	Get() ([]*DataSource, error)
	GetByNames(names []string) ([]*DataSource, []string, error)
	Store(source []*DataSource) error
	RemoveByNames(names []string) error
}
