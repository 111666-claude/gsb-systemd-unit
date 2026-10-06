# systemd-unit

解析 systemd 风格的 unit 文件，并按顺序叠加多份 drop-in：标量后层覆盖前层，列表键按顺序累积。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/unitfile --sample=cont
go run ./cmd/unitfile --sample=overlay
go run ./cmd/unitfile --sample=sectioncase
go run ./cmd/unitfile --sample=quote
go run ./cmd/unitfile --sample=work
```

## 口径

- **续行**：行尾是 `\` 时，下一行去掉行首空白后接在本行后面。
- **节**：`[Section]` 的节名大小写不敏感，查询按小写匹配；任何节之前的内容忽略。
- **叠加**：先读基础文件，再按顺序读 drop-in，后层的同名键覆盖前层。
- **列表键**：`Environment` 这类列表键按出现顺序累积，取值时给全部。
- **引号**：值两端成对的引号去掉。

## 不变量

- 同一节的同名标量取最后一次出现的值。
- 列表键的值按出现顺序返回。
- `scanned` 不随键数乘查询次数放大：八百个键查八百次的 `scanned` 不超过 6000。

## 输出契约

`Get` 取标量（最后一次），`All` 取列表键全部值。场景打印一行 JSON。
`--sample=work` 打印 `{"keys": N, "scanned": N}`。
