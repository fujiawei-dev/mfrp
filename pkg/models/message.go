/*
 * @Date: 2022.01.22 18:10
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 18:10
 */

package models

const (
	Success = iota
	NotExists
	InvalidPassword
	AlreadyInUse
	Unexpected
	Unsupported
)

type GeneralRes struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

const (
	ControlConn = iota
	WorkConn
)

type ClientCtlReq struct {
	Type      int    `json:"Type"`
	Password  string `json:"Password"`
	ProxyName string `json:"ProxyName"`
}

type ClientCtlRes struct {
	GeneralRes
}

type ServerCtlReq struct {
	Type int `json:"Type"`
}
