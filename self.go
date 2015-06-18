package main

import (
	"fmt"
	"os"

	"github.com/mperham/inspeqtor/services"
)

func init() {
	services.SupportedInits = append(services.SupportedInits, func() (services.InitSystem, error) {
		return &Self{}, nil
	})
}

// Self provides serivce lookup for the current process, useful
// when running Inspeqtor from the Makefile,
// otherwise it can't find itself as a service.
type Self struct{}

func (m *Self) Name() string { return "self" }

func (m *Self) Restart(name string) error {
	return fmt.Errorf("Cannot restart myself")
}

func (m *Self) Reload(name string) error {
	return fmt.Errorf("Cannot reload myself")
}

func (m *Self) LookupService(name string) (*services.ProcessStatus, error) {
	if name == "inspeqtor" {
		return &services.ProcessStatus{
			Pid:    os.Getpid(),
			Status: services.Up}, nil
	}
	return nil, &services.ServiceError{Init: m.Name(), Name: name, Err: services.ErrServiceNotFound}
}
