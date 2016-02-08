package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

const arg0 = "uselessd"

func (*TestSuite) TestcommandLine_Parse(c *C) {
	var out string
	var err error

	//
	// Test that version command line flag works
	//
	args := &commandLine{}
	cli := []string{
		arg0, "-V",
	}

	out, err = args.Parse(cli)
	c.Assert(err, IsNil)
	c.Check(out, Equals, "uselessd v0.0.1 built with go1.5.3\nCopyright 2015 Tim Heckman\n")

	args = &commandLine{}
	cli = []string{
		arg0, "--version",
	}

	out, err = args.Parse(cli)
	c.Assert(err, IsNil)
	c.Check(out, Equals, "uselessd v0.0.1 built with go1.5.3\nCopyright 2015 Tim Heckman\n")

	//
	// Test that short flags work
	//
	args = &commandLine{}
	cli = []string{
		arg0,
		"-H", "192.0.2.1",
	}

	out, err = args.Parse(cli)
	c.Assert(err, IsNil)
	c.Check(out, Equals, "")
	c.Check(args.Host, Equals, "192.0.2.1")

	//
	// Test that long flags work
	//
	args = &commandLine{}
	cli = []string{
		arg0,
		"--host", "192.0.2.1",
	}

	out, err = args.Parse(cli)
	c.Assert(err, IsNil)
	c.Check(out, Equals, "")
	c.Check(args.Host, Equals, "192.0.2.1")
}
