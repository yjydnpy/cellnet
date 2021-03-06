package main

import (
	"github.com/davyxu/pbmeta"
)

func printFile(gen *Generator, file *pbmeta.FileDescriptor) {

	gen.Println("// Generated by github.com/davyxu/cellnet/protoc-gen-msg")
	gen.Println("// DO NOT EDIT!")
	gen.Println("// Source: ", file.FileName())

	gen.In()
	for _, v := range file.Define.Dependency {
		gen.Println("// ", v)
	}
	gen.Out()

	gen.Println("package ", file.PackageName())
	gen.Println()
	gen.Println("import (")
	gen.In()
	gen.Println("\"github.com/davyxu/cellnet\"")
	gen.Out()
	gen.Println(")")

	gen.Println()
	gen.Println("func init() {")
	gen.In()

	for i := 0; i < file.MessageCount(); i++ {

		msg := file.Message(i)

		gen.Println("cellnet.RegisterMessageMeta(\"", file.PackageName(), ".", msg.Name(), "\", (*", msg.Name(), ")(nil))")
	}

	gen.Out()
	gen.Println("}")

}
