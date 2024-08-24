package pid_test

import (
	"fmt"
	"time"

	"go.einride.tech/pid"
)

func ExampleController() {
	// Create a PID controller.
	c := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
			Timeout:          time.Second,
		},
	}
	// Update the PID controller.
	c.Update(pid.ControllerInput{
		ReferenceSignal: 10,
		ActualSignal:    0,
	})
	fmt.Printf("%+v\n", c.State)
	// Reset the PID controller.
	c.Reset()
	fmt.Printf("%+v\n", c.State)

}
