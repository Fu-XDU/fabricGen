package utils

import "github.com/urfave/cli"

var (
	WorkingDirFlag = cli.StringFlag{
		Name:   "dir, d",
		Usage:  "Working dir",
		Value:  "./",
		EnvVar: "DIR",
	}
	EndPointFlag = cli.StringFlag{
		Name:   "end, e",
		Usage:  "EndPoint",
		Value:  "example.com",
		EnvVar: "END",
	}
	OrgNameFlag = cli.StringFlag{
		Name:   "org, o",
		Usage:  "Org name",
		Value:  "org1",
		EnvVar: "ORG",
	}
	BatchTimeoutFlag = cli.StringFlag{
		Name:   "timeout, t",
		Usage:  "Batch timeout",
		Value:  "1s",
		EnvVar: "TIMEOUT",
	}
	MaxMessageCountFlag = cli.StringFlag{
		Name:   "maxMessageCount, m",
		Usage:  "Max message count",
		Value:  "500",
		EnvVar: "M",
	}
	AbsoluteMaxBytesFlag = cli.StringFlag{
		Name:   "absoluteMaxBytes, am",
		Usage:  "Absolute max bytes",
		Value:  "98 MB",
		EnvVar: "AM",
	}
	PreferredMaxBytesFlag = cli.StringFlag{
		Name:   "preferredMaxBytes, pm",
		Usage:  "Preferred max bytes",
		Value:  "8193 KB",
		EnvVar: "PM",
	}
	PeerCountFlag = cli.IntFlag{
		Name:   "peercount, p",
		Usage:  "Peer count",
		Value:  4,
		EnvVar: "PEER_COUNT",
	}
	SeqFlag = cli.IntFlag{
		Name:   "seq, s",
		Usage:  "Seq",
		Value:  7050,
		EnvVar: "SEQ",
	}
	ChainCodeFlag = cli.StringSliceFlag{
		Name:   "chaincode, cc",
		Usage:  "Chaincodes",
		EnvVar: "CHAINCODE",
	}
	CpuLimitFlag = cli.Float64Flag{
		Name:   "cpu",
		Usage:  "CPU limits",
		Value:  1,
		EnvVar: "CPU",
	}
)
