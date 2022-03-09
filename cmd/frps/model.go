package main

import (
	"github.com/marsofsnow/eos-go"

)

var (
	ContractCaller               = eos.AccountName(`adxio.im`)
	ContractCallerPermissin = eos.PermissionName("active")
	ContractCallerPermissionLevel  = []eos.PermissionLevel{
		{Actor: ContractCaller, Permission: ContractCallerPermissin},
	}

	// im contract
	//查询
	IMContract = "adxio.im"
	TableFrps = "frps1"

)



type ChainFrps struct {
	Id int64 `json:"id"`
	Ip string `json:"ip"`
	IpHash string `json:"ip_hash"`
	UnusedPorts []uint32 `json:"unused_ports"`
	UsedPorts []string 	`json:"used_ports"`
	Status uint32 `json:"status"`
}


type ActionAddFrpsAddr struct {
	FrpsIp string `json:"frps_ip"`
	UnusedPorts []uint32 `json:"unused_ports"`
}

func (a *ActionAddFrpsAddr) to_eos_action() *eos.Action{
	return &eos.Action{
		Account:       ContractCaller,
		Name:          eos.ActionName("addfrpaddr"),
		Authorization: ContractCallerPermissionLevel,
		ActionData:    eos.NewActionData(a),
	}
}

type ActionTakePort struct {
	FrpsIp string `json:"frps_ip"`
	Port uint32 `json:"port"`
	FrpcIp string `json:"frpc_ip"`
}
func (a *ActionTakePort) to_eos_action() *eos.Action{
	return &eos.Action{
		Account:       ContractCaller,
		Name:          eos.ActionName("takeport"),
		Authorization: ContractCallerPermissionLevel,
		ActionData:    eos.NewActionData(a),
	}
}

type ActionReturnPort struct {
	FrpsIp string `json:"frps_ip"`
	Port uint32 `json:"port"`
}

func (a *ActionReturnPort) to_eos_action() *eos.Action{
	return &eos.Action{
		Account:       ContractCaller,
		Name:          eos.ActionName("returnport"),
		Authorization: ContractCallerPermissionLevel,
		ActionData:    eos.NewActionData(a),
	}
}