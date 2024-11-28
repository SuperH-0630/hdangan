//go:build !windows

package systeminit

func hidePath(filepath string) error {
	return nil
}
