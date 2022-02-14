package fabricGen

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	orgName    string
	endPoint   string
	outDir     string
	peerConut  int
	beginPort  int
	chaincodes []string
}

func NewConfigtx(orgName, endPoint, outDir string, peerCount int, seq int, chaincodes []string) (*Config, error) {
	if peerCount < 1 {
		return nil, errors.New("peer count should > 1")
	}
	outDir, _ = filepath.Abs(outDir)
	cmd := fmt.Sprintf("mkdir -p %s", outDir)
	_, err := ExecShell(cmd)
	if err != nil {
		return nil, err
	}
	return &Config{
		orgName:    orgName,
		endPoint:   endPoint,
		outDir:     outDir,
		peerConut:  peerCount,
		beginPort:  seq,
		chaincodes: chaincodes,
	}, nil
}

func (c *Config) Gen() (err error) {
	err = c.genCryptoConfigFile()
	if err != nil {
		return
	}
	err = c.genConfigtxFile()
	if err != nil {
		return
	}
	err = c.cleanData()
	err = c.genCryptoConfig()
	if err != nil {
		return
	}
	err = c.genConfigTx()
	if err != nil {
		return
	}
	err = c.genComposeFile()
	return
}

func (c *Config) genCryptoConfig() (err error) {
	cmd := cryptogen + "generate --config=/data/crypto-config.yaml --output=/data/crypto-config"
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)
	return
}
func (c *Config) genComposeFile() (err error) {
	beginPort := c.beginPort
	composeStr := compose
	composeStr = strings.ReplaceAll(composeStr, "#[OrgName]", c.orgName)
	composeStr = strings.ReplaceAll(composeStr, "#[Endpoint]", c.endPoint)
	composeStr = strings.ReplaceAll(composeStr, "#[OrdererPort]", strconv.Itoa(beginPort))
	for _, cc := range c.chaincodes {
		composeStr = strings.ReplaceAll(composeStr, "#[CCName]", cc)
		composeStr = strings.ReplaceAll(composeStr, "- ./crypto-config:/opt/crypto-config", "- ./chaincode/#[CCName]:/opt/gopath/src/github.com/#[CCName]\n      - ./crypto-config:/opt/crypto-config")
	}
	composeStr = strings.ReplaceAll(composeStr, "- ./chaincode/#[CCName]:/opt/gopath/src/github.com/#[CCName]", "")

	for i := 0; i < c.peerConut; i++ {
		peerStr := peer
		peerStr = strings.ReplaceAll(peerStr, "#[OrgName]", c.orgName)
		peerStr = strings.ReplaceAll(peerStr, "#[PeerNum]", strconv.Itoa(i))
		peerStr = strings.ReplaceAll(peerStr, "#[Endpoint]", c.endPoint)
		peerStr = strings.ReplaceAll(peerStr, "#[PeerPort]", strconv.Itoa(beginPort+1))
		peerStr = strings.ReplaceAll(peerStr, "#[PeerPort2]", strconv.Itoa(beginPort+2))
		beginPort += 2
		composeStr += peerStr
	}
	err = toFile(c.outDir+"/docker-compose.yaml", composeStr)
	return
}
func (c *Config) genConfigTx() (err error) {
	// Make dir
	cmd := "mkdir -p %s/configtx"
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)

	// Generator system channel
	cmd = configtxgen + "-profile TwoOrgsOrdererGenesis -outputBlock /data/configtx/genesis.block -channelID systemch"
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)

	// Generator common channel
	cmd = configtxgen + "-profile Channel -outputCreateChannelTx /data/configtx/hello.tx -channelID hello"
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)

	// Set anchor peer
	cmd = configtxgen + "-profile Channel -outputAnchorPeersUpdate /data/configtx/#[OrgName]MSPanchors_hello.tx -channelID hello -asOrg #[OrgName]MSP"
	cmd = strings.ReplaceAll(cmd, "#[OrgName]", c.orgName)
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)
	return nil
}

func (c *Config) genConfigtxFile() (err error) {
	tx := configtx
	tx = strings.ReplaceAll(tx, "#[OrgName]", c.orgName)
	tx = strings.ReplaceAll(tx, "#[Endpoint]", c.endPoint)

	err = toFile(c.outDir+"/configtx.yaml", tx)
	return
}

func (c *Config) genCryptoConfigFile() (err error) {
	cryptoConf := cryptoConfig
	cryptoConf = strings.ReplaceAll(cryptoConf, "#[OrgName]", c.orgName)
	cryptoConf = strings.ReplaceAll(cryptoConf, "#[Endpoint]", c.endPoint)
	cryptoConf = strings.ReplaceAll(cryptoConf, "#[PeerCount]", strconv.Itoa(c.peerConut))

	err = toFile(c.outDir+"/crypto-config.yaml", cryptoConf)
	return
}

func (c *Config) cleanData() (err error) {
	t := time.Now().Format(time.RFC3339)
	path := fmt.Sprintf("%s/configtx", c.outDir)
	_ = os.Rename(path, "/Users/fuming/.Trash/configtx "+t)
	path = fmt.Sprintf("%s/crypto-config", c.outDir)
	_ = os.Rename(path, "/Users/fuming/.Trash/crypto-config "+t)
	return
}

func toFile(dir, content string) (err error) {
	file, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	_, err = write.WriteString(content)
	if err != nil {
		return
	}
	// Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	return
}
