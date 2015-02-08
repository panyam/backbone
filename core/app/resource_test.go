package app

import (
	. "gopkg.in/check.v1"
)

type ResourceSuite struct{}

var _ = Suite(&ResourceSuite{})

func (s *ResourceSuite) TestUpload(c *C) {
	resHandler, err := NewResourceHandler()
	c.Assert(err, Equals, nil)
	c.Assert(resHandler.Upload("path", "url"), Equals, nil)
}
