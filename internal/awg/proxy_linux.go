//go:build linux

package awg

import "syscall"

// setReuseAddr sets SO_REUSEADDR on the socket before bind.
// Prevents EADDRINUSE when reconnecting while dup'd fd still holds the port.
func setReuseAddr(_ string, _ string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})
	return err
}

// shutdownAllFDs closes all registered blocking fds to unblock recvmmsg/sendmmsg.
// shutdown(SHUT_RDWR) alone doesn't reliably unblock UDP recvmmsg;
// close(fd) guarantees EBADF wakeup.
func (p *Proxy) shutdownAllFDs() {
	p.shutdownMu.Lock()
	for _, fd := range p.shutdownFDs {
		syscall.Shutdown(fd, syscall.SHUT_RDWR)
		syscall.Close(fd)
	}
	p.shutdownFDs = p.shutdownFDs[:0]
	p.shutdownMu.Unlock()
}
