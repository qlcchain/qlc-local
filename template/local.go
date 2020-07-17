package template

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/qlcchain/qlc-go-sdk/pkg/types"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//type Node struct {
//	cfgPath             string
//	Process             *os.Process
//	GenesisChainAccount *types.Account // genesis qlc block account
//	GenesisGasAccount   *types.Account // genesis gas block account
//	Account             *types.Account // node run with the account
//	Client              *qlcchain.QLCClient
//	ListenAddress       string
//	GRPCListenAddress   string
//	Port                string
//	para                *NodeParam
//	lock                *sync.Mutex
//	started             bool
//	logger              *zap.SugaredLogger
//}

// Create docker-compose.yml file by parameters
// Parameters:
// - dir: directory of docker-compose.yml to be created
// - nodeCount: number of qlc nodes to run
// - repCount: number of representation nodes to run
// - ptmCount: number of ptm nodes to run
// - qlcVersion: docker image version of qlc, if set "", default is "latest"
// - ptmVersion: docker image version of ptm, if set "", default is "latest"
func Template(dir string, nodeCount, repCount, ptmCount int, qlcVersion, ptmVersion string) ([]QlcNode, []PtmNode, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return nil, nil, err
	}

	ptmParams := initPTMs(ptmCount, ptmVersion)
	ptmUrls := make([]string, 0)
	for _, n := range ptmParams {
		ptmUrls = append(ptmUrls, n.cEndpoint)
	}
	qlcParams := initNodes(nodeCount, repCount, qlcVersion, ptmUrls)
	data := map[string]interface{}{
		"Nodes": qlcParams,
		"Ptms":  ptmParams,
	}
	writer := new(bytes.Buffer)
	err = tmpl.Execute(writer, data)
	if err != nil {
		return nil, nil, err
	}

	if err := ioutil.WriteFile(dir, writer.Bytes(), os.ModePerm); err != nil {
		return nil, nil, err
	}
	fmt.Println(writer.String())
	return toQLCNodes(qlcParams), toPTMNodes(ptmParams), nil
}

func initPTMs(count int, ptmVersion string) []PtmParam {
	params := make([]PtmParam, 0)
	if count == 0 {
		return params
	}
	nodeIndex := generateNum()
	port1Index := 9181
	port2Index := 9183
	ipv4Address := 21
	if ptmVersion == "" {
		ptmVersion = "latest"
	}
	image := fmt.Sprintf("qlcchain/ptm:%s", ptmVersion)
	containerName := fmt.Sprintf("ptm_node%d", nodeIndex)
	param1 := PtmParam{
		Name:          containerName,
		Image:         image,
		ContainerName: containerName,
		Command:       commond(true, nodeIndex, port1Index, port2Index, ipv4Address),
		Port1:         fmt.Sprintf("%d:%d", port1Index, port1Index),
		Port2:         fmt.Sprintf("%d:%d", port2Index, port2Index),
		cEndpoint:     fmt.Sprintf("http://10.0.0.%d:%d", ipv4Address, port2Index),
		Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address),
		Volumes:       fmt.Sprintf("./qlcptm-%d:/home/qlcptm%d/", nodeIndex, nodeIndex),
	}
	params = append(params, param1)
	for i := 1; i < count; i++ {
		port1Indext := port1Index + 100*i
		port2Indext := port2Index + 100*i
		containerName := fmt.Sprintf("ptm_node%d", nodeIndex+i)
		paramt := PtmParam{
			Name:          containerName,
			Image:         image,
			ContainerName: containerName,
			Command:       commond(false, nodeIndex+i, port1Indext, port2Indext, ipv4Address+i),
			Port1:         fmt.Sprintf("%d:%d", port1Indext, port1Indext),
			Port2:         fmt.Sprintf("%d:%d", port2Indext, port2Indext),
			cEndpoint:     fmt.Sprintf("http://10.0.0.%d:%d", ipv4Address+i, port2Indext),
			Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address+i),
			Volumes:       fmt.Sprintf("./qlcptm-%d:/home/qlcptm%d/", nodeIndex+i, nodeIndex+i),
		}
		params = append(params, paramt)
	}
	return params
}

func commond(first bool, index int, port1Index, port2Index, ipv4Address int) string {
	coms := make([]string, 0)
	javaStr := "java -jar /tessera/tessera-app.jar"
	coms = append(coms, javaStr)
	confile := fmt.Sprintf("-configfile /home/qlcptm%d/ptm.json", index)
	coms = append(coms, confile)
	sc := fmt.Sprintf("-o serverConfigs[0].serverAddress=http://127.0.0.1:%d", port1Index)
	coms = append(coms, sc)
	scb := fmt.Sprintf("-o serverConfigs[0].bindingAddress=http://127.0.0.1:%d", port1Index)
	coms = append(coms, scb)
	sc2 := fmt.Sprintf("-o serverConfigs[2].serverAddress=http://10.0.0.%d:%d", ipv4Address, port2Index)
	coms = append(coms, sc2)
	sc2b := fmt.Sprintf("-o serverConfigs[2].bindingAddress=http://0.0.0.0:%d", port2Index)
	coms = append(coms, sc2b)
	if first {
		url := fmt.Sprintf("-o peer[0].url=http://10.0.0.%d:%d", ipv4Address+1, port2Index+100)
		coms = append(coms, url)
	} else {
		url := fmt.Sprintf("-o peer[0].url=http://10.0.0.%d:%d", ipv4Address-1, port2Index-100)
		coms = append(coms, url)
	}
	return "[\"" + strings.Join(coms, "\",\"") + "\"]"
}

func initNodes(count, repCount int, qlcVesion string, ptmUrls []string) []NodeParam {
	params := make([]NodeParam, 0)
	if count == 0 {
		return params
	}
	nodeIndex := generateNum()
	portP2PIndex := 19034
	portHttpIndex := 19035
	portWsIndex := 19036
	portBootNodeIndex := 19037
	ipv4Address := 10
	if qlcVesion == "" {
		qlcVesion = "latest"
	}
	ptmUrl := ""
	if len(ptmUrls) > 0 {
		ptmUrl = ptmUrls[0]
	}
	seed := "46b31acd0a3bf072e7bea611a86074e7afae5ff95610f5f870208f2fd9357418"
	image := fmt.Sprintf("qlcchain/go-qlc-test:%s", qlcVesion)
	bootparam := NodeParam{
		Name:          fmt.Sprintf("qlcchain_node%d", nodeIndex),
		Image:         image,
		ContainerName: fmt.Sprintf("qlcchain_node%d", nodeIndex),
		Command:       command(true, portBootNodeIndex, portP2PIndex, ptmUrl, seed),
		PortP2P:       portStr(portP2PIndex, 19734),
		PortHttp:      portStr(portHttpIndex, 19735),
		PortWs:        portStr(portWsIndex, 19736),
		Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address),
		Volumes:       fmt.Sprintf("./go-qlc-test-%d:/qlcchain/.gqlcchain_test", nodeIndex),
		accountSeed:   seed,
	}
	params = append(params, bootparam)
	if count > 1 {
		for i := 1; i < count; i++ {
			ptmUrlt := ""
			if len(ptmUrls) > i {
				ptmUrlt = ptmUrls[i]
			}
			seed := ""
			if repCount > i {
				seedNew, _ := types.NewSeed()
				seed = hex.EncodeToString(seedNew[:])
			}
			paramt := NodeParam{
				Name:          fmt.Sprintf("qlcchain_node%d", nodeIndex+i),
				Image:         image,
				ContainerName: fmt.Sprintf("qlcchain_node%d", nodeIndex+i),
				Command:       command(false, portBootNodeIndex+100*i, portP2PIndex+100*i, ptmUrlt, seed),
				PortP2P:       portStr(portP2PIndex+100*i, 19734),
				PortHttp:      portStr(portHttpIndex+100*i, 19735),
				PortWs:        portStr(portWsIndex+100*i, 19736),
				Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address+i),
				Volumes:       fmt.Sprintf("./go-qlc-test-%d:/root/.gqlcchain_test", nodeIndex+i),
				accountSeed:   seed,
			}
			params = append(params, paramt)
		}
	}
	return params
}

func portStr(port int, origin int) string {
	return fmt.Sprintf("%s:%s", strconv.Itoa(port), strconv.Itoa(origin))
}

func command(isBootNode bool, portBootNode, p2pPort int, ptmUrl string, seed string) string {
	coms := make([]string, 0)
	key, id, _ := identityConfig()
	configParams := "--configParams=logLevel=info;rpc.rpcEnabled=true;p2p.isBootNode=%s;p2p.bootNode=['http://127.0.0.1:%d/bootNode'];p2p.bootNodeHttpServer=0.0.0.0:%d;p2p.listen=/ip4/0.0.0.0/tcp/%d;p2p.identity.peerId=%s;p2p.identity.privateKey=%s"
	if ptmUrl == "" {
		configParams := fmt.Sprintf(configParams, strconv.FormatBool(isBootNode), portBootNode, portBootNode, p2pPort, id, key)
		coms = append(coms, configParams)
	} else {
		configParams = configParams + ";privacy.enable=true;privacy.ptmNode=%s"
		configParams := fmt.Sprintf(configParams, strconv.FormatBool(isBootNode), portBootNode, portBootNode, p2pPort, id, key, ptmUrl)
		coms = append(coms, configParams)
	}
	if seed != "" {
		seed := fmt.Sprintf("--seed=%s", seed)
		coms = append(coms, seed)
	}
	return "[\"" + strings.Join(coms, "\",\"") + "\"]"
}

// identityConfig initializes a new identity.
func identityConfig() (string, string, error) {
	sk, pk, err := ic.GenerateKeyPair(ic.RSA, 2048)
	if err != nil {
		return "", "", err
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := sk.Bytes()
	if err != nil {
		return "", "", err
	}
	privKey := base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return "", "", err
	}
	peerID := id.Pretty()
	return privKey, peerID, nil
}

func generateNum() int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(100)
	return randNum
}

func toQLCNodes(params []NodeParam) []QlcNode {
	nodes := make([]QlcNode, 0)
	for _, param := range params {
		account := new(types.Account)
		if param.accountSeed != "" {
			sByte, _ := hex.DecodeString(param.accountSeed)
			seed, _ := types.BytesToSeed(sByte)
			account, _ = seed.Account(0)
		} else {
			account = nil
		}
		node := QlcNode{
			HTTPEndpoint:  fmt.Sprintf("http://127.0.0.1:%s", getPort(param.PortHttp)),
			WSEndpoint:    fmt.Sprintf("http://127.0.0.1:%s", getPort(param.PortWs)),
			ContainerName: param.ContainerName,
			Account:       account,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func toPTMNodes(params []PtmParam) []PtmNode {
	nodes := make([]PtmNode, 0)
	for _, param := range params {
		node := PtmNode{
			EndPoint:      fmt.Sprintf("http://127.0.0.1:%s", getPort(param.Port2)),
			ContainerName: param.ContainerName,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func getPort(value string) string {
	vs := strings.Split(value, ":")
	return vs[0]
}
