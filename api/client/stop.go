package client

import (
	"fmt"
	"net/url"
	"strconv"

	flag "github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/utils"
)

// CmdStop stops one or more running containers.
//
// A running container is stopped by first sending SIGTERM and then SIGKILL if the container fails to stop within a grace period (the default is 10 seconds).
//
// Usage: docker stop [OPTIONS] CONTAINER [CONTAINER...]
func (cli *DockerCli) CmdStop(args ...string) error {
	cmd := cli.Subcmd("stop", "CONTAINER [CONTAINER...]", "Stop a running container by sending SIGTERM and then SIGKILL after a\ngrace period", true)
	nSeconds := cmd.Int([]string{"t", "-time"}, 10, "Seconds to wait for stop before killing it")
	cmd.Require(flag.Min, 1)

	utils.ParseFlags(cmd, args, true)

	v := url.Values{}
	v.Set("t", strconv.Itoa(*nSeconds))

	var encounteredError error
	for _, name := range cmd.Args() {
		_, _, err := readBody(cli.call("POST", "/containers/"+name+"/stop?"+v.Encode(), nil, nil))
		if err != nil {
			fmt.Fprintf(cli.err, "%s\n", err)
			encounteredError = fmt.Errorf("Error: failed to stop one or more containers")
		} else {
			fmt.Fprintf(cli.out, "%s\n", name)
		}
	}
	return encounteredError
}
