//+build !windows

package nsocket

import "net"

// questionable quality placeholders
func Listen(name string, opt error) (net.Listener, error)  { return nil, nil }
func DialWindows(name string, opt error) (net.Conn, error) { return nil, nil }
