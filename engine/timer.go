package engine

import (
	"time"
)

type Timer struct {
	currentFrames int
	targetFrames  int
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		currentFrames: 0,
		targetFrames:  int(d.Milliseconds()) * 60 / 1000,
	}
}

func (t *Timer) Update() {
	t.currentFrames++
}

func (t *Timer) IsReady() bool {
	return t.currentFrames >= t.targetFrames
}

func (t *Timer) Reset() {
	t.currentFrames = 0
}
