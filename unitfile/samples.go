package unitfile

import (
	"fmt"
	"strings"
)

// Scenario 是一个固定场景。
type Scenario struct {
	Name string
	Run  func() map[string]any
}

// Samples 返回全部场景。
func Samples() []Scenario {
	return []Scenario{
		{Name: "cont", Run: func() map[string]any {
			base := "[Service]\nExecStart=/bin/echo a \\\n  b\n"
			u, _ := Parse(base, nil, &Counter{})
			return map[string]any{"endsb": strings.HasSuffix(u.Get("Service", "ExecStart", &Counter{}), "b")}
		}},
		{Name: "overlay", Run: func() map[string]any {
			base := "[Service]\nX=1\n"
			u, _ := Parse(base, []string{"[Service]\nX=2\n"}, &Counter{})
			return map[string]any{"x": u.Get("Service", "X", &Counter{})}
		}},
		{Name: "sectioncase", Run: func() map[string]any {
			base := "[Service]\nUser=root\n"
			u, _ := Parse(base, nil, &Counter{})
			return map[string]any{"user": u.Get("service", "User", &Counter{})}
		}},
		{Name: "quote", Run: func() map[string]any {
			base := "[Service]\nEnvironment=\"A=1\"\n"
			u, _ := Parse(base, nil, &Counter{})
			return map[string]any{"len": len(u.Get("Service", "Environment", &Counter{}))}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	var b strings.Builder
	b.WriteString("[Service]\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "K%04d=v\n", i)
	}
	u, _ := Parse(b.String(), nil, &Counter{})
	c := &Counter{}
	for i := 0; i < n; i++ {
		u.Get("Service", fmt.Sprintf("K%04d", i), c)
	}
	return map[string]any{"keys": n, "scanned": c.Scanned}
}
