package dbs

import (
	"testing"
	"fmt"
)

func TestSQL(t *testing.T) {
	var c1 = SQL("a=?", 1)
	c1.Append("AND b=?", 11)
	fmt.Println(c1.ToSQL())
}

func TestAND(t *testing.T) {
	var wa1 = AND()
	wa1.Appends(SQL("c1=?", 10), SQL("c2=?", 20), OR(SQL("c1=?", 11), SQL("c2=?", 21)), AND(SQL("c3=?", 30)))
	fmt.Println(wa1.ToSQL())

	var wa2 = AND()
	wa2.Append("c1=?", 10000)
	wa2.Append("c2=?", "1111")
	wa2.Appends(SQL("c3=?", 30000))
	fmt.Println(wa2.ToSQL())
}

func TestOR(t *testing.T) {
	var wo1 = OR()
	wo1.Appends(SQL("c1=?", 10), SQL("c2=?", 20), AND(SQL("c3=?", 30), SQL("c4=?", 40)))
	fmt.Println(wo1.ToSQL())
}

func TestIN(t *testing.T) {
	var wi1 = IN("c1", []int{1, 2, 3, 4, 5})
	fmt.Println(wi1.ToSQL())

	var wi2 = AND(IN("c1", []int{1, 2, 3, 4, 5}), SQL("c2=?", 10))
	fmt.Println(wi2.ToSQL())
}

func TestCase(t *testing.T) {
	var c1 = Case("a")
	c1.When("10", SQL("a=?", 1000))
	c1.When("20", "200")
	c1.Else("0")

	fmt.Println(c1.ToSQL())
}