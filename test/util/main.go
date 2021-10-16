package main

import (
	"fmt"
	"gonpy/trader/util"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println(util.RoundTo(1000.2, 0.1))
	fmt.Println(util.RoundTo(1000.23, 0.1))
	fmt.Println(util.RoundTo(1000.0, 1))

	v := viper.New()
	// v.AddConfigPath("../../")
	// v.SetConfigType("dotenv")
	// v.SetConfigName("env")
	// v.AutomaticEnv()
	// v.SetConfigFile("../../.env")
	v.SetConfigFile("../../docker-compose.yml")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	// vv := v.Get("mongo_port")
	// vv := v.GetInt("MONGO_PORT")
	// fmt.Println(vv)

	// dockComposeConfig := &DockerComposeYml{}
	// if err:=v.Unmarshal(dockComposeConfig);err!=nil{
	// 	panic(err)
	// }
	// pp := strings.Split(dockComposeConfig.Services.Influxdb.Ports[0], ":")[0]
	// fmt.Println(1, pp)

	dd := util.GetDockerComposeYml("../../docker-compose.yml")
	fmt.Println(dd)
}
