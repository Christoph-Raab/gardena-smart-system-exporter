package device

import (
	"fmt"
	"reflect"
)

const (
	TypeSensor       = "SENSOR"
	AttrSoilHumidity = "soilHumidity"
	AttrSoilTemp     = "soilTemperature"
)

type Sensor struct {
	soilHumidity    float64
	soilTemperature float64
	common          Common
}

func (s Sensor) GetDeviceId() string {
	return s.common.id
}

func (s Sensor) GetFloatAttr(key string) (float64, error) {
	switch key {
	case AttrSoilHumidity:
		return s.soilHumidity, nil
	case AttrSoilTemp:
		return s.soilTemperature, nil
	default:
		return s.common.getFloatAttr(key)
	}
}

func (s Sensor) GetStrAttr(key string) (string, error) {
	// has no string attributes itself
	return s.common.getStrAttr(key)
}

func (s Sensor) GetDeviceType() string {
	return TypeSensor
}

// SensorFrom creates a Sensor af a map of attributes. Attribute values are excepted to be
// interfaces that can be converted with th device.xFromVal methods.
// If a value is of unexpected kind an error is returned.
func SensorFrom(in map[string]any) (Sensor, error) {
	var s Sensor
	f, err := floatFromVal(reflect.ValueOf(in[AttrSoilHumidity]))
	if err != nil {
		return Sensor{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrSoilHumidity, in, err)
	}
	s.soilHumidity = f
	f, err = floatFromVal(reflect.ValueOf(in[AttrSoilTemp]))
	if err != nil {
		return Sensor{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrSoilTemp, in, err)
	}
	s.soilTemperature = f
	c, err := commonFrom(in)
	if err != nil {
		return Sensor{}, fmt.Errorf("unable to generate common attributes, got err:\n%w", err)
	}
	s.common = c
	return s, nil
}
