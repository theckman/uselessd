package uselessecho

import (
	"io/ioutil"
	"net"
	"runtime"
	"testing"

	"github.com/Sirupsen/logrus"

	. "gopkg.in/check.v1"
)

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (*TestSuite) SetUpSuite(c *C) {
	logrus.SetOutput(ioutil.Discard)
}

func (*TestSuite) TestServer_ListenAndServe(c *C) {
	srvr := &Server{
		Host:       "127.0.0.1",
		Port:       10008,
		BufferSize: 9,
	}

	go func() { c.Assert(srvr.ListenAndServe(), IsNil) }()

	// try to ensure the above goroutine has had time to process
	runtime.Gosched()

	//
	// Test UDP echo
	//
	conn, err := net.Dial("udp", "127.0.0.1:10008")
	c.Assert(err, IsNil)

	count, err := conn.Write([]byte("ohaithere"))
	c.Assert(err, IsNil)
	c.Check(count, Equals, 9)

	// change buffer size now for later test we can change
	// this now because the byte slice has already been
	// created with this value for the first connection
	srvr.BufferSize = 4

	buffer := make([]byte, 9)

	count, err = conn.Read(buffer)
	c.Assert(err, IsNil)
	c.Check(count, Equals, 9)
	c.Check(string(buffer), Equals, "ohaithere")

	//
	// Test smaller BufferSize
	//
	count, err = conn.Write([]byte("ohaithere"))
	c.Assert(err, IsNil)
	c.Check(count, Equals, 9)

	buffer = make([]byte, 9)

	count, err = conn.Read(buffer)
	c.Assert(err, IsNil)
	c.Check(count, Equals, 4)
	c.Check(string(buffer), Equals, "ohai\x00\x00\x00\x00\x00")

	c.Assert(conn.Close(), IsNil)

	//
	// Test TCP Echo
	//
	conn, err = net.Dial("tcp", "127.0.0.1:10008")
	c.Assert(err, IsNil)

	count, err = conn.Write([]byte("OHAI"))
	c.Assert(err, IsNil)
	c.Check(count, Equals, 4)

	buffer = make([]byte, 4)

	count, err = conn.Read(buffer)
	c.Assert(err, IsNil)
	c.Check(count, Equals, 4)
	c.Check(string(buffer), Equals, "OHAI")
}
