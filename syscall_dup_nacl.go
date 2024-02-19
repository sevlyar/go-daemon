// +build nacl

package daemon

func syscallDup(oldfd int, newfd int) (err error) {
	return nil
}
