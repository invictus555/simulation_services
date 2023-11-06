package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"simulation_services/ker"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // 注册信号
	/* **************************** 调用业务函数 开始处 ****************************
	************** 注意：业务函数需要在协程中运行，避免阻塞后续的业务函数 ************* */

	go func() {
		ker.SDKFetchKerRuleGroupSimulationService(false) // 启动ker模拟程序
	}()

	/* ************************** 调用业务函数 结束处 ************************** */
	for { // 处理信号
		select {
		case <-sigCh: // 收到信号，退出
			os.Exit(0)
		default: // 主线程继续运行
			time.Sleep(5 * time.Second)
		}
	}
}
