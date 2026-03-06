//go:build !linux

package awg

import "syscall"

func setReuseAddr(_ string, _ string, _ syscall.RawConn) error { return nil }

func (p *Proxy) shutdownAllFDs() {
	p.shutdownMu.Lock()
	p.shutdownFDs = p.shutdownFDs[:0]
	p.shutdownMu.Unlock()
}
