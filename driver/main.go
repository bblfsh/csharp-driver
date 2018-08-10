package main

import (
	_ "github.com/bblfsh/csharp-driver/driver/impl"
	"github.com/bblfsh/csharp-driver/driver/normalizer"

	"gopkg.in/bblfsh/sdk.v2/driver/server"
)

func main() {
	server.Run(normalizer.Transforms)
}
