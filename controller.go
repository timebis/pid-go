package pid

import (
	"math"
	"time"
)

// Controller implements a basic PID controller.
type Controller struct {
	// Config for the Controller.
	Config ControllerConfig
	// State of the Controller.
	State ControllerState
}

// ControllerConfig contains configurable parameters for a Controller.
type ControllerConfig struct {
	// ProportionalGain determines ratio of output response to error signal.
	ProportionalGain float64
	// IntegralGain determines previous error's affect on output.
	IntegralGain float64
	// DerivativeGain decreases the sensitivity to large reference changes.
	DerivativeGain float64
	// Timeout is the maximum time between updates before the controller resets.
	Timeout time.Duration
}

// ControllerState holds mutable state for a Controller.
type ControllerState struct {
	// ControlError is the difference between reference and current value.
	ControlError float64
	// ControlErrorIntegral is the integrated control error over time.
	ControlErrorIntegral float64
	// ControlErrorDerivative is the rate of change of the control error.
	ControlErrorDerivative float64
	// ControlSignal is the current control signal output of the controller.
	ControlSignal float64
	// LastUpdateTime is the time of the last update.
	LastUpdateTime time.Time
}

// ControllerInput holds the input parameters to a Controller.
type ControllerInput struct {
	// ReferenceSignal is the reference value for the signal to control.
	ReferenceSignal float64
	// ActualSignal is the actual value of the signal to control.
	ActualSignal float64
}

// Update the controller state.
func (c *Controller) Update(input ControllerInput) {

	if math.IsNaN(input.ReferenceSignal) || math.IsNaN(input.ActualSignal) ||
		math.IsInf(input.ReferenceSignal, 0) || math.IsInf(input.ActualSignal, 0) {
		return
	}

	samplingInterval := time.Since(c.State.LastUpdateTime)
	samplingInterval_sec := samplingInterval.Seconds()
	previousError := c.State.ControlError

	c.State.ControlError = input.ReferenceSignal - input.ActualSignal
	c.State.ControlErrorDerivative = (c.State.ControlError - previousError) / samplingInterval_sec

	if samplingInterval < c.Config.Timeout || c.Config.Timeout == 0 && !c.State.LastUpdateTime.IsZero() {
		c.State.ControlErrorIntegral += c.State.ControlError * samplingInterval_sec
	} else {
		// keep the previous integral term if the timeout is reached : avoid a sudden change in the control signal.
		c.State.ControlErrorIntegral = c.State.ControlErrorIntegral
	}

	c.State.ControlSignal = c.Config.ProportionalGain*c.State.ControlError +
		c.Config.IntegralGain*c.State.ControlErrorIntegral +
		c.Config.DerivativeGain*c.State.ControlErrorDerivative

	c.State.LastUpdateTime = time.Now()

}

// Reset the controller state.
func (c *Controller) Reset() {
	c.State = ControllerState{}
}
