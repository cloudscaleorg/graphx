package admin

import (
	"github.com/cloudscaleorg/graphx/etcd"
	"github.com/cloudscaleorg/graphx/registry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Admin exports all methods for administering a GraphX cluster.
//
// This Admin utilizes the etcd package to persist GraphX application resources.
type Admin struct {
	dsmap    *etcd.DSMap
	chartmap *etcd.ChartMap
	beReg    registry.Backend
	logger   zerolog.Logger
}

// NewAdmin is a constructor for an Admin.
func NewAdmin(dsmap *etcd.DSMap, chartmap *etcd.ChartMap, beReg registry.Backend) *Admin {
	return &Admin{
		dsmap:    dsmap,
		chartmap: chartmap,
		beReg:    beReg,
		logger:   log.With().Str("component", "admin").Logger(),
	}
}
