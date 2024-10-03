module nknovh

go 1.21

toolchain go1.22.1

replace (
	nknovh-engine v1.1.0 => ./internal/nknovh-engine
	nknovh-wasm v1.0.0 => ./internal/nknovh-wasm
	templater v1.0.0 => ./internal/templater
	xwasmapi v1.0.0 => ./internal/xwasmapi
)

require (
	nknovh-engine v1.1.0
	nknovh-wasm v1.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/bramvdbogaerde/go-scp v1.4.0 // indirect
	github.com/fvbommel/sortorder v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/itchyny/base58-go v0.0.5 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/nknorg/ncp-go v1.0.5 // indirect
	github.com/nknorg/nkn-sdk-go v1.4.7 // indirect
	github.com/nknorg/nkn/v2 v2.1.7 // indirect
	github.com/nknorg/nkngomobile v0.0.0-20220615081414-671ad1afdfa9 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pbnjay/memory v0.0.0-20190104145345-974d429e7ae4 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	templater v1.0.0 // indirect
	xwasmapi v1.0.0 // indirect
)
