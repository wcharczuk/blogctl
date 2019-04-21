package mathutil

import "math"

const (
	_pi   = math.Pi
	_2pi  = 2 * math.Pi
	_3pi4 = (3 * math.Pi) / 4.0
	_4pi3 = (4 * math.Pi) / 3.0
	_3pi2 = (3 * math.Pi) / 2.0
	_5pi4 = (5 * math.Pi) / 4.0
	_7pi4 = (7 * math.Pi) / 4.0
	_pi2  = math.Pi / 2.0
	_pi4  = math.Pi / 4.0
	_d2r  = (math.Pi / 180.0)
	_r2d  = (180.0 / math.Pi)

	// Epsilon represents the minimum amount of relevant delta we care about.
	Epsilon = 0.00000001
)
