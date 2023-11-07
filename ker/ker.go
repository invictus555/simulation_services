package ker

import (
	"fmt"
	"time"

	"simulation_services/utils"
)

func init() {
	fetchKerServiceHostAddr()
	fetchSDKFetchRuleGroupRequestList()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			<-ticker.C
			fetchKerServiceHostAddr()
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

		addr, err := getKerServiceHostAddr()
		if err != nil || len(addr) == 0 {
			continue
		}

		request := getSDKFetchRuleGroupRequest(randomly)
		body, err := serializeSDKFetchRuleGroupRequest2JSON(request) // 序列化SDK Fetch RuleGroup的请求参数
		if err != nil || len(body) == 0 {
			continue
		}

		response, err := utils.DoHttpGetMethodV2(getAddrForFetchRuleGroup(addr), []byte(body))
		if err != nil || len(response) == 0 {
			fmt.Println("Do http POST failed, err:", err)
		} else {
			fmt.Printf("Successfully[%s]:randomly:%v PSM:%-32s module:%-6s IP:%-48s--->\t%-48s\n", utils.GetNowTime(), randomly, request.SdkHostPsm, request.ModuleName, addr, request.AddrIpv6)
		}
	}
}
