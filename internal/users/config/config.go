package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Pagination struct {
	mu           sync.RWMutex
	defaultLimit uint64
	maxLimit     uint64
}

func NewPaginationConfig(v *viper.Viper) *Pagination {
	p := &Pagination{}
	p.UpdatePaginationConfig(v)
	return p
}

func (p *Pagination) UpdatePaginationConfig(v *viper.Viper) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.defaultLimit = v.GetUint64("api.default_limit")
	p.maxLimit = v.GetUint64("api.max_allowed_limit")
}

func (p *Pagination) DefaultLimit() uint64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.defaultLimit
}

func (p *Pagination) MaxLimit() uint64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.maxLimit
}
