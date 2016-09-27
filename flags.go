package common

import (
        "github.com/codegangsta/cli"
)

var (
        FlAddr = cli.StringFlag{
                Name:  "addr",
                Usage: "<ip>:<port> to listen on",
                Value: "10.254.9.56:9088",
        }

        Fllistaddr = cli.StringFlag{
                Name:  "listaddr",
                Usage: "",
                Value: "http://10.254.9.54:6061",
        }
)