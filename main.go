package main

import (
	"fabricGen/generator"
	"log"
)

const workingDir = "/Users/fuming/Downloads/fabricdev2"

func main() {
	conf, _ := generator.NewConfigtx("org1", "flxdu.cn", workingDir+"/org1", 4, 7500, []string{"cc1", "cc2"})
	err := conf.LifecycleDeploy()
	if err != nil {
		log.Fatal(err)
	}
	return
	conf, _ = generator.NewConfigtx("org2", "flxdu2.cn", workingDir+"/org2", 4, 7509, []string{"cc1", "cc2"})
	err = conf.LifecycleDeploy()
	if err != nil {
		log.Fatal(err)
	}
	conf, _ = generator.NewConfigtx("org3", "flxdu3.cn", workingDir+"/org3", 4, 7518, []string{"cc1", "cc2"})
	err = conf.LifecycleDeploy()
	if err != nil {
		log.Fatal(err)
	}
	conf, _ = generator.NewConfigtx("org4", "flxdu4.cn", workingDir+"/org4", 4, 7527, []string{"cc1", "cc2"})
	err = conf.LifecycleDeploy()
	if err != nil {
		log.Fatal(err)
	}
}
