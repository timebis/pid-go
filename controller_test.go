package pid

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

// Custom assertion function to check if a value is within a specified range
func AssertInRange(t *testing.T, value, min, max float64) {
	assert.Assert(t, value >= min && value <= max,
		fmt.Sprintf("expected %v to be between %v and %v", value, min, max))
}

func TestSpeedControl_ControlLoop_OutputIncrease(t *testing.T) {
	// Given a pidControl with reference value and update interval, dt
	pidControl := Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
		},
		State: ControllerState{
			LastUpdateTime: time.Now(),
		},
	}

	time.Sleep(100 * time.Millisecond)

	// Check output value when output increase is needed
	pidControl.Update(ControllerInput{
		ReferenceSignal: 10,
		ActualSignal:    0,
	})

	AssertInRange(t, pidControl.State.ControlSignal, 120, 122)
	assert.Equal(t, float64(10), pidControl.State.ControlError)
}

func TestTimeout(t *testing.T) {
	// Given a pidControl with reference value and update interval, dt
	pidControl := Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
			Timeout:          1 * time.Second,
		},
		State: ControllerState{
			LastUpdateTime: time.Now().Add(-2 * time.Second),
		},
	}

	// Check output value when output increase is needed
	pidControl.Update(ControllerInput{
		ReferenceSignal: 10,
		ActualSignal:    0,
	})

	AssertInRange(t, pidControl.State.ControlSignal, 24.999, 25.0001)
	assert.Equal(t, float64(10), pidControl.State.ControlError)
}

func TestTimeoutNotSpecified(t *testing.T) {
	// Given a pidControl with reference value and update interval, dt
	pidControl := Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
		},
		State: ControllerState{
			LastUpdateTime: time.Now(),
		},
	}

	time.Sleep(100 * time.Millisecond)

	// Check output value when output increase is needed
	pidControl.Update(ControllerInput{
		ReferenceSignal: 10,
		ActualSignal:    0,
	})

	AssertInRange(t, pidControl.State.ControlSignal, 120, 122)
	assert.Equal(t, float64(10), pidControl.State.ControlError)
}

func TestLastUpdateEmpty(t *testing.T) {
	// Given a pidControl with reference value and update interval, dt
	pidControl := Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
			Timeout:          time.Second,
		},
	}

	time.Sleep(100 * time.Millisecond)

	// Check output value when output increase is needed
	pidControl.Update(ControllerInput{
		ReferenceSignal: 10,
		ActualSignal:    0,
	})

	AssertInRange(t, pidControl.State.ControlSignal, 19.999, 20.0001)
	assert.Equal(t, float64(10), pidControl.State.ControlError)
}

func TestSpeedControl_ControlLoop_OutputConstant(t *testing.T) {
	// Given a pidControl with reference output and update interval, dt
	pidControl := Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
		},
	}
	// Check output value when output value decrease is needed
	pidControl.Update(ControllerInput{
		ReferenceSignal: 10,
		ActualSignal:    10,
	})
	assert.Equal(t, float64(0), pidControl.State.ControlSignal)
	assert.Equal(t, float64(0), pidControl.State.ControlError)
}

func TestSimpleController_Reset(t *testing.T) {
	// Given a Controller with stored values not equal to 0
	c := &Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
		},
		State: ControllerState{
			ControlErrorIntegral:   10,
			ControlErrorDerivative: 10,
			ControlError:           10,
		},
	}
	// And a duplicate Controller with empty values
	expectedController := &Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
		},
	}
	// When resetting stored values
	c.Reset()
	// Then
	assert.Equal(t, expectedController.State, c.State)
}

func TestNaNInput(t *testing.T) {

	// Given a pidControl with reference value and update interval, dt
	pidControl := Controller{
		Config: ControllerConfig{
			ProportionalGain: 2.0,
			IntegralGain:     1.0,
			DerivativeGain:   1.0,
		},
		State: ControllerState{
			ControlError:  11,
			ControlSignal: 122,
		},
	}

	var z float64
	// Check output value when output value decrease is needed
	pidControl.Update(ControllerInput{
		ReferenceSignal: 1 / z,
		ActualSignal:    2,
	})

	assert.Equal(t, float64(122), pidControl.State.ControlSignal)
	assert.Equal(t, float64(11), pidControl.State.ControlError)

	pidControl.Update(ControllerInput{
		ReferenceSignal: 3,
		ActualSignal:    z / z,
	})
	assert.Equal(t, float64(122), pidControl.State.ControlSignal)
	assert.Equal(t, float64(11), pidControl.State.ControlError)

}
