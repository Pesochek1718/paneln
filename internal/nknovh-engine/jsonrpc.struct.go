package nknovh_engine

import (
		"encoding/json"
		"time"
		"net/http"
)

type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method string      `json:"method"`
	Params *json.RawMessage `json:"params,omitempty"`
	Id     int      `json:"id"`
}

type JsonRPCConf struct {
	Timeout time.Duration
	Ip string
	Method string
	Params *json.RawMessage
	Client *http.Client
	UnmarshalData interface{}
}

type NodeSt struct {
	State NodeState
	Neighbor NodeNeighbor
}

type RPCErrorState struct {
	Code int `json:"code,omitempty"`
	Data string `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	WalletAddress string `json:"walletAddress,omitempty"`
}

type NodeState struct {
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error *RPCErrorState `json:"error,omitempty"`
	Result  struct {
		Addr               string `json:"addr"`
		Currtimestamp      int    `json:"currTimeStamp"`
		Height             int    `json:"height"`
		ID                 string `json:"id"`
		Jsonrpcport        int    `json:"jsonRpcPort"`
		LedgerMode         string `json:"ledgerMode"`
		ProposalSubmitted  int    `json:"proposalSubmitted"`
		ProtocolVersion    int    `json:"protocolVersion"`
		PublicKey          string `json:"publicKey"`
		RelayMessageCount  uint64 `json:"relayMessageCount"`
		SyncState          string `json:"syncState"`
		Tlsjsonrpcdomain   string `json:"tlsJsonRpcDomain"`
		Tlsjsonrpcport     int    `json:"tlsJsonRpcPort"`
		Tlswebsocketdomain string `json:"tlsWebsocketDomain"`
		Tlswebsocketport   int    `json:"tlsWebsocketPort"`
		Uptime             int    `json:"uptime"`
		Version            string `json:"version"`
		Websocketport      int    `json:"websocketPort"`
	} `json:"result"`
}


type NodeNeighbor struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error *RPCErrorState `json:"error,omitempty"`
	Result []struct {
		Addr               string `json:"addr"`
		Height             int    `json:"height"`
		ID                 string `json:"id"`
		Isoutbound         bool   `json:"isOutbound"`
		Jsonrpcport        int    `json:"jsonRpcPort"`
		LedgerMode         string `json:"ledgerMode"`
		Protocolversion    int    `json:"protocolVersion"`
		PublicKey          string `json:"publicKey"`
		RoundTripTime      int    `json:"roundTripTime"`
		SyncState          string `json:"syncState"`
		Tlsjsonrpcdomain   string `json:"tlsJsonRpcDomain"`
		Tlsjsonrpcport     int    `json:"tlsJsonRpcPort"`
		Tlswebsocketdomain string `json:"tlsWebsocketDomain"`
		Tlswebsocketport   int    `json:"tlsWebsocketPort"`
		Websocketport      int    `json:"websocketPort"`
		ConnTime           int    `json:"connTime"`
	} `json: "result"`
}
