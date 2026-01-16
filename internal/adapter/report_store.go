package adapter

import (
	m "github.com/mouse-blink/gooze/internal/model"
)

type ReportStore interface {
	SaveReports(path m.Path, reports []m.ReportV2) error
	LoadReports(path m.Path) ([]m.ReportV2, error)
}

type reportStore struct{}

func NewReportStore() ReportStore {
	return &reportStore{}
}

func (rs *reportStore) SaveReports(path m.Path, reports []m.ReportV2) error {
	return nil
}

func (rs *reportStore) LoadReports(path m.Path) ([]m.ReportV2, error) {
	return nil, nil
}
