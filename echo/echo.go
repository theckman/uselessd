package uselessecho

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

// DefaultBufferSize is the default buffer size. This avoids large data payloads
// from consuming tons of memory. For UDP connections the rest is discarded,
// whereas TCP will let us take it in separate pieces.
const DefaultBufferSize int = 64

// DefaultPort is the default TCP and UDP port for the echo server.
const DefaultPort uint16 = 7

// Server is a echo server as defined per RFC 862. It takes the data received
// and echoes it back.
type Server struct {
	Host       string
	Port       uint16
	BufferSize int

	udpConn     *net.UDPConn
	tcpListener *net.TCPListener

	mu sync.Mutex
}

func (s *Server) readUDP(data net.PacketConn) ([]byte, net.Addr, error) {
	if s.BufferSize < 1 {
		s.BufferSize = DefaultBufferSize
	}

	buf := make([]byte, s.BufferSize)
	count, addr, err := data.ReadFrom(buf)

	if err != nil {
		return nil, nil, err
	}

	return buf[0:count], addr, nil
}

func (s *Server) handleUDPEcho() {
	for {
		data, addr, err := s.readUDP(s.udpConn)

		if err != nil {
			return
		}

		connUUID := uuid.New()

		logrus.WithFields(logrus.Fields{
			"server":      "echo",
			"proto":       "udp",
			"bytes_read":  len(data),
			"remote_addr": addr.String(),
			"uuid":        connUUID,
		}).Info("new echo read")

		count, err := s.udpConn.WriteTo(data, addr)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"server": "echo",
				"error":  err.Error(),
				"uuid":   connUUID,
			}).Error("error writing data")
			return
		}

		logrus.WithFields(logrus.Fields{
			"server":     "echo",
			"bytes_sent": count,
			"uuid":       connUUID,
		}).Info("new echo replied")
	}
}

func (s *Server) handleTCPEcho(conn net.Conn) {
	connUUID := uuid.New()

	logrus.WithFields(logrus.Fields{
		"server":      "echo",
		"proto":       "tcp",
		"remote_addr": conn.RemoteAddr().String(),
		"conn_uuid":   connUUID,
	}).Info("new connection")

	if s.BufferSize < 1 {
		s.BufferSize = DefaultBufferSize
	}

	for {
		buf := make([]byte, s.BufferSize)
		count, err := conn.Read(buf)

		if err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"server":     "echo",
			"bytes_read": count,
			"conn_uuid":  connUUID,
		}).Info("echo read")

		count, err = conn.Write(buf[0:count])

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"server":    "echo",
				"error":     err.Error(),
				"conn_uuid": connUUID,
			}).Error("error writing data")
			return
		}

		logrus.WithFields(logrus.Fields{
			"server":     "echo",
			"bytes_sent": count,
			"conn_uuid":  connUUID,
		}).Info("echo replied")
	}
}

func (s *Server) tcpListen() {
	if s.tcpListener == nil {
		return
	}

	for {
		conn, err := s.tcpListener.Accept()

		if err != nil {
			return
		}

		go s.handleTCPEcho(conn)
	}
}

// ListenAndServe is the function that listens and serves the echo service.
// This function blocks until Close() is called.
func (s *Server) ListenAndServe() error {
	s.mu.Lock()

	if s.udpConn != nil || s.tcpListener != nil {
		s.mu.Unlock()
		return errors.New("this echo Server looks dirty; please create a new one")
	}

	if s.Port == 0 {
		s.Port = DefaultPort
	}

	tL, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))

	if err != nil {
		s.mu.Unlock()
		return err
	}

	s.tcpListener = tL.(*net.TCPListener)

	uL, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", s.Host, s.Port))

	if err != nil {
		s.mu.Unlock()
		return err
	}

	s.udpConn = uL.(*net.UDPConn)

	s.mu.Unlock()

	go s.handleUDPEcho()
	s.tcpListen()

	return err
}

// Close is a function to close the TCP and UDP listeners for the server.
func (s *Server) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if s.udpConn != nil {
		if err = s.udpConn.Close(); err != nil {
			return
		}
	}

	if s.tcpListener != nil {
		if err = s.tcpListener.Close(); err != nil {
			return
		}
	}

	return
}
