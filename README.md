# Beat
Toy program language

#### 処理の流れ
```text
tokens, err := tokenize(input)
nodes, err := parse(tokens)
types, err := typeCheck(nodes)
obj, err := createObj(nodes, types)

// [option] err = exportObj(path, obj)
// [option] obj, err = loadObj(path)

obj, err = link([]obj{obj, ...})
asm, err := compile(obj)

vm, err := newVm(asm)
exitCode, err := vm.run()
```