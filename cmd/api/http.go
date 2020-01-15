package api

import (
	"demo/api"
	"demo/pkg/config"
	"demo/pkg/jaeger"

	"github.com/spf13/cobra"
)

// ServerCmd http server
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "start http api server",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// 初始化依赖扩展
		jaeger.Init(
			config.GetString("jaeger.service_name"),
			config.GetString("jaeger.sampler_type"),
			config.GetFloat64("jaeger.sampler_param"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 启动入口
		api.Run()
	},
}
