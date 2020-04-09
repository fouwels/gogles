package ioman

import "time"

type controllerInternals struct {
	lastState            EnumState
	lastStateChangeN     time.Time
	lastStateChangeNmin1 time.Time
	flowAverageTotal     float64
	flowAverageN         uint64
}

type controller struct {
	i controllerInternals
}

func newController() *controller {
	return &controller{
		i: controllerInternals{},
	}
}

func (c *controller) states(sensors Sensors) EnumState {
	breathin := sensors.Flow.Val > _breathInFlowThreshold
	breathout := sensors.Flow.Val < _breathOutFlowThreshold

	state := StateRest
	lastState := c.i.lastState

	if breathin {
		state = StateBreathingIn
		//logf("SPECIAL", "Breathing In %v", sensors.Flow.Val)
	}
	if breathout {
		state = StateBreathingOut
		//logf("SPECIAL", "Breathing Out %v", sensors.Flow.Val)
	}
	if breathin && breathout {
		state = StateError
		logf("ioman:states", "Breathing Error %v", sensors.Flow.Val)
	}

	if state != lastState {
		c.i.lastStateChangeNmin1 = c.i.lastStateChangeN
		c.i.lastStateChangeN = sensors.Flow.Timestamp
	}
	c.i.lastState = state

	return state
}

func (c *controller) calculate(sensors Sensors, state EnumState) Calculated {

	calc := Calculated{}

	// If moving into breathing state, reset counters
	if state == StateBreathingIn && c.i.lastState != StateBreathingIn {
		c.i.flowAverageTotal = 0
		c.i.flowAverageN = 0
	}

	// If in breathing state, increment counters
	if state == StateBreathingIn {
		c.i.flowAverageTotal += float64(sensors.Flow.Val)
		c.i.flowAverageN++
	}

	// If leaving breathing state, calculate integrated flow.
	if state != StateBreathingIn && c.i.lastState == StateBreathingIn && state != StateError {
		duration := sensors.Flow.Timestamp.Sub(c.i.lastStateChangeNmin1)
		flow := (c.i.flowAverageTotal / float64(c.i.flowAverageN)) * duration.Minutes()

		calc.FlowIntegrated = flow
		calc.FlowIntegratedTimestamp = sensors.Flow.Timestamp
	}

	return calc
}
