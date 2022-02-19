package main

import (
	"fabricGen/generator"
	"fabricGen/utils"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

/*
  BatchTimeout: 2s
  BatchSize:
    MaxMessageCount: 128000
    AbsoluteMaxBytes: 500 MB
    PreferredMaxBytes: 512000 KB
*/

const (
	clientIdentifier = "fabricStarter" // Client identifier to advertise over the network
	clientVersion    = "1.0.0"
	clientUsage      = "fabric starter"
	author           = "Fu Ming"
	email            = "fuming.137@icloud.com"
)

var (
	app       = cli.NewApp()
	baseFlags = []cli.Flag{
		utils.WorkingDirFlag,
		utils.EndPointFlag,
		utils.OrgNameFlag,
		utils.PeerCountFlag,
		utils.SeqFlag,
		utils.ChainCodeFlag,
		utils.BatchTimeoutFlag,
		utils.MaxMessageCountFlag,
		utils.AbsoluteMaxBytesFlag,
		utils.PreferredMaxBytesFlag,
	}
)

func init() {
	app.Action = starter
	app.Name = clientIdentifier
	app.Version = clientVersion
	app.Usage = clientUsage
	app.Author = author
	app.Email = email
	app.Commands = []cli.Command{}
	app.Flags = append(app.Flags, baseFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func starter(ctx *cli.Context) error {
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	err := deploy(ctx.String("o"), ctx.String("e"), ctx.String("d"), ctx.Int("p"), ctx.Int("s"), ctx.StringSlice("cc"), ctx.String("t"), ctx.String("m"), ctx.String("am"), ctx.String("pm"))
	return err
}

func deploy(orgName, endPoint, workingDir string, peerCount, seq int, chaincodes []string, batchTimeout, maxMessageCount, absoluteMaxBytes, preferredMaxBytes string) (err error) {
	batchSize := generator.NewBatchSize(maxMessageCount, absoluteMaxBytes, preferredMaxBytes)
	conf, _ := generator.NewConfigtx(orgName, endPoint, workingDir+"/"+orgName, peerCount, seq, chaincodes, batchTimeout, batchSize)
	err = conf.LifecycleDeploy()
	if err != nil {
		log.Println(err)
	}
	return
}
