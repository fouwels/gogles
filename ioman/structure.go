package ioman

import "time"

//DataPacket ..
type DataPacket struct {
	Valid      bool
	Timestamp  time.Time
	Sensors    Sensors
	Calculated Calculated
	State      enumState
	Stats      Stats
}

//Sensors ..
type Sensors struct {
	Flow Flow
}

//Flow ..
type Flow struct {
	Val float32
	Err error
	CRC uint8
}

//Calculated ..
type Calculated struct {
	FlowIntegrated          float64
	FlowIntegratedTimestamp time.Time
}

//EnumState ..
type enumState int

func (e enumState) String() string {
	switch int(e) {
	case 0:
		return "Error"
	case 1:
		return "Breathing In"
	case 2:
		return "Breathing Out"
	case 3:
		return "At Rest"
	default:
		return "Enum Error"
	}
}

const (
	//StateError ..
	StateError enumState = iota
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

type internalPacket struct {
	lastState        enumState
	lastStateChange  time.Time
	flowAverageTotal float64
	flowAverageN     uint64
}
