package generator

const (
	configtx = `Organizations:
  - &OrdererOrg
    Name: OrdererOrg
    ID: OrdererMSP
    MSPDir: crypto-config/ordererOrganizations/msp
    Policies:
      Readers:
        Type: Signature
        Rule: "OR('OrdererMSP.member')"
      Writers:
        Type: Signature
        Rule: "OR('OrdererMSP.member')"
      Admins:
        Type: Signature
        Rule: "OR('OrdererMSP.admin')"

  - &Org
    Name: #[OrgNameUpper]MSP
    ID: #[OrgName]MSP
    MSPDir: crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/msp
    AnchorPeers:
      - Host: peer#[PeerNum].#[OrgName].#[Endpoint]
        Port: 7051
    Policies:
      Readers:
        Type: Signature
        Rule: "OR('#[OrgName]MSP.admin', '#[OrgName]MSP.peer', '#[OrgName]MSP.client')"
      Writers:
        Type: Signature
        Rule: "OR('#[OrgName]MSP.admin', '#[OrgName]MSP.client')"
      Admins:
        Type: Signature
        Rule: "OR('#[OrgName]MSP.admin')"
      Endorsement:
        Type: Signature
        Rule: "OR('#[OrgName]MSP.peer')"

Capabilities:
  Channel: &ChannelCapabilities
    V2_0: true
  Orderer: &OrdererCapabilities
    V2_0: true
  Application: &ApplicationCapabilities
    V2_0: true

Application: &ApplicationDefaults
  Organizations:
  Policies:
    Readers:
      Type: ImplicitMeta
      Rule: "ANY Readers"
    Writers:
      Type: ImplicitMeta
      Rule: "ANY Writers"
    Admins:
      Type: ImplicitMeta
      Rule: "MAJORITY Admins"
    LifecycleEndorsement:
      Type: ImplicitMeta
      Rule: "ANY Endorsement"
    Endorsement:
      Type: ImplicitMeta
      Rule: "ANY Endorsement"
  Capabilities:
    <<: *ApplicationCapabilities

Channel: &ChannelDefaults
  Policies:
    Readers:
      Type: ImplicitMeta
      Rule: "ANY Readers"
    Writers:
      Type: ImplicitMeta
      Rule: "ANY Writers"
    Admins:
      Type: ImplicitMeta
      Rule: "MAJORITY Admins"
  Capabilities:
    <<: *ChannelCapabilities

Orderer: &OrdererDefaults
  OrdererType: solo
  Addresses:
    - orderer.#[Endpoint]:7050
  BatchTimeout: #[BatchTimeout]
  BatchSize:
    MaxMessageCount: #[MaxMessageCount]
    AbsoluteMaxBytes: #[AbsoluteMaxBytes]
    PreferredMaxBytes: #[PreferredMaxBytes]
  Kafka:
    Brokers:
      - 127.0.0.1:9092
  Organizations:
  Policies:
    Readers:
      Type: ImplicitMeta
      Rule: "ANY Readers"
    Writers:
      Type: ImplicitMeta
      Rule: "ANY Writers"
    Admins:
      Type: ImplicitMeta
      Rule: "MAJORITY Admins"
    BlockValidation:
      Type: ImplicitMeta
      Rule: "ANY Writers"
  Capabilities:
    <<: *OrdererCapabilities

Profiles:
  # TwoOrgsOrdererGenesis配置文件用于创建系统通道创世块
  TwoOrgsOrdererGenesis:
    <<: *ChannelDefaults
    # 定义排序服务
    Orderer:
      <<: *OrdererDefaults
      # 定义排序服务的管理员
      Organizations:
        - *OrdererOrg
    # 创建一个名为SampleConsortium的联盟,包含两个组织Org,#[OrgName]
    Consortiums:
      SampleConsortium:
        Organizations:
          - *Org

  Channel:
    <<: *ChannelDefaults
    Consortium: SampleConsortium
    Application:
      <<: *ApplicationDefaults
      Organizations:
        - *Org`
	cryptoConfig = `OrdererOrgs:
  - Name: Orderer
    Domin: #[Endpoint]
    Specs:
      - Hostname: orderer

PeerOrgs:
  - Name: #[OrgName]
    Domain: #[OrgName].#[Endpoint]
    EnableNodeOUs: true
    Template:
      Count: #[PeerCount]
    Users:
      Count: 1`
	configtxgen = "docker run --privileged=true -v %s:/data -e FABRIC_CFG_PATH=/data --rm hyperledger/fabric-tools:2.4 configtxgen "
	cryptogen   = "docker run --rm --privileged=true -v %s:/data hyperledger/fabric-tools:2.4 cryptogen "
	compose     = `version: '3.7'

networks:
  byfn:

services:
  orderer.#[Endpoint]:
    container_name: orderer.#[Endpoint]
    image: hyperledger/fabric-orderer:2.4
    restart: always
    privileged: true
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/etc/hyperledger/configtx/genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/msp/orderer/msp
    volumes:
      - ./configtx/genesis.block:/etc/hyperledger/configtx/genesis.block
      - ./crypto-config/ordererOrganizations/orderers/orderer./msp:/etc/hyperledger/msp/orderer/msp
      - ./crypto-config/ordererOrganizations/orderers/orderer./tls:/etc/hyperledger/msp/orderer/tls
      - ./var/orderer.#[Endpoint]:/var/hyperledger/production/orderer
    ports:
      - "#[OrdererPort]:7050"
    deploy:
      resources:
        limits:
          cpus: '#[cpu_limit]'
    command: orderer start
    logging:
      driver: none
    networks:
      - byfn

  cli:
    container_name: cli.#[Endpoint]
    image: hyperledger/fabric-tools:2.4
    privileged: true
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - FABRIC_LOGGING_SPEC=DEBUG
      - CORE_PEER_ID=cli
      - CORE_CHAINCODE_KEEPALIVE=10
    command: /bin/bash
    volumes:
      - ./chaincode/#[CCName]:/opt/gopath/src/github.com/#[CCName]
      - ./crypto-config:/opt/crypto-config
      - ./configtx:/opt/configtx
    logging:
      driver: none
    networks:
      - byfn
`
	peer = `  peer#[PeerNum].#[OrgName].#[Endpoint]:
    container_name: peer#[PeerNum].#[OrgName].#[Endpoint]
    image: hyperledger/fabric-peer:2.4
    privileged: true
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=#[OrgName]_byfn
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/msp

      - CORE_PEER_ID=peer#[PeerNum].#[OrgName].#[Endpoint]
      - CORE_PEER_ADDRESS=peer#[PeerNum].#[OrgName].#[Endpoint]:7051
      - CORE_PEER_LOCALMSPID=#[OrgName]MSP
      - CORE_PEER_CHAINCODELISTENADDRESS=peer#[PeerNum].#[OrgName].#[Endpoint]:7052
      
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.#[OrgName].#[Endpoint]:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer#[PeerNum].#[OrgName].#[Endpoint]:7051
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_STATE_ENABLED=true

    command: peer node start
    volumes:
      - /var/run/docker.sock:/host/var/run/docker.sock
      - ./crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/peers/peer#[PeerNum].#[OrgName].#[Endpoint]/msp:/etc/hyperledger/peer/msp
      - ./crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/peers/peer#[PeerNum].#[OrgName].#[Endpoint]/tls:/etc/hyperledger/peer/tls
      - ./var/peer#[PeerNum].#[OrgName].#[Endpoint]:/var/hyperledger/production
      #- ./var/etc/hyperledger/fabric:/etc/hyperledger/fabric
      - ./var/hyperledger:/var/hyperledger
    ports:
      - "#[PeerPort]:7051"
      - "#[PeerPort2]:7053"
    deploy:
      resources:
        limits:
          cpus: '#[cpu_limit]'
    logging:
      driver: none
    networks:
      - byfn

`
	sdkConfig = `name: "fabric-sdk-app-network"

version: 1.0.0

client:
  organization: #[OrgName]
  logging:
    level: info
  peer:
    timeout:
      connection: 10s
      response: 180s
      discovery:
        greylistExpiry: 10s
  cryptoconfig:
    path: #[workingDir]/crypto-config
  credentialStore:
    path: "/tmp/state-store"
    cryptoStore:
      path: /tmp/msp

  BCCSP:
    security:
      enabled: true
      default:
        provider: "SW"
      hashAlgorithm: "SHA2"
      softVerify: true
      level: 256

  tlsCerts:
    systemCertPool: true
    client:
      key:
        path: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/tls/client.key
      cert:
        path: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/tls/client.crt

peerPermission: &peerPermission
  endorsingPeer: true
  chaincodeQuery: true
  ledgerQuery: true
  eventSource: true

channelPolicy: &channelPolicy
  discovery:
    maxTargets: 2
    retryOpts:
      attempts: 4
      initialBackoff: 500ms
      maxBackoff: 5s
      backoffFactor: 2.0

  selection:
    SortingStrategy: BlockHeightPriority
    Balancer: RoundRobin
    BlockHeightLagThreshold: 5

  queryChannelConfig:
    minResponses: 1
    maxTargets: 1
    retryOpts:
      attempts: 5
      initialBackoff: 500ms
      maxBackoff: 5s
      backoffFactor: 2.0

  eventService:
    resolverStrategy: PreferOrg
    balancer: RoundRobin
    blockHeightLagThreshold: 2
    reconnectBlockHeightLagThreshold: 5
    peerMonitorPeriod: 3s


grpcOptions: &grpcOptions
  keep-alive-time: 5s
  keep-alive-timeout: 6s
  keep-alive-permit: false
  fail-fast: true
  allow-insecure: true

channels:
  hello:
    peers:
      peer#[PeerNum].#[OrgName].#[Endpoint]:
        <<: *peerPermission
    policies:
      <<: *channelPolicy
    configpath: #[workingDir]/configtx/hello.tx

organizations:
  #[OrgName]:
    mspID: #[OrgName]MSP
    users:
      User1:
        cert:
          path: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/msp/cacerts/ca.#[OrgName].#[Endpoint]-cert.pem
    cryptoPath: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/Admin@#[OrgName].#[Endpoint]/msp
    peers:
      - peer#[PeerNum].#[OrgName].#[Endpoint]

  ordererorg:
    mspID: OrdererMSP
    users:
      Admin:
        cert:
          path: #[workingDir]/crypto-config/ordererOrganizations/users/Admin@/msp/cacerts/ca.-cert.pem
    cryptoPath: #[workingDir]/crypto-config/ordererOrganizations/users/Admin@/msp

orderers:
  orderer.#[Endpoint]:
    url: orderer.#[Endpoint]:#[OrdererPort]
    grpcOptions:
      ssl-target-name-override: orderer.#[Endpoint]
      hostnameOverride: orderer.#[Endpoint]
      <<: *grpcOptions
    tlsCACerts:
      path: #[workingDir]/crypto-config/ordererOrganizations/tlsca/tlsca.-cert.pem

peers:
  peer#[PeerNum].#[OrgName].#[Endpoint]:
    url: peer#[PeerNum].#[OrgName].#[Endpoint]:#[PeerPort]
    eventUrl: peer#[PeerNum].#[OrgName].#[Endpoint]:#[PeerEventPort]
    grpcOptions:
      ssl-target-name-override: peer#[PeerNum].#[OrgName].#[Endpoint]
      hostnameOverride: peer#[PeerNum].#[OrgName].#[Endpoint]
      <<: *grpcOptions
    tlsCACerts:
      path: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/tlsca/tlsca.#[OrgName].#[Endpoint]-cert.pem
entityMatchers:
  orderer:
    - pattern: orderer.#[Endpoint]:(\d+)
      urlSubstitutionExp: orderer.#[Endpoint]:#[OrdererPort]
      sslTargetOverrideUrlSubstitutionExp: orderer.#[Endpoint]
      mappedHost: orderer.#[Endpoint]

  peer:
    - pattern: (\w+).#[OrgName].#[Endpoint]:(\d+)
      urlSubstitutionExp: ${1}.#[OrgName].#[Endpoint]:#[PeerPort]
      sslTargetOverrideUrlSubstitutionExp: ${1}.#[OrgName].#[Endpoint]
      mappedHost: peer#[PeerNum].#[OrgName].#[Endpoint]

`
	tapeConfig = `peer#[PeerNum]: &peer#[PeerNum]
  addr: 127.0.0.1:#[PeerPort]
orderer1: &orderer1
  addr: 127.0.0.1:#[OrdererPort]

# Nodes to interact with
endorsers:
  - *peer0
# we might support multi-committer in the future for more complex test scenario,
# i.e. consider tx committed only if it's done on >50% of nodes. But for now,
# it seems sufficient to support single committer.
committers:
  - *peer0
commitThreshold: 1

orderer: *orderer1


# Invocation configs
channel: hello
chaincode: cc1
args:
  - transfer
  - peer0.#[OrgName].#[Endpoint]
  - peer1.#[OrgName].#[Endpoint]
  - 0x1
mspid: #[OrgName]MSP
private_key: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/msp/keystore/priv_sk
sign_cert: #[workingDir]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/msp/signcerts/User1@#[OrgName].#[Endpoint]-cert.pem
num_of_conn: 10
client_per_conn: 10
`
	caliperCompose = `version: "2"

networks:
  caliper.com:
    external:
      name: #[OrgName]_byfn

services:
  caliper:
    container_name: caliper
    image: hyperledger/caliper:0.4.2
    command: launch manager --caliper-fabric-gateway-enabled
    environment:
      - CALIPER_BIND_SUT=fabric:2.2
      - CALIPER_BENCHCONFIG=#[OrgName]/caliper-workspace/benchmarks/config.yaml
      - CALIPER_NETWORKCONFIG=#[OrgName]/caliper-workspace/networks/network-config.yaml
    volumes:
      - ../../:/hyperledger/caliper/workspace
    networks:
        - caliper.com
`
	caliperConfig = `---
test:
  name: simple
  description:
    This is an example benchmark for caliper, to test the backend DLT's
    performance with simple account opening & querying transactions
  workers:
    type: local
    number: 4
  rounds:
    - label: open
      description: Test description for the opening of an account through the deployed chaincode
      txNumber: 3000
      rateControl:
        type: maximum-rate
        opts:
          tps: 300
          step: 50
          sampleInterval: 1
          includeFailed: true
      arguments:
        money: 100
      workload:
        module: #[OrgName]/caliper-workspace/benchmarks/open.js
#监测docker容器的资源使用情况
monitor:
  type:
  - docker
  docker:  
    containers:
    - all
    charting:
    bar:
      metrics: [Memory(avg), CPU%(avg)]
    polar:
      metrics: [all]`
	openjs = `'use strict';

const {WorkloadModuleInterface} = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleInterface {
    constructor() {
        super();
        this.workerIndex = -1;
        this.totalWorkers = -1;
        this.roundIndex = -1;
        this.roundArguments = undefined;
        this.sutAdapter = undefined;
        this.sutContext = undefined;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        this.workerIndex = workerIndex;
        this.totalWorkers = totalWorkers;
        this.roundIndex = roundIndex;
        this.roundArguments = roundArguments;
        this.sutAdapter = sutAdapter;
        this.sutContext = sutContext;
    }

    async submitTransaction() {
        // TX arguments for "cc1"
        let requestSettings = {
            contractId: 'cc1',
            contractFunction: 'transfer',
            contractArguments: ["peer0.#[OrgName].#[Endpoint]", "peer1.#[OrgName].#[Endpoint]", "0x1"],
            //readOnly: true,
            //invokerIdentity: 'client0.org2.example.com',
            timeout: 10
        };

        await this.sutAdapter.sendRequests(requestSettings);
    }

    async cleanupWorkloadModule() {
        // NOOP
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;`
	caliperNetworkConfig = `name: Fabric
version: "1.0"
mutual-tls: false

caliper:
  blockchain: fabric


info:
  Version: 2.1
clients:
  client.#[OrgName].#[Endpoint]:
    client:
      organization: #[OrgName]
      credentialStore:
        path: /tmp/hfc-kvs/#[OrgName]
        cryptoStore:
          path: /tmp/hfc-cvs/#[OrgName]
      clientPrivateKey:
        path: #[OrgName]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/msp/keystore/priv_sk
      clientSignedCert:
        path: #[OrgName]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/User1@#[OrgName].#[Endpoint]/msp/signcerts/User1@#[OrgName].#[Endpoint]-cert.pem


channels:
  hello:
    created: true
    orderers:
      - orderer.#[Endpoint]
    peers:
      peer#[PeerNum].#[OrgName].#[Endpoint]:
        eventSource: true

    contracts:
      - id: cc1
        version: "1.0"
        language: golang
        path: chaincode/cc1

organizations:
  #[OrgName]:
    mspid: #[OrgName]MSP
    peers:
      - peer#[PeerNum].#[OrgName].#[Endpoint]
    adminPrivateKey:
      path: #[OrgName]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/Admin@#[OrgName].#[Endpoint]/msp/keystore/priv_sk
    signedCert:
      path: #[OrgName]/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/Admin@#[OrgName].#[Endpoint]/msp/signcerts/Admin@#[OrgName].#[Endpoint]-cert.pem



orderers:
  orderer.#[Endpoint]:
    url: grpc://orderer.#[Endpoint]:7050
    grpcOptions:
      ssl-target-name-override: orderer.#[Endpoint]

peers:
  peer#[PeerNum].#[OrgName].#[Endpoint]:
    url: grpc://peer#[PeerNum].#[OrgName].#[Endpoint]:7051
    grpcOptions:
      ssl-target-name-override: peer#[PeerNum].#[OrgName].#[Endpoint]
      grpc.keepalive_time_ms: 600000

`
	dockerPeerCmdPrefix = "docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/users/Admin@#[OrgName].#[Endpoint]/msp -e CORE_PEER_LOCALMSPID=#[OrgName]MSP -e CORE_PEER_ADDRESS=peer#[PeerNum].#[OrgName].#[Endpoint]:7051 cli.#[Endpoint] peer "
	createChannel       = dockerPeerCmdPrefix + "channel create -o orderer.#[Endpoint]:7050 -c hello -f /opt/configtx/hello.tx"
	joinChannel         = dockerPeerCmdPrefix + "channel join -b /opt/configtx/hello.block"
	setAnchorPeer       = dockerPeerCmdPrefix + "channel update -o orderer.#[Endpoint]:7050 -c hello -f /opt/configtx/#[OrgNameUpper]MSPanchors_hello.tx"
	installCC           = dockerPeerCmdPrefix + "lifecycle chaincode install /opt/gopath/src/github.com/#[CCName]/#[CCName].tar.gz"
	approveCC           = dockerPeerCmdPrefix + "lifecycle chaincode approveformyorg --channelID hello --name #[CCName] --version 1.0 --init-required --package-id #[packageID] --sequence 1"
	commitCC            = dockerPeerCmdPrefix + "lifecycle chaincode commit -o orderer.#[Endpoint]:7050 --channelID hello --name #[CCName] --version 1.0 --sequence 1 --init-required"
	initCC              = dockerPeerCmdPrefix + "chaincode invoke -o orderer.#[Endpoint]:7050 --isInit -C hello -n #[CCName] -c '{\"Args\":[\"peer0.#[OrgName].#[Endpoint]\", \"0xffff\",\"peer1.#[OrgName].#[Endpoint]\", \"0xffff\"]}'"
	//invokeCC            = dockerPeerCmdPrefix + "chaincode invoke -o orderer.#[Endpoint]:7050 -C hello -n #[CCName] -c '{\"Args\":#[Args]}'"
)
