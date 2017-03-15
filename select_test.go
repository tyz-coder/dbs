package dba

import (
	"testing"
	"fmt"
)

func TestSelectBuilder_Select(t *testing.T) {
	var sb = NewSelectBuilder()
	sb.Selects("id", "name").Select("email").From("user").Where("id=?", 1)
	fmt.Println(sb.ToSQL())
}