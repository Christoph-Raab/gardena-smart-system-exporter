package device

import (
	"fmt"
	"reflect"
)

const (
	CommonType       = "COMMON"
	AttrType         = "type"
	AttrId           = "id"
	AttrName         = "name"
	AttrBatteryLevel = "batteryLevel"
	AttrBatteryState = "batteryState"
	AttrRFLinkLevel  = "rfLinkLevel"
	AttrSerial       = "serial"
	AttrModelType    = "modelType"
	AttrRFLinkState  = "rfLinkState"
)

type Common struct {
	id           string
	name         string
	batteryLevel float64
	batteryState string
	rfLinkLevel  float64
	serial       string
	modelType    string
	rfLinkState  string
}

func (c Common) getFloatAttr(key string) (float64, error) {
	switch key {
	case AttrBatteryLevel:
		return c.batteryLevel, nil
	case AttrRFLinkLevel:
		return c.rfLinkLevel, nil
	default:
		return 0, fmt.Errorf("unsupported float attribute %s", key)
	}
}

func (c Common) getStrAttr(key string) (string, error) {
	switch key {
	case AttrName:
		return c.name, nil
	case AttrBatteryState:
		return c.batteryState, nil
	case AttrSerial:
		return c.serial, nil
	case AttrModelType:
		return c.modelType, nil
	case AttrRFLinkState:
		return c.rfLinkState, nil
	default:
		return "", fmt.Errorf("unsuppported string attribute %s", key)
	}
}

func commonFrom(in map[string]any) (Common, error) {
	var c Common
	str, err := strFromVal(reflect.ValueOf(in[AttrId]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrId, in, err)
	}
	c.id = str
	str, err = strFromVal(reflect.ValueOf(in[AttrName]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrName, in, err)
	}
	c.name = str
	f, err := floatFromVal(reflect.ValueOf(in[AttrBatteryLevel]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrBatteryLevel, in, err)
	}
	c.batteryLevel = f
	str, err = strFromVal(reflect.ValueOf(in[AttrBatteryState]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrBatteryState, in, err)
	}
	c.batteryState = str
	f, err = floatFromVal(reflect.ValueOf(in[AttrRFLinkLevel]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrRFLinkLevel, in, err)
	}
	c.rfLinkLevel = f
	str, err = strFromVal(reflect.ValueOf(in[AttrSerial]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrSerial, in, err)
	}
	c.serial = str
	str, err = strFromVal(reflect.ValueOf(in[AttrModelType]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrModelType, in, err)
	}
	c.modelType = str
	str, err = strFromVal(reflect.ValueOf(in[AttrRFLinkState]))
	if err != nil {
		return Common{}, fmt.Errorf("unable to get attr '%s' from map %v, got err:\n%w", AttrRFLinkState, in, err)
	}
	c.rfLinkState = str
	return c, nil
}
