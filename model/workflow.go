package model

import (
	"time"
)

type Workflow struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Name        string     `json:"name"`
	Namespace   string     `json:"namespace"`
	Replicas    int32      `json:"replicas"`
	Deployment  string     `json:"deployment"`
	Service     string     `json:"service"`
	Ingress     string     `json:"ingress"`
	ServiceType string     `json:"service_type" gorm:"column:service_type"`
}

func (*Workflow) TableName() string {
	return "workflow"
}
