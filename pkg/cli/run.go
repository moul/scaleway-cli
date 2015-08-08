// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"strings"

	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdRun = &Command{
	Exec:        runRun,
	UsageLine:   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
	Description: "Run a command in a new server",
	Help:        "Run a command in a new server.",
	Examples: `
    $ scw run ubuntu-trusty
    $ scw run --gateway=myotherserver ubuntu-trusty
    $ scw run ubuntu-trusty bash
    $ scw run --name=mydocker docker docker run moul/nyancat:armhf
    $ scw run --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB bash
    $ scw run --attach alpine
    $ scw run --detach alpine
`,
}

func init() {
	cmdRun.Flag.StringVar(&runCreateName, []string{"-name"}, "", "Assign a name")
	cmdRun.Flag.StringVar(&runCreateBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdRun.Flag.StringVar(&runCreateEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdRun.Flag.StringVar(&runCreateVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdRun.Flag.BoolVar(&runHelpFlag, []string{"h", "-help"}, false, "Print usage")
	cmdRun.Flag.BoolVar(&runAttachFlag, []string{"a", "-attach"}, false, "Attach to serial console")
	cmdRun.Flag.BoolVar(&runDetachFlag, []string{"d", "-detach"}, false, "Run server in background and print server ID")
	cmdRun.Flag.StringVar(&runGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	// FIXME: handle start --timeout
}

// Flags
var runCreateName string       // --name flag
var runCreateBootscript string // --bootscript flag
var runCreateEnv string        // -e, --env flag
var runCreateVolume string     // -v, --volume flag
var runHelpFlag bool           // -h, --help flag
var runAttachFlag bool         // -a, --attach flag
var runDetachFlag bool         // -d, --detach flag
var runGateway string          // -g, --gateway flag

func runRun(cmd *Command, rawArgs []string) {
	if runHelpFlag {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}
	if runAttachFlag && len(rawArgs) > 1 {
		log.Fatalf("Conflicting options: -a and COMMAND")
	}
	if runAttachFlag && runDetachFlag {
		log.Fatalf("Conflicting options: -a and -d")
	}
	if runDetachFlag && len(rawArgs) > 1 {
		log.Fatalf("Conflicting options: -d and COMMAND")
	}

	args := commands.RunArgs{
		Attach:     runAttachFlag,
		Bootscript: runCreateBootscript,
		Command:    rawArgs[1:],
		Detach:     runDetachFlag,
		Gateway:    runGateway,
		Image:      rawArgs[0],
		Name:       runCreateName,
		Tags:       strings.Split(runCreateEnv, " "),
		Volumes:    strings.Split(runCreateVolume, " "),
		// FIXME: DynamicIPRequired
		// FIXME: Timeout
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.Run(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'run': %v", err)
	}
}
