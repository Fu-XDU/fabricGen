package generator

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
	chaincodes map[string]*chaincode
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
		chaincodes: newChaincodes(chaincodes),
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
	if err != nil {
		return
	}
	err = c.GenSDKFile()
	return
}

func (c *Config) genCryptoConfig() (err error) {
	cmd := cryptogen + "generate --config=/data/crypto-config.yaml --output=/data/crypto-config"
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)
	return
}
func (c *Config) GenSDKFile() (err error) {
	beginPort := c.beginPort
	sdkConf := sdkConfig
	temple1 := "peer#[PeerNum].#[OrgName].#[Endpoint]:\n        <<: *peerPermission"
	temple2 := "peer#[PeerNum].#[OrgName].#[Endpoint]:\n    url: localhost:#[PeerPort]\n    eventUrl: localhost:#[PeerEventPort]\n    grpcOptions:\n      ssl-target-name-override: peer#[PeerNum].#[OrgName].#[Endpoint]\n      hostnameOverride: peer#[PeerNum].#[OrgName].#[Endpoint]\n      <<: *grpcOptions\n    tlsCACerts:\n      path: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/tlsca/tlsca.#[OrgName].#[Endpoint]-cert.pem"
	temple3 := "    - pattern: (\\w+).#[OrgName].#[Endpoint]:(\\d+)\n      urlSubstitutionExp: ${1}.#[OrgName].#[Endpoint]:#[PeerPort]\n      sslTargetOverrideUrlSubstitutionExp: ${1}.#[OrgName].#[Endpoint]\n      mappedHost: peer#[PeerNum].#[OrgName].#[Endpoint]\n"
	temple4 := "      - peer#[PeerNum].#[OrgName].#[Endpoint]\n"
	for i := 0; i < c.peerConut; i++ {
		newStr := fmt.Sprintf("peer%v.#[OrgName].#[Endpoint]:\n        <<: *peerPermission\n      "+temple1, i)
		sdkConf = strings.ReplaceAll(sdkConf, temple1, newStr)

		newStr = temple2 + "\n\n  "
		newStr = strings.ReplaceAll(newStr, "#[PeerNum]", strconv.Itoa(i))
		newStr = strings.ReplaceAll(newStr, "#[PeerPort]", strconv.Itoa(beginPort+1))
		newStr = strings.ReplaceAll(newStr, "#[PeerEventPort]", strconv.Itoa(beginPort+2))

		newStr = newStr + temple2
		sdkConf = strings.ReplaceAll(sdkConf, temple2, newStr)

		newStr = temple3 + "\n"
		newStr = strings.ReplaceAll(newStr, "#[PeerNum]", strconv.Itoa(i))
		newStr = strings.ReplaceAll(newStr, "#[PeerPort]", strconv.Itoa(beginPort+1))
		newStr = newStr + temple3
		sdkConf = strings.ReplaceAll(sdkConf, temple3, newStr)

		newStr = temple4
		newStr = strings.ReplaceAll(newStr, "#[PeerNum]", strconv.Itoa(i))
		newStr = newStr + temple4
		sdkConf = strings.ReplaceAll(sdkConf, temple4, newStr)
		beginPort += 2
	}
	sdkConf = strings.ReplaceAll(sdkConf, temple1, "")
	sdkConf = strings.ReplaceAll(sdkConf, temple2, "")
	sdkConf = strings.ReplaceAll(sdkConf, temple3, "")
	sdkConf = strings.ReplaceAll(sdkConf, temple4, "")
	sdkConf = c.replace(sdkConf)
	sdkConf = strings.ReplaceAll(sdkConf, "#[workingDir]", c.outDir)
	sdkConf = strings.ReplaceAll(sdkConf, "#[OrdererPort]", strconv.Itoa(c.beginPort))
	err = toFile(c.outDir+"/sdk-config.yaml", sdkConf)
	return
}
func (c *Config) genComposeFile() (err error) {
	beginPort := c.beginPort
	composeStr := c.replace(compose)
	composeStr = strings.ReplaceAll(composeStr, "#[OrdererPort]", strconv.Itoa(beginPort))
	for name := range c.chaincodes {
		composeStr = strings.ReplaceAll(composeStr, "#[CCName]", name)
		composeStr = strings.ReplaceAll(composeStr, "- ./crypto-config:/opt/crypto-config", "- ./chaincode/#[CCName]:/opt/gopath/src/github.com/#[CCName]\n      - ./crypto-config:/opt/crypto-config")
	}
	composeStr = strings.ReplaceAll(composeStr, "- ./chaincode/#[CCName]:/opt/gopath/src/github.com/#[CCName]", "")

	for i := 0; i < c.peerConut; i++ {
		peerStr := c.replace(peer)
		peerStr = strings.ReplaceAll(peerStr, "#[PeerNum]", strconv.Itoa(i))
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
	cmd = c.replace(cmd)
	cmd = fmt.Sprintf(cmd, c.outDir)
	_, err = ExecShell(cmd)
	return nil
}

func (c *Config) genConfigtxFile() (err error) {
	tx := c.replace(configtx)

	err = toFile(c.outDir+"/configtx.yaml", tx)
	return
}

func (c *Config) genCryptoConfigFile() (err error) {
	cryptoConf := c.replace(cryptoConfig)

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
func (c *Config) replace(str string) string {
	str = strings.ReplaceAll(str, "#[OrgName]", c.orgName)
	str = strings.ReplaceAll(str, "#[Endpoint]", c.endPoint)
	str = strings.ReplaceAll(str, "#[PeerCount]", strconv.Itoa(c.peerConut))
	return str
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
