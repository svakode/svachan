package framework

import (
	"errors"
	"math"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/minio/minio/pkg/disk"
	"github.com/svakode/svachan/dictionary"
)

type Server struct {
	Memory MemoryService
	CPU    CPUService
	Disk   DiskService
}

func NewServer(m MemoryService, c CPUService, d DiskService) *Server {
	return &Server{
		Memory: m,
		CPU:    c,
		Disk:   d,
	}
}

type MemoryService interface {
	Get() (*Memory, error)
}

type Memory struct {
	Used       float64
	Total      float64
	Percentage float64
}

func NewMemory() MemoryService {
	return &Memory{}
}

func (r *Memory) Get() (*Memory, error) {
	ram, err := memory.Get()
	if err != nil {
		return nil, errors.New(dictionary.GeneralError)
	}

	r.Used = float64(ram.Used) / math.Pow(1024, 3)
	r.Total = float64(ram.Total) / math.Pow(1024, 3)
	r.Percentage = (r.Used / r.Total) * 100
	return r, nil
}

type CPUService interface {
	Get() (*CPU, error)
}

type CPU struct {
	Used float64
}

func NewCPU() CPUService {
	return &CPU{}
}

func (c *CPU) Get() (*CPU, error) {
	cpuInfo, err := cpu.Get()
	if err != nil {
		return nil, errors.New(dictionary.GeneralError)
	}
	c.Used = (float64(cpuInfo.User) + float64(cpuInfo.User)) / float64(cpuInfo.Total) * 100
	return c, nil
}

type DiskService interface {
	Get() (*Disk, error)
}

type Disk struct {
	Free       float64
	Total      float64
	Percentage float64
}

func NewDisk() DiskService {
	return &Disk{}
}

func (d *Disk) Get() (*Disk, error) {
	diskInfo, err := disk.GetInfo("/")
	if err != nil {
		return nil, errors.New(dictionary.GeneralError)
	}
	d.Free = float64(diskInfo.Free) / math.Pow(1024, 3)
	d.Total = float64(diskInfo.Total) / math.Pow(1024, 3)
	d.Percentage = ((d.Total - d.Free) / d.Total) * 100
	return d, nil
}
