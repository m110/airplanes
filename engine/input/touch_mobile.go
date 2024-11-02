//go:build (android || ios || (darwin && arm) || (darwin && arm64)) && !js

package input

func isTouchPrimaryInput() bool {
	return true
}
