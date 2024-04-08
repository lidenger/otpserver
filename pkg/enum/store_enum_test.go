package enum

import "testing"

func TestConvertStore(t *testing.T) {
	t.Log(ConvertStore("mysql"))
	t.Log(ConvertStore("pgsql"))
	t.Log(ConvertStore("local"))
	t.Log(ConvertStore("memory"))
	t.Log(ConvertStore("xxx"))
}
