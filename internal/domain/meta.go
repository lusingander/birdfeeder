package domain

import "time"

type MetaRepository interface {
	ReadMeta() (*Meta, error)
}

type Meta struct {
	LastUpdate time.Time
}

func (m *Meta) FormattedLastUpdate() string {
	return m.LastUpdate.Format("2006/01/02 15:04:05")
}
