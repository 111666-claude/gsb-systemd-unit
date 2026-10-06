package unitfile

import "testing"

func TestSectionStored(t *testing.T) {
	u, _ := Parse("[Service]\nUser=root\n", nil, &Counter{})
	if got := u.Get("Service", "User", &Counter{}); got != "root" {
		t.Fatalf("User 应是 root：%q", got)
	}
}

func TestGetMissing(t *testing.T) {
	u, _ := Parse("[Service]\nUser=root\n", nil, &Counter{})
	if got := u.Get("Service", "Nope", &Counter{}); got != "" {
		t.Fatalf("缺失键应是空：%q", got)
	}
}

func TestAllCounts(t *testing.T) {
	u, _ := Parse("[Service]\nEnvironment=A=1\n", nil, &Counter{})
	if got := len(u.All("Service", "Environment", &Counter{})); got != 1 {
		t.Fatalf("应有一个值：%d", got)
	}
}

func TestValueTrimmed(t *testing.T) {
	u, _ := Parse("[Service]\nX= 1 \n", nil, &Counter{})
	if got := u.Get("Service", "X", &Counter{}); got != "1" {
		t.Fatalf("值应去掉两端空白：%q", got)
	}
}

func TestCommentSkipped(t *testing.T) {
	u, _ := Parse("# c\n[Service]\nX=1\n", nil, &Counter{})
	if got := u.Get("Service", "X", &Counter{}); got != "1" {
		t.Fatalf("注释不该影响解析：%q", got)
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	Parse("[Service]\nX=1\n", nil, c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
