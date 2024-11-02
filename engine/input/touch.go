//go:build !android && !ios

package input

func IsTouchPrimaryInput() bool {
	return false
}
