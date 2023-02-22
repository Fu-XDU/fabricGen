package utils

import "github.com/urfave/cli"

var (
	WorkingDirFlag = cli.StringFlag{
		Name:  "dir, d",
		Usage: "Working dir",
		Value: "./",
	}
	EndPointFlag = cli.StringFlag{
		Name:  "end, e",
		Usage: "EndPoint",
		Value: "example.com",
	}
	OrgNameFlag = cli.StringFlag{
		Name:  "org, o",
		Usage: "Org name",
		Value: "org1",
	}
	BatchTimeoutFlag = cli.StringFlag{
		Name:  "timeout, t",
		Usage: "Batch timeout",
		Value: "1s",
	}
	MaxMessageCountFlag = cli.StringFlag{
		Name:  "maxMessageCount, m",
		Usage: "Max message count",
		Value: "500",
	}
	AbsoluteMaxBytesFlag = cli.StringFlag{
		Name:  "absoluteMaxBytes, am",
		Usage: "Absolute max bytes",
		Value: "98 MB",
	}
	PreferredMaxBytesFlag = cli.StringFlag{
		Name:  "preferredMaxBytes, pm",
		Usage: "Preferred max bytes",
		Value: "8193 KB",
	}
	PeerCountFlag = cli.IntFlag{
		Name:  "peercount, p",
		Usage: "Peer count",
		Value: 4,
	}
	SeqFlag = cli.IntFlag{
		Name:  "seq, s",
		Usage: "Seq",
		Value: 7050,
	}
	ChainCodeFlag = cli.StringSliceFlag{
		Name:  "chaincode, cc",
		Usage: "Chaincodes",
	}
)
