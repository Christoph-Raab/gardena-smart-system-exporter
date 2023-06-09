package device

import (
	"fmt"
	"reflect"
)

const (
	TypeMower          = "MOWER"
	AttrState          = "state"
	AttrActivity       = "activity"
	AttrOperatingHours = "operatingHours"
)

type Mower struct {
	state          string
	activity       string
	operatingHours float64
	common         Common
}

func (m Mower) GetDeviceId() string {
	return m.common.id
}

func (m Mower) GetFloatAttr(key string) (float64, error) {
	switch key {
	case AttrOperatingHours:
		return m.operatingHours, nil
	default:
		return m.common.getFloatAttr(key)
	}
}

func (m Mower) GetStrAttr(key string) (string, error) {
	switch key {
	case AttrState:
		return m.state, nil
	case AttrActivity:
		return m.activity, nil
	default:
		return m.common.getStrAttr(key)
	}
}

func (m Mower) GetDeviceType() string {
	return TypeMower
}

// MowerFrom creates a Mower af a map of attributes. Attribute values are excepted to be
// interfaces that can be converted with th device.xFromVal methods.
// If a value is of unexpected kind an error is returned.
func MowerFrom(in map[string]any) (Mower, error) {
	var m Mower
	var f float64
	str, err := strFromVal(reflect.ValueOf(in[AttrState]))
	if err != nil {
		return Mower{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrState, in, err)
	}
	m.state = str
	str, err = strFromVal(reflect.ValueOf(in[AttrActivity]))
	if err != nil {
		return Mower{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrActivity, in, err)
	}
	m.activity = str
	f, err = floatFromVal(reflect.ValueOf(in[AttrOperatingHours]))
	if err != nil {
		return Mower{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrOperatingHours, in, err)
	}
	m.operatingHours = f
	c, err := commonFrom(in)
	if err != nil {
		return Mower{}, fmt.Errorf("unable to generate common attributes, got err:\n%w", err)
	}
	m.common = c
	return m, nil
}
