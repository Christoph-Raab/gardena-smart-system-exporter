package device

import (
	"fmt"
	"reflect"
)

const Type = "DEVICE"

type Device interface {
	GetDeviceId() string
	GetDeviceType() string
	GetFloatAttr(key string) (float64, error)
	GetStrAttr(key string) (string, error)
}

// Factory create a gardena device from a given map of attributes. Currently supported
// devices are:
// - SENSOR
// - MOWER
func Factory(in map[string]any) (Device, error) {
	switch in[AttrType] {
	case TypeSensor:
		s, err := SensorFrom(in)
		if err != nil {
			return nil, fmt.Errorf("unable to create sensor from %v, got err:\n%w", in, err)
		}
		return s, nil
	case TypeMower:
		m, err := MowerFrom(in)
		if err != nil {
			return nil, fmt.Errorf("unable to create mower from %v, got err:\n%w", in, err)
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unsupported device type %s", in[AttrType])
	}
}

// floatFromVal excepts a reflect.Value of kind Float64 and returns the float value
// otherwise it returns an error
func floatFromVal(in reflect.Value) (float64, error) {
	if in.Kind() != reflect.Float64 {
		return 0, fmt.Errorf("excepted float64 value, got %v", in)
	}
	return in.Float(), nil
}

// strFromVal excepts a reflect.Value of kind String and returns the string value
// otherwise it returns an error
func strFromVal(in reflect.Value) (string, error) {
	if in.Kind() != reflect.String {
		return "", fmt.Errorf("excepted string value, got %v", in)
	}
	return in.String(), nil
}
