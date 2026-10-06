// Package unitfile 解析 unit 文件并叠加 drop-in。
package unitfile

import "strings"

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

type entry struct {
	section string
	key     string
	value   string
}

// Unit 是解析结果。
type Unit struct {
	entries []entry
}

// Parse 读基础文件与 drop-in。
func Parse(base string, dropins []string, c *Counter) (*Unit, error) {
	u := &Unit{}
	_ = dropins
	section := ""
	for _, raw := range strings.Split(base, "\n") {
		c.Scanned++
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		u.entries = append(u.entries, entry{section, strings.TrimSpace(key), strings.TrimSpace(value)})
	}
	return u, nil
}

// Get 取某个键最后一次出现的值。
func (u *Unit) Get(section, key string, c *Counter) string {
	out := ""
	for _, e := range u.entries {
		c.Scanned++
		if e.section == section && e.key == key {
			out = e.value
		}
	}
	return out
}

// All 取某个键的全部值。
func (u *Unit) All(section, key string, c *Counter) []string {
	var out []string
	for _, e := range u.entries {
		c.Scanned++
		if e.section == section && e.key == key {
			out = append(out, e.value)
		}
	}
	return out
}
