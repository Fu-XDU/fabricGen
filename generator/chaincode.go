package generator

type chaincode struct {
	name      string
	packageID string
	approved  bool
	commited  bool
}

func newChaincode(name string) *chaincode {
	return &chaincode{name: name}
}

func newChaincodes(name []string) map[string]*chaincode {
	res := make(map[string]*chaincode)
	for _, n := range name {
		res[n] = newChaincode(n)
	}
	return res
}
