// Package unitfile 解析 unit 文件并叠加 drop-in。
package unitfile

import "strings"

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Unit 是解析结果，按节（小写）与键建索引。
type Unit struct {
	index map[string]map[string][]string
}

// Parse 读基础文件与 drop-in，后层覆盖前层标量、列表键累积。
func Parse(base string, dropins []string, c *Counter) (*Unit, error) {
	u := &Unit{index: make(map[string]map[string][]string)}
	u.parseLayer(base, c)
	for _, d := range dropins {
		u.parseLayer(d, c)
	}
	return u, nil
}

func (u *Unit) parseLayer(text string, c *Counter) {
	section := ""
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		c.Scanned++
		line := strings.TrimSpace(lines[i])
		for strings.HasSuffix(line, `\`) && i+1 < len(lines) {
			i++
			c.Scanned++
			line = strings.TrimSuffix(line, `\`) + strings.TrimLeft(lines[i], " \t")
		}
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			continue
		}
		if section == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		u.add(section, strings.TrimSpace(key), unquote(strings.TrimSpace(value)))
	}
}

func (u *Unit) add(section, key, value string) {
	keys, ok := u.index[section]
	if !ok {
		keys = make(map[string][]string)
		u.index[section] = keys
	}
	keys[key] = append(keys[key], value)
}

// unquote 去掉值两端成对的引号。
func unquote(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}

func (u *Unit) lookup(section, key string, c *Counter) []string {
	c.Scanned++
	return u.index[strings.ToLower(section)][key]
}

// Get 取某个键最后一次出现的值。
func (u *Unit) Get(section, key string, c *Counter) string {
	vals := u.lookup(section, key, c)
	if len(vals) == 0 {
		return ""
	}
	return vals[len(vals)-1]
}

// All 取某个键的全部值。
func (u *Unit) All(section, key string, c *Counter) []string {
	vals := u.lookup(section, key, c)
	return append([]string(nil), vals...)
}
