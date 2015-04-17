package gen

import (
	// "fmt"
	// "go/parser"
	// "go/token"
	. "gopkg.in/check.v1"
	"log"
	"testing"
)

type TestSuite struct {
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	log.Println("Setting up tests....")
	TestingT(t)
}

func (s *TestSuite) SetUpSuite(c *C) {
}

func (s *TestSuite) SetUpTest(c *C) {
}

func (s *TestSuite) TearDownTest(c *C) {
}

// Tests begin

func (s *TestSuite) TestNewTypeSystem(c *C) {
	ts := NewTypeSystem()
	c.Assert(ts, Not(IsNil))
}

func (s *TestSuite) TestNewType(c *C) {
	t := NewType(BasicType, "", "int64", nil)
	c.Assert(t, Not(IsNil))
}

func (s *TestSuite) TestAddBasicType(c *C) {
	ts := NewTypeSystem()
	t := NewType(BasicType, "", "int64", nil)
	t = ts.AddType(t)
	c.Assert(t.Id, Not(Equals), 0)
	t2 := ts.GetType("", "int64")
	c.Assert(t2, Not(IsNil))
	c.Assert(t.TypeClass, Equals, t2.TypeClass)
	c.Assert(t.Package, Equals, t2.Package)
	c.Assert(t.Name, Equals, t2.Name)
	c.Assert(t.TypeData, Equals, t2.TypeData)
}
