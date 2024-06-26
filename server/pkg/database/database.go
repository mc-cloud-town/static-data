package database

import (
	"strconv"
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (m *Model) StringID() string {
	return strconv.Itoa(int(m.ID))
}
