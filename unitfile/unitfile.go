// Package unitfile 解析 unit 文件并叠加 drop-in。
package unitfile

import "strings"

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Unit 是解析结果，按节与键走索引。
type Unit struct {
	index map[string]map[string][]string
}

// Parse 读基础文件与 drop-in，后层同名键覆盖前层。
func Parse(base string, dropins []string, c *Counter) (*Unit, error) {
	u := &Unit{index: make(map[string]map[string][]string)}
	u.parseInto(base, c)
	for _, d := range dropins {
		u.parseInto(d, c)
	}
	return u, nil
}

// logicalLines 合并行尾反斜杠的续行。
func logicalLines(text string, c *Counter) []string {
	raws := strings.Split(text, "\n")
	var lines []string
	for i := 0; i < len(raws); i++ {
		c.Scanned++
		line := raws[i]
		for {
			trimmed := strings.TrimRight(line, " \t\r")
			if !strings.HasSuffix(trimmed, "\\") || i+1 >= len(raws) {
				break
			}
			i++
			c.Scanned++
			line = strings.TrimSuffix(trimmed, "\\") + strings.TrimLeft(raws[i], " \t")
		}
		lines = append(lines, line)
	}
	return lines
}

// unquote 去掉值两端成对的引号。
func unquote(v string) string {
	if len(v) >= 2 {
		first, last := v[0], v[len(v)-1]
		if first == last && (first == '"' || first == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}

func (u *Unit) parseInto(text string, c *Counter) {
	section := ""
	inSection := false
	for _, raw := range logicalLines(text, c) {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = unquote(strings.TrimSpace(value))
		sec := u.index[section]
		if sec == nil {
			sec = make(map[string][]string)
			u.index[section] = sec
		}
		sec[key] = append(sec[key], value)
	}
}

// Get 取某个键最后一次出现的值。
func (u *Unit) Get(section, key string, c *Counter) string {
	c.Scanned++
	values := u.index[strings.ToLower(section)][key]
	if len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

// All 取某个键的全部值。
func (u *Unit) All(section, key string, c *Counter) []string {
	c.Scanned++
	values := u.index[strings.ToLower(section)][key]
	out := make([]string, len(values))
	copy(out, values)
	return out
}
