package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/svakode/svachan/framework"
)

type MemoryService struct {
	mock.Mock
}

func (m *MemoryService) Get() (*framework.Memory, error) {
	res := m.Called()
	if res.Get(0) == nil {
		return nil, res.Error(1)
	}
	return res.Get(0).(*framework.Memory), res.Error(1)
}

type CPUService struct {
	mock.Mock
}

func (c *CPUService) Get() (*framework.CPU, error) {
	res := c.Called()
	if res.Get(0) == nil {
		return nil, res.Error(1)
	}
	return res.Get(0).(*framework.CPU), res.Error(1)
}

type DiskService struct {
	mock.Mock
}

func (d *DiskService) Get() (*framework.Disk, error) {
	res := d.Called()
	if res.Get(0) == nil {
		return nil, res.Error(1)
	}
	return res.Get(0).(*framework.Disk), res.Error(1)
}
