package ioman

import (
	"container/ring"
	"time"
)

const _boxcarRatio float64 = 0.1 //division of of sample rate
const _breathInFlowThreshold = 5
const _breathOutFlowThreshold = 0

type calcStore struct {
	flowAverageTotal float64
	flowAverageN     uint64
}

type bufferStore struct {
	flowMovingAverage float64
	flowRing          *ring.Ring //N size circular buffer of previous flow measurements
}
type stateStore struct {
	state           EnumState
	stateChange     time.Time
	lastState       EnumState //state n-1
	lastStateChange time.Time
}

type controller struct {
	state       stateStore
	buffer      bufferStore
	calc        calcStore
	sampledRate time.Duration
}

func newController(sampledRate time.Duration) *controller {

	blen := int(1/sampledRate.Seconds()) * 1
	logf("controller", "Creating a new buffer of len %v", blen)

	fr := ring.New(blen)
	for i := 0; i < fr.Len()+1; i++ {
		fr.Value = float64(0)
		fr = fr.Next()
	}
	return &controller{

		buffer: bufferStore{
			flowRing: fr,
		},

		sampledRate: sampledRate,
	}
}

var counter int = 0

var ticks int = 0

func (c *controller) buffers(sensors Sensors) {
	c.buffer.flowRing = c.buffer.flowRing.Next()
	c.buffer.flowRing.Value = float64(sensors.Flow.Val)

	//Calculate current trailing moving average
	boxcar := int(1 / (c.sampledRate.Seconds()) * _boxcarRatio) //Calculate boxcar width as ratio _boxcarRatio of sample rate
	sum := float64(0)

	c.buffer.flowRing = c.buffer.flowRing.Move(-(boxcar)) //Index to [N - boxcar]
	for s := 0; s < boxcar; s++ {                         //Advance from [N - boxcar] to N
		val := c.buffer.flowRing.Value.(float64)
		sum = sum + val
		c.buffer.flowRing = c.buffer.flowRing.Next()
	}

	c.buffer.flowMovingAverage = sum / float64(boxcar)
}

func (c *controller) states(sensors Sensors) EnumState {

	newstate := StateRest

	if c.buffer.flowMovingAverage > _breathInFlowThreshold {
		newstate = StateBreathingIn
	}

	if newstate != c.state.state {
		c.state.lastStateChange = c.state.stateChange
		c.state.stateChange = sensors.Flow.Timestamp
		logf("controller", "state changed from %v to %v", c.state.lastState, newstate)
	}

	c.state.lastState = c.state.state
	c.state.state = newstate

	return newstate
}

func (c *controller) calculate(sensors Sensors) Calculated {

	calc := Calculated{}

	// If moving into breathing state, reset counters
	if c.state.state == StateBreathingIn && c.state.lastState != StateBreathingIn {
		c.calc.flowAverageTotal = 0
		c.calc.flowAverageN = 0
	}

	// If in breathing state, increment counters
	if c.state.state == StateBreathingIn {
		c.calc.flowAverageTotal += float64(sensors.Flow.Val)
		c.calc.flowAverageN++
	}

	// If leaving breathing state, calculate integrated flow.
	if c.state.state == StateRest && c.state.lastState == StateBreathingIn {
		duration := c.state.stateChange.Sub(c.state.lastStateChange)

		flow := (c.calc.flowAverageTotal / float64(c.calc.flowAverageN)) * duration.Minutes()

		calc.FlowIntegrated = flow
		calc.FlowIntegratedTimestamp = c.state.stateChange

		logf("controller", "breath calculated as %v liters at %v ", calc.FlowIntegrated, calc.FlowIntegratedTimestamp)
	}

	return calc
}
