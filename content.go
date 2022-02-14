package fabricGen

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
    Name: #[OrgName]MSP
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
  BatchTimeout: 0.2s
  BatchSize:
    MaxMessageCount: 10
    AbsoluteMaxBytes: 99MB
    PreferredMaxBytes: 512KB
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
    # 创建一个名为SampleConsortium的联盟,包含两个组织Org,Org2
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
	configtxgen = "docker run -v %s:/data -e FABRIC_CFG_PATH=/data --rm hyperledger/fabric-tools:2.4 configtxgen "
	cryptogen   = "docker run --rm -v %s:/data hyperledger/fabric-tools:2.4 cryptogen "
	compose     = `version: '3'

networks:
  byfn:

services:
  orderer.#[Endpoint]:
    container_name: orderer.#[Endpoint]
    image: hyperledger/fabric-orderer:2.4
    restart: always
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
    command: orderer start
    networks:
      - byfn

  cli:
    container_name: cli.#[Endpoint]
    image: hyperledger/fabric-tools:2.4
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
    networks:
      - byfn

`
	peer = `  peer#[PeerNum].#[OrgName].#[Endpoint]:
    container_name: peer#[PeerNum].#[OrgName].#[Endpoint]
    image: hyperledger/fabric-peer:2.4
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
    command: peer node start
    volumes:
      - /var/run/docker.sock:/host/var/run/docker.sock
      - ./crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/peers/peer#[PeerNum].#[OrgName].#[Endpoint]/msp:/etc/hyperledger/peer/msp
      - ./crypto-config/peerOrganizations/#[OrgName].#[Endpoint]/peers/peer#[PeerNum].#[OrgName].#[Endpoint]/tls:/etc/hyperledger/peer/tls
      - ./var/peer#[PeerNum].#[OrgName].#[Endpoint]:/var/hyperledger/production
    ports:
      - "#[PeerPort]:7051"
      - "#[PeerPort2]:7053"
    networks:
      - byfn

`
)
