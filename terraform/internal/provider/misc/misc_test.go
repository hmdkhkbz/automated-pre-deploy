package misc

import (
	"github.com/praserx/ipconv"
	"net"
	"testing"
)

func TestIPComparison(t *testing.T) {
	start := net.ParseIP("10.255.255.10")
	end := net.ParseIP("10.255.255.40")

	s, err := ipconv.IPv4ToInt(start)
	if err != nil {
		t.Fatal(err)
	}
	e, err := ipconv.IPv4ToInt(end)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(s, e)
}
