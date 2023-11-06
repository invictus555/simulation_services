package ker

import (
	"fmt"
	"time"

	"simulation_services/utils"
)

func init() {
	fetchSDKFetchRuleGroupRequestList()
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			<-ticker.C
			fetchSDKFetchRuleGroupRequestList()
		}
	}()
}

// SDKFetchKerRuleGroupSimulationService 模拟调用ker服务[机房粒度]
func SDKFetchKerRuleGroupSimulationService(randomly bool) {
	for {
		if requestAssist == nil || len(requestAssist.RequestList) == 0 {
			fmt.Println("request parameters is empty")
			break
		}

		addr := getKerServiceHostAddr()
		if val, ok := kerServiceStatus[addr]; ok && !val { // 无效ker服务
			continue
		}

		url := getAddrForFetchRuleGroup(addr)
		request := getSDKFetchRuleGroupRequest(randomly)
		body, err := serializeSDKFetchRuleGroupRequest2JSON(request) // 序列化SDK Fetch RuleGroup的请求参数
		if err != nil || len(body) == 0 {
			continue
		}

		response, err := utils.DoHttpGetMethodV2(url, []byte(body))
		if err != nil || len(response) == 0 {
			fmt.Println("Do http POST failed, err:", err)
		} else {
			fmt.Printf("Successfully:randomly:%v, PSM:%-32s,module:%-6s,IP:%s->%s\n", randomly, request.SdkHostPsm, request.ModuleName, addr, request.AddrIpv6)
		}
		kerServiceStatus[addr] = err == nil // 记录ker服务是否有效
	}
}
