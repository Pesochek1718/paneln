package nknovh_engine

import (
	"encoding/json"
	"log"
	"net/http"
)

type ProxyData struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type ProxyResponse struct {
	Code    int         `json:"code"`
	Data    []ProxyData `json:"data"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
}

func getIP() (string, string) {
	lunaproxy := "https://tq.lunaproxy.com/getflowip?neek=1295458&num=1&type=2&sep=1&regions=us&ip_si=2&level=1&sb="
	resp, err := http.Get(lunaproxy)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var proxyResp ProxyResponse
	err = json.NewDecoder(resp.Body).Decode(&proxyResp)
	if err != nil {
		log.Fatalln(err)
	}

	if !proxyResp.Success {
		log.Fatalf("Ошибка запроса: %s", proxyResp.Msg)
	}

	if len(proxyResp.Data) > 0 {
		return proxyResp.Data[0].IP, proxyResp.Data[0].Port
	} else {
		log.Println("Данные не найдены")
		return "", ""
	}
}
