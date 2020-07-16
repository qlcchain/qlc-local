package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
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

func Template(dir string, nodeCount, repCount, ptmCount int) ([]QlcNode, []PtmNode, error) {
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		return nil, nil, err
	}

	qlcParams, qlcNodes := initNodes(nodeCount)
	ptmParams, ptmNodes := initPTMs(ptmCount)
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
	return qlcNodes, ptmNodes, nil
}

func initPTMs(count int) ([]PTMParam, []PtmNode) {
	params := make([]PTMParam, 0)
	nodes := make([]PtmNode, 0)
	if count == 0 {
		return params, nodes
	}
	nodeIndex := 1
	port1Index := 9181
	port2Index := 9183
	ipv4Address := 20
	image := "qlcchain/ptm:0.10.5"
	param1 := PTMParam{
		Name:          fmt.Sprintf("ptm_node%d", nodeIndex),
		Image:         image,
		ContainerName: fmt.Sprintf("ptm_node%d", nodeIndex),
		Command:       commond(true, nodeIndex, port1Index, port2Index, ipv4Address),
		Port1:         fmt.Sprintf("%d:%d", port1Index, port1Index),
		Port2:         fmt.Sprintf("%d:%d", port2Index, port2Index),
		Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address),
		Volumes:       fmt.Sprintf("./qlcptm-%d:/home/qlcptm%d/", nodeIndex, nodeIndex),
	}
	params = append(params, param1)
	node1 := PtmNode{
		Url: fmt.Sprintf("http://127.0.0.1:%d", port2Index),
	}
	nodes = append(nodes, node1)
	for i := 1; i < count; i++ {
		port1Indext := port1Index + 100*i
		port2Indext := port2Index + 100*i
		paramt := PTMParam{
			Name:          fmt.Sprintf("ptm_node%d", nodeIndex+i),
			Image:         image,
			ContainerName: fmt.Sprintf("ptm_node%d", nodeIndex+i),
			Command:       commond(false, nodeIndex+i, port1Indext, port2Indext, ipv4Address+i),
			Port1:         fmt.Sprintf("%d:%d", port1Indext, port1Indext),
			Port2:         fmt.Sprintf("%d:%d", port2Indext, port2Indext),
			Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address+i),
			Volumes:       fmt.Sprintf("./qlcptm-%d:/home/qlcptm%d/", nodeIndex+i, nodeIndex+i),
		}
		params = append(params, paramt)
		nodet := PtmNode{
			Url: fmt.Sprintf("http://127.0.0.1:%d", port2Indext),
		}
		nodes = append(nodes, nodet)
	}
	return params, nodes
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

func initNodes(count int) ([]NodeParam, []QlcNode) {
	params := make([]NodeParam, 0)
	nodes := make([]QlcNode, 0)
	if count == 0 {
		return params, nodes
	}
	nodeIndex := 1
	portP2PIndex := 19034
	portHttpIndex := 19035
	portWsIndex := 19036
	portBootNodeIndex := 19037
	ipv4Address := 10
	image := "qlcchain/go-qlc-test:latest"
	bootparam := NodeParam{
		Name:          fmt.Sprintf("qlcchain_node%d", nodeIndex),
		Image:         image,
		ContainerName: fmt.Sprintf("qlcchain_node%d", nodeIndex),
		Command:       command(true, portBootNodeIndex, portP2PIndex),
		PortP2P:       portStr(portP2PIndex, 19734),
		PortHttp:      portStr(portHttpIndex, 19735),
		PortWs:        portStr(portWsIndex, 19736),
		Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address),
		Volumes:       fmt.Sprintf("./go-qlc-test-%d:/qlcchain/.gqlcchain_test", nodeIndex),
	}
	params = append(params, bootparam)
	if count > 1 {
		for i := 1; i < count; i++ {
			paramt := NodeParam{
				Name:          fmt.Sprintf("qlcchain_node%d", nodeIndex+i),
				Image:         image,
				ContainerName: fmt.Sprintf("qlcchain_node%d", nodeIndex+i),
				Command:       command(false, portBootNodeIndex+100*i, portP2PIndex+100*i),
				PortP2P:       portStr(portP2PIndex+100*i, 19734),
				PortHttp:      portStr(portHttpIndex+100*i, 19735),
				PortWs:        portStr(portWsIndex+100*i, 19736),
				Ipv4Address:   fmt.Sprintf("10.0.0.%d", ipv4Address+i),
				Volumes:       fmt.Sprintf("./go-qlc-test-%d:/root/.gqlcchain_test", nodeIndex+i),
			}
			params = append(params, paramt)
		}
	}
	return params, nodes
}

func portStr(port int, origin int) string {
	return fmt.Sprintf("%s:%s", strconv.Itoa(port), strconv.Itoa(origin))
}

func command(isBootNode bool, portBootNode, p2pPort int) string {
	coms := make([]string, 0)
	key, id, _ := identityConfig()
	configParams := fmt.Sprintf("--configParams=logLevel=info;rpc.rpcEnabled=true;p2p.isBootNode=%s;p2p.bootNode=['http://127.0.0.1:%d/bootNode'];p2p.bootNodeHttpServer=0.0.0.0:%d;p2p.listen=/ip4/0.0.0.0/tcp/%d;p2p.identity.peerId=%s;p2p.identity.privateKey=%s",
		strconv.FormatBool(isBootNode), portBootNode, portBootNode, p2pPort, id, key)
	coms = append(coms, configParams)
	if isBootNode {
		seed := "--seed=46b31acd0a3bf072e7bea611a86074e7afae5ff95610f5f870208f2fd9357418"
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
