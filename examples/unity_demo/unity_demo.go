package main

import (
	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/examples/unity_demo/npc"
	"github.com/xiaonanln/goworld/examples/unity_demo/player"
)

func main() {
	goworld.RegisterSpace(&player.MySpace{}) // 注册自定义的Space类型

	goworld.RegisterService("OnlineService", &player.OnlineService{}, 3)
	goworld.RegisterService("SpaceService", &player.SpaceService{}, 3)
	// 注册Account类型
	goworld.RegisterEntity("Account", &player.Account{})
	// 注册Monster类型
	goworld.RegisterEntity("Monster", &npc.Monster{})
	// 注册Avatar类型，并定义属性
	goworld.RegisterEntity("Player", &player.Player{})
	// 运行游戏服务器
	goworld.Run()
}
