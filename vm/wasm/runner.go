package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

func execHandler(name string) (*wasm.Module, error) {
	switch name {
	case "add":
		raw, err := ioutil.ReadFile("add.wasm")
		if err != nil {
			return nil, err
		}
		add, err := wasm.ReadModule(bytes.NewReader(raw), nil)
		if err != nil {
			return nil, fmt.Errorf("could not read wasm %q module: %v", name, err)
		}
		return add, nil
	case "go":
		// create a whole new module, called "go", from scratch.
		// this module will contain one exported function "print",
		// implemented itself in pure Go.
		print := func(proc *exec.Process, v int32) {
			fmt.Printf("result = %v\n", v)
		}

		m := wasm.NewModule()
		m.Types = &wasm.SectionTypes{
			Entries: []wasm.FunctionSig{
				{
					Form:       0, // value for the 'func' type constructor.
					ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
				},
			},
		}
		m.FunctionIndexSpace = []wasm.Function{
			{
				Sig:  &m.Types.Entries[0],
				Host: reflect.ValueOf(print),
				Body: &wasm.FunctionBody{},
			},
		}
		m.Export = &wasm.SectionExports{
			Entries: map[string]wasm.ExportEntry{
				"print": {
					FieldStr: "print",
					Kind:     wasm.ExternalFunction,
					Index:    0,
				},
			},
		}
		return m, nil
	default:
		return nil, fmt.Errorf("module %q unknown", name)
	}
}

func main() {
	// TODO: Compile from wast to wasm.
	raw, err := ioutil.ReadFile("main.wasm")
	if err != nil {
		panic(err)
	}
	m, err := wasm.ReadModule(bytes.NewReader(raw), execHandler)
	if err != nil {
		panic(err)
	}
	vm, err := exec.NewVM(m)
	if err != nil {
		log.Fatalf("could not create wagon vm: %v", err)
	}
	const fct1 = 2 // index of function fct1.
	out, err := vm.ExecCode(fct1)
	if err != nil {
		log.Fatalf("could not execute fct1(): %v", err)
	}
	fmt.Printf("fct1() -> %v\n", out)

	const fct2 = 3 // index of function fct2.
	out, err = vm.ExecCode(fct2, 40, 6)
	if err != nil {
		log.Fatalf("could not execute fct2(40, 6): %v", err)
	}
	fmt.Printf("fct2() -> %v\n", out)

	const fct3 = 4 // index of function fct3.
	out, err = vm.ExecCode(fct3, 42, 42)
	if err != nil {
		log.Fatalf("could not execute fct3(42, 42): %v", err)
	}
	fmt.Printf("fct3() -> %v\n", out)

}
