package ker

import (
	"fmt"
	"math/rand"
	"time"

	"simulation_services/utils"
)

func init() {
	fetchDataAssistance()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			<-ticker.C
			fetchDataAssistance()
		}
	}()
}

// SDKFetchKerRuleGroupSimulationService 模拟调用ker服务[机房粒度]
func SDKFetchKerRuleGroupSimulationService(randomly bool) {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
	for {
		if dataAssistance == nil {
			fmt.Println("parameter(dataAssistance) is empty")
			break
		}
		if len(dataAssistance.RequestList) == 0 || len(dataAssistance.kerHosts) == 0 {
			fmt.Println("request inner parameters are empty")
			break
		}

		request := getSDKFetchRuleGroupRequest(randomly)
		body, err := serializeSDKFetchRuleGroupRequest2JSON(request) // 序列化SDK Fetch RuleGroup的请求参数
		if err != nil || len(body) == 0 {
			continue
		}

		url := "https://vip-boe-i18n.byted.org/ker/api/sdk/fetch_rule_groups"
		response, err := utils.DoHttpGetMethodV2(url, []byte(body))
		if err != nil || len(response) == 0 {
			fmt.Println("Do http POST failed, err:", err)
		} else {
			fmt.Printf("Successfully[%s]:randomly:%v PSM:%-32s module:%-6s IP:%s\n", utils.GetNowTime(), randomly, request.SdkHostPsm, request.ModuleName, request.AddrIpv6)
		}
	}
}
