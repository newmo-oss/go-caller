package caller

import "runtime"

func NewFrameForTest(rf runtime.Frame) Frame {
	return Frame{
		frame: rf,
	}
}
