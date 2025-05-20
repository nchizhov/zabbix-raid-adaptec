package main

import (
	"adapter"
	"log"
	"support"
)

func main() {
	log.SetFlags(0)

	args := support.ParseArgs()

	if args.Update {
		adapter.Update()
	}
	if args.GetInfo {
		adapter.GetInfo(args)
	}
	if args.Discovery {
		adapter.Discovery(args)
	}
	if args.Info {
		adapter.Info(args)
	}
}
