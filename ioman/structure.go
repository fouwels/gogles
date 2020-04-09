package ioman

import "time"

//DataPacket ..
type DataPacket struct {
	Valid      bool
	Timestamp  time.Time
	Sensors    Sensors
	Calculated Calculated
	State      EnumState
	Stats      Stats
}

//Sensors ..
type Sensors struct {
	Flow Flow
	ADC  ADC
}

//Flow ..
type Flow struct {
	Val       float64
	CRC       uint8
	Err       error
	Timestamp time.Time
}

//ADC ..
type ADC struct {
	Vals      []uint16
	Timestamp time.Time
	Err       error
}

//Calculated ..
type Calculated struct {
	FlowIntegrated          float64
	FlowIntegratedTimestamp time.Time
}

//EnumState ..
type EnumState int

func (e EnumState) String() string {
	switch int(e) {
	case 0:
		return "Error"
	case 1:
		return "Breathing In"
	case 2:
		return "Breathing Out"
	case 3:
		return "Rest"
	default:
		return "Enum Error"
	}
}

const (
	//StateError ..
	StateError EnumState = iota
	//StateBreathingIn ..
	StateBreathingIn
	//StateBreathingOut ..
	StateBreathingOut
	//StateRest ..
	StateRest
)

//Stats ..
type Stats struct {
	OkReads     uint64
	FailedReads uint64
}
