package ker

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
	"unsafe"

	"code.byted.org/gopkg/consul"
	"code.byted.org/gopkg/env"
)

var (
	counter             = -1 // 循环计数
	requestAssist       *RequestInfoAssist
	kerServiceHostAddrs []string
)

const (
	sdkType = "ker_dec_sdk"
)

var moduleWithPsmList = map[string][]string{
	"vod": {
		"tiktok.vod.ker",
		"toutiao.videoarch.smart_player",
	},
	"image": {
		"tiktok.vod.ker",
	},
}

type SDKFetchRuleGroupRequest struct {
	AddrIpv6    string `json:"addr_ipv6"`
	AddrIpv4    string `json:"addr_ipv4"`
	PodName     string `json:"pod_name"`
	SdkType     string `json:"sdk_type"`
	SdkHostPsm  string `json:"sdk_host_psm"`
	ModuleName  string `json:"module_name"`
	ClusterName string `json:"cluster_name"`
}

type RequestInfoAssist struct {
	RequestList []*SDKFetchRuleGroupRequest
}

func fetchKerServiceHostAddr() {
	idc := consul.IDC(env.IDC())
	ipv6Endpoints, err := consulLookupIPV6(idc, "tiktok.vod.ker") // 解析ipv6
	if err != nil {
		panic(err)
	}

	if len(ipv6Endpoints.Addrs()) == 0 {
		return
	}

	addrs := ipv6Endpoints.Addrs()
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&(kerServiceHostAddrs))), *(*unsafe.Pointer)(unsafe.Pointer(&addrs))) // #nosec
}

func fetchSDKFetchRuleGroupRequestList() {
	ReuestList := newSDKFetchRuleGroupRequest()
	if len(ReuestList) == 0 {
		return
	}

	info := &RequestInfoAssist{}
	info.RequestList = append(info.RequestList, ReuestList...)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&(requestAssist))), *(*unsafe.Pointer)(unsafe.Pointer(&info))) // #nosec
}

func newSDKFetchRuleGroupRequest() []*SDKFetchRuleGroupRequest {
	idc := consul.IDC(env.IDC())
	var outputs []*SDKFetchRuleGroupRequest
	for module, psmList := range moduleWithPsmList {
		if len(module) == 0 || len(psmList) == 0 {
			continue
		}
		for _, psm := range psmList {
			ipv6Endpoints, err := consulLookupIPV6(idc, psm) // 解析ipv6
			if err != nil || len(ipv6Endpoints) == 0 {
				continue
			}
			ipv4Endpoints, err := consulLookupIPV4(idc, psm) // 解析ipv4
			if err != nil || len(ipv6Endpoints) == 0 {
				continue
			}
			for _, ipv6Endpoint := range ipv6Endpoints {
				_, port, err := net.SplitHostPort(ipv6Endpoint.Addr)
				if err != nil {
					continue
				}
				key1 := ipv6Endpoint.Cluster + ipv6Endpoint.Env + port
				for _, ipv4Endpoint := range ipv4Endpoints {
					_, port, err := net.SplitHostPort(ipv4Endpoint.Addr)
					if err != nil {
						continue
					}
					key2 := ipv4Endpoint.Cluster + ipv4Endpoint.Env + port
					if key1 == key2 {
						outputs = append(outputs, &SDKFetchRuleGroupRequest{
							SdkHostPsm:  psm,
							ModuleName:  module,
							SdkType:     sdkType,
							AddrIpv6:    ipv6Endpoint.Addr,
							AddrIpv4:    ipv4Endpoint.Addr,
							PodName:     ipv4Endpoint.Env,
							ClusterName: ipv4Endpoint.Cluster,
						})
					}
				}

			}
		}
	}

	return outputs
}

func consulLookupIPV6(idc consul.IDC, psm string) (consul.Endpoints, error) {
	endpoints, err := consul.LookupName(psm,
		consul.WithIDC(idc),
		consul.WithNoCache(true),
		consul.WithUnique(consul.V6),
		consul.WithNoConsulHash(true),
		consul.WithAddrFamily(consul.V6),
	)
	return endpoints, err
}

func consulLookupIPV4(idc consul.IDC, psm string) (consul.Endpoints, error) {
	endpoints, err := consul.LookupName(psm,
		consul.WithIDC(idc),
		consul.WithNoCache(true),
		consul.WithUnique(consul.V4),
		consul.WithNoConsulHash(true),
		consul.WithAddrFamily(consul.V4),
	)
	return endpoints, err
}

func serializeSDKFetchRuleGroupRequest2JSON(req *SDKFetchRuleGroupRequest) (string, error) {
	body, err := json.Marshal(req)
	return string(body), err
}

func getKerServiceHostAddr() (string, error) {
	if len(kerServiceHostAddrs) == 0 {
		return "", errors.New("ker hosts are empty")
	}
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
	index := rand.Intn(len(kerServiceHostAddrs))
	return kerServiceHostAddrs[index], nil
}

// getAddrForFetchRuleGroup构造请求地址
func getAddrForFetchRuleGroup(addr string) string {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("http://%s:%s/ker/api/sdk/fetch_rule_groups", ip, port)
}

// getSDKFetchRuleGroupRequest [随机地/顺序地] 获取request参数
func getSDKFetchRuleGroupRequest(randomly bool) *SDKFetchRuleGroupRequest {
	index := 0           // 记录随机访问的下标
	sleepTimeSecond := 0 // 默认sleep 1s
	maxIndex := len(requestAssist.RequestList)

	if randomly {
		index = rand.Intn(maxIndex) // 生成随机整数0~maxIndex-1
		sleepTimeSecond = index % 5 // 最多sleep 5s
	} else {
		counter = counter + 1
		index = counter % maxIndex // 顺序性获取参数
	}

	request := requestAssist.RequestList[index] // 获取一个服务实例的请求参数
	if request != nil && sleepTimeSecond > 0 {
		time.Sleep(time.Second * time.Duration(sleepTimeSecond))
	}
	return request
}
