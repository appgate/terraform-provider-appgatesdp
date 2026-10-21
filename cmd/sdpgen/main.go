// Command sdpgen generates the provider from the SDP Admin API specifications.
//
//	sdpgen generate [-version v25]   regenerate all versions, or one version
//	sdpgen pin [-version v25]        pin versions to the current head of their branch
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/appgate/terraform-provider-appgatesdp/internal/codegen"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "sdpgen:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: sdpgen generate|pin [-version vNN[,vNN]] [-root dir]")
	}
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	version := fs.String("version", "", "comma separated API versions (default: all registered versions)")
	root := fs.String("root", ".", "repository root")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	var versions []string
	if *version != "" {
		versions = strings.Split(*version, ",")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	opts := codegen.Options{
		Root:     *root,
		Versions: versions,
		Log:      func(format string, a ...any) { fmt.Fprintf(os.Stderr, format+"\n", a...) },
	}
	switch args[0] {
	case "generate":
		return codegen.Generate(ctx, opts)
	case "pin":
		return codegen.Pin(ctx, opts)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
