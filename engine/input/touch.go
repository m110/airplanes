//go:build ((darwin && !arm && !arm64) || freebsd || linux || windows || js) && !android && !ios

package input

func isTouchPrimaryInput() bool {
	return false
}
