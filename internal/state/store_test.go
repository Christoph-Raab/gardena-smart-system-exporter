package state

import (
	"encoding/json"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena/device"
	"os"
	"testing"
)

func TestStoreDevicesFromState(t *testing.T) {
	location, err := os.ReadFile("../../test/location.json")
	if err != nil {
		t.Fatal("Unable to read location.json file", err)
	}
	state := gardena.State{}
	if err := json.Unmarshal(location, &state); err != nil {
		t.Fatal("Unable to create state from json", err)
	}
	s := NewStore()
	err = s.StoreDevices(state)
	if err != nil {
		t.Fatal("Unable to store state", err)
	}

	// Start testing the store
	if len(s.devices) != 2 {
		t.Fatalf("Excepted two devices, found %d", len(s.devices))
	}

	// Verify sensor
	idOne := "dev-1-id"
	sensor := s.devices[idOne]
	id := sensor.GetDeviceId()
	if id != idOne {
		t.Fatalf("Excepted sensor to have id %s, got %s", idOne, id)
	}
	deviceType := sensor.GetDeviceType()
	if deviceType != device.TypeSensor {
		t.Fatalf("Expected device with id %s to be of type %s, got %s", idOne, device.TypeSensor, deviceType)
	}
	soH, err := sensor.GetFloatAttr(device.AttrSoilHumidity)
	if err != nil {
		t.Fatalf("Unexcepted error for attr %s:\n%v", device.AttrSoilHumidity, err)
	}
	if soH != 95 {
		t.Fatalf("Excepted 95 as value for %s, got %v", device.AttrSoilHumidity, soH)
	}
	n, err := sensor.GetStrAttr(device.AttrName)
	if err != nil {
		t.Fatalf("Unexcepted error for attr %s:\n%v", device.AttrName, err)
	}
	if n != "Sensor01" {
		t.Fatalf("Excepted Sensor01 as value for %s, got %v", device.AttrName, n)
	}

	// Verify mower
	idTwo := "dev-2-id"
	mower := s.devices[idTwo]
	id = mower.GetDeviceId()
	if id != idTwo {
		t.Fatalf("Excepted mower to have id %s, got %s", idTwo, id)
	}
	deviceType = mower.GetDeviceType()
	if deviceType != device.TypeMower {
		t.Fatalf("Expected device with id %s to be of type %s, got %s", idTwo, device.TypeMower, deviceType)
	}
	opH, err := mower.GetFloatAttr(device.AttrOperatingHours)
	if err != nil {
		t.Fatalf("Unexcepted error for attr %s:\n%v", device.AttrOperatingHours, err)
	}
	if opH != 435 {
		t.Fatalf("Excepted 435 as value for %s, got %v", device.AttrOperatingHours, soH)
	}
	n, err = mower.GetStrAttr(device.AttrName)
	if err != nil {
		t.Fatalf("Unexcepted error for attr %s:\n%v", device.AttrName, err)
	}
	if n != "SILENO" {
		t.Fatalf("Excepted SILENO as value for %s, got %v", device.AttrName, n)
	}

	// Verify failure
	str, err := mower.GetStrAttr("rand")
	if str != "" || err == nil {
		t.Fatalf("Expected empty value and error for str attribute with key 'rand'")
	}
	f, err := sensor.GetFloatAttr("foo")
	if f != 0 || err == nil {
		t.Fatalf("Expected empty value and error for float attribute with key 'foo'")
	}
}
