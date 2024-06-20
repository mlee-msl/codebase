package main

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIsNilPtr(t *testing.T) {
	convey.Convey("case1", t, func() {
		b := IsNilPtr(nil)
		convey.So(b, convey.ShouldBeTrue)
	})

	convey.Convey("case2", t, func() {
		b := IsNilPtr((*int)(nil))
		convey.So(b, convey.ShouldBeTrue)
	})

	convey.Convey("case3", t, func() {
		type (
			st1 struct{}
			st2 struct{ int string }
		)

		b := IsNilPtr((*st1)(nil))
		convey.So(b, convey.ShouldBeTrue)

		b = IsNilPtr((*st2)(nil))
		convey.So(b, convey.ShouldBeTrue)

		b = IsNilPtr(st2{})
		convey.So(b, convey.ShouldBeFalse)

		b = IsNilPtr(&st2{})
		convey.So(b, convey.ShouldBeFalse)
	})

	convey.Convey("case4", t, func() {
		b := IsNilPtr((*[]int)(nil))
		convey.So(b, convey.ShouldBeTrue)

		b = IsNilPtr(([]int)(nil))
		convey.So(b, convey.ShouldBeFalse)

		b = IsNilPtr((chan int)(nil))
		convey.So(b, convey.ShouldBeFalse)

		b = IsNilPtr((map[int]int)(nil))
		convey.So(b, convey.ShouldBeFalse)

		b = IsNilPtr((func())(nil))
		convey.So(b, convey.ShouldBeFalse)

		b = IsNilPtr((any)(nil))
		convey.So(b, convey.ShouldBeTrue)

		b = IsNilPtr(1)
		convey.So(b, convey.ShouldBeFalse)

		b = IsNilPtr("futu")
		convey.So(b, convey.ShouldBeFalse)
	})
}