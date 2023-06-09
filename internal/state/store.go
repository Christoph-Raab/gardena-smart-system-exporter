package state

import (
	"fmt"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena/device"
)

type Store struct {
	devices map[string]device.Device
}

// NewStore creates a new map[string]device.Device
func NewStore() Store {
	var s Store
	s.devices = make(map[string]device.Device)
	return s
}

// StoreDevices adds all devices for a give location state to the store
func (s *Store) StoreDevices(location gardena.State) error {
	devices := s.devicesFrom(location)
	for k, v := range devices {
		if err := s.addDevice(k, v); err != nil {
			return fmt.Errorf("Unable to add all devices to internal store for location %s, got err\n%v", location.Data.Id, err)
		}
	}
	return nil
}

// addDevice adds a given id/map of attributes to the store.
// It uses the factory method of gardena.device to create an
// actual device (MOWER, SENSOR, ...) from the input.
// If a device with the given id is already in the store
// an error is returned.
func (s *Store) addDevice(id string, attrs map[string]any) error {
	if s.devices[id] != nil {
		return fmt.Errorf("device with id %s already in store", id)
	}
	d, err := device.Factory(attrs)
	if err != nil {
		return fmt.Errorf("unable to create device with device with id %s with factory, got err:\n%w", id, err)
	}
	s.devices[id] = d
	return nil
}

// devicesFrom generates a map of devices by merging DEVICE, <type> and COMMON into on map of
// attributes with the attribute name as key. The attribute value is either a str or a float.
// Those devices are collected in a map with the device ID as key.
func (s *Store) devicesFrom(locationData gardena.State) map[string]map[string]any {
	devices := make(map[string]map[string]any)
	for _, d := range locationData.Included {
		m := devices[d.Id]
		if m == nil {
			m = make(map[string]any)
		}
		m[device.AttrId] = d.Id
		if d.Type != device.CommonType && d.Type != device.Type {
			m[device.AttrType] = d.Type
		}
		for k, v := range d.Attributes {
			m[k] = v.Value
		}
		devices[d.Id] = m
	}
	return devices
}
