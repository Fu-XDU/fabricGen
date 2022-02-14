package generator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (c *Config) LifecycleDeploy() (err error) {
	err = c.Gen()
	if err != nil {
		return
	}
	err = c.Up()
	err = c.CreateAndJoinChannel()
	err = c.DeployCC("cc1")
	return nil
}

func (c *Config) Up() (err error) {
	cmd := "PROJECT_NAME=%s docker-compose -f %s/docker-compose.yaml up -d"
	cmd = fmt.Sprintf(cmd, c.orgName, c.outDir)
	_, err = ExecShell(cmd)
	log.Printf("docker compose up success")
	return
}

func (c *Config) CreateAndJoinChannel() (err error) {
	err = c.CreateChannel()
	err = c.JoinChannel()
	return
}

func (c *Config) CreateChannel() (err error) {
	cmd := c.replace(createChannel)
	_, err = ExecShell(cmd)
	cmd = "docker exec cli.#[EndPoint] mv /go/hello.block /opt/configtx/"
	cmd = strings.ReplaceAll(cmd, "#[EndPoint]", c.endPoint)
	_, err = ExecShell(cmd)
	log.Printf("create channel success")
	return
}

func (c *Config) JoinChannel() (err error) {
	cmd := c.replace(joinChannel)
	for i := 0; i < c.peerConut; i++ {
		cmd2 := strings.ReplaceAll(cmd, "#[PeerNum]", strconv.Itoa(i))
		_, err = ExecShell(cmd2)
		log.Printf("peer [peer%s.%s] join channel hello", strconv.Itoa(i), c.endPoint)
	}
	return
}

func (c *Config) DeployCC(ccName string) (err error) {
	err = c.InstallCC(ccName)
	err = c.ApproveCC(ccName)
	err = c.CommitCC(ccName)
	err = c.InitCC(ccName)
	return
}

func (c *Config) InstallCC(ccName string) (err error) {
	if _, ok := c.chaincodes[ccName]; !ok {
		return errors.New(fmt.Sprintf("chaincode [%s] should implement first", ccName))
	}
	cmd := c.replace(installCC)
	cmd = strings.ReplaceAll(cmd, "#[CCName]", ccName)
	for i := 0; i < c.peerConut; i++ {
		cmd2 := strings.ReplaceAll(cmd, "#[PeerNum]", strconv.Itoa(i))
		o, _ := ExecShell(cmd2)
		if c.chaincodes[ccName].packageID == "" {
			c.chaincodes[ccName].packageID = getPackageID(o)
		}
		log.Printf("peer [peer%s.%s] install chaincode [%s]", strconv.Itoa(i), c.endPoint, ccName)
	}
	return
}

func (c *Config) ApproveCC(ccName string) (err error) {
	if _, ok := c.chaincodes[ccName]; !ok {
		return errors.New(fmt.Sprintf("chaincode [%s] should implement first", ccName))
	}

	if c.chaincodes[ccName].packageID == "" {
		return errors.New(fmt.Sprintf("chaincode [%s] should install first", ccName))
	}

	cmd := c.replace(approveCC)
	cmd = strings.ReplaceAll(cmd, "#[CCName]", ccName)
	cmd = strings.ReplaceAll(cmd, "#[PeerNum]", "0")
	cmd = strings.ReplaceAll(cmd, "#[packageID]", c.chaincodes[ccName].packageID)
	_, err = ExecShell(cmd)
	c.chaincodes[ccName].approved = true
	log.Printf("peer [peer0.%s] approve chaincode [%s]", c.endPoint, ccName)
	return
}

func (c *Config) CommitCC(ccName string) (err error) {
	if _, ok := c.chaincodes[ccName]; !ok {
		return errors.New(fmt.Sprintf("chaincode [%s] should implement first", ccName))
	}

	if c.chaincodes[ccName].packageID == "" {
		return errors.New(fmt.Sprintf("chaincode [%s] should install first", ccName))
	}

	if !c.chaincodes[ccName].approved {
		return errors.New(fmt.Sprintf("chaincode [%s] should approve first", ccName))
	}

	cmd := c.replace(commitCC)
	cmd = strings.ReplaceAll(cmd, "#[CCName]", ccName)
	cmd = strings.ReplaceAll(cmd, "#[PeerNum]", "0")
	_, err = ExecShell(cmd)
	c.chaincodes[ccName].commited = true
	log.Printf("peer [peer0.%s] commit chaincode [%s]", c.endPoint, ccName)
	return
}

func (c *Config) InitCC(ccName string) (err error) {
	if _, ok := c.chaincodes[ccName]; !ok {
		return errors.New(fmt.Sprintf("chaincode [%s] should implement first", ccName))
	}

	if c.chaincodes[ccName].packageID == "" {
		return errors.New(fmt.Sprintf("chaincode [%s] should install first", ccName))
	}

	if !c.chaincodes[ccName].approved {
		return errors.New(fmt.Sprintf("chaincode [%s] should approve first", ccName))
	}

	if !c.chaincodes[ccName].commited {
		return errors.New(fmt.Sprintf("chaincode [%s] should commit first", ccName))
	}

	cmd := c.replace(initCC)
	cmd = strings.ReplaceAll(cmd, "#[CCName]", ccName)
	cmd = strings.ReplaceAll(cmd, "#[PeerNum]", "0")
	_, err = ExecShell(cmd)
	log.Printf("peer [peer0.%s] init chaincode [%s]", c.endPoint, ccName)
	return
}

func getPackageID(output string) (id string) {
	s := strings.Split(output, "Chaincode code package identifier: ")
	if len(s) == 1 {
		log.Fatalf("can not get package id from output %s", output)
	}
	id = strings.TrimRight(s[1], "\n")
	return
}
