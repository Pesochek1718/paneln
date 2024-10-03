package nknovh_wasm

import (
	"fmt"
	"strconv"
	"syscall/js"
)

func (c *CLIENT) RegisterJSFuncs() {

	js.Global().Set("preSortNodes", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		if obj := x[0].Type(); obj != js.TypeObject {
			return false
		}
		if !x[0].Truthy() {
			return false
		}
		c.PreSortNodes(&x[0])
		c.mux.StartView.Lock()
		c.ParseAll()
		c.mux.StartView.Unlock()
		return true
	}))
	js.Global().Set("showModal", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		fmt.Println("showModal:"+x[0].String())
		c.ShowHideModal(x[0].String(), "show")
		return nil
	}))
	js.Global().Set("closeModal", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		c.ShowHideModal(x[0].String(), "hide")
		return nil
	}))
	js.Global().Set("genId", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.WsGenId()
		return nil
	}))
	js.Global().Set("toggleEnter", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		x := js.Global().Get("document").Call("getElementById", "s-enter")
		y := js.Global().Get("document").Call("getElementById", "toggleEnter")
		x.Set("innerHTML", y.Get("innerHTML"))
		return nil
	}))
	js.Global().Set("enterById", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		x := js.Global().Get("document").Call("getElementById", "hashId")
		hash := x.Get("value").String()
		c.WsAuth(hash)
		return nil
	}))
	js.Global().Set("logout", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.WsSend("logout")
		return nil
	}))

	js.Global().Set("setLanguage", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		c.SetLanguage(x[0].String(), x[1].String())
		return nil
	}))
	js.Global().Set("setEntriesPerPage", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		if num, err := strconv.Atoi(x[0].String()); err == nil {
			c.SetEntriesPerPage(num)
			return nil
		}
		fmt.Println("Cannot convert string to int")
		return false
	}))
	js.Global().Set("setPage", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		num := x[0].Int()
		c.SetPage(num)
		return nil
	}))
	js.Global().Set("prevPage", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.SetPage(c.CurrentPage - 1)
		return nil
	}))
	js.Global().Set("nextPage", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.SetPage(c.CurrentPage + 1)
		return nil
	}))

	js.Global().Set("addWalletLabels", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.AddWalletLabels(true)
		return nil
	}))
	js.Global().Set("toggleCheckBox", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.ToggleCheckBox()
		return nil
	}))
	js.Global().Set("addNodes", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.AddNodes()
		return nil
	}))
	js.Global().Set("createServer", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.CreateServer()
		return nil
	}))
	js.Global().Set("rmNodes", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.RmNodes()
		return nil
	}))
	js.Global().Set("saveSettings", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.SaveSettings()
		return nil
	}))
	js.Global().Set("hideAttention", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.HideAttention()
		return nil
	}))
	js.Global().Set("switchTab", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		c.SwitchTab(x[0].String())
		return nil
	}))
	js.Global().Set("checkOnline", js.FuncOf(func(_ js.Value, x []js.Value) interface{} {
		c.checkOnline(x[0].Int())
		return nil
	}))
	return
}

// "Db":
// 	[{
// 		"Host": "localhost",
// 		"Login": "caesar",
// 		"Database": "nkn",
// 		"Password": "caesar",
// 		"DbType":  "mysql",
// 		"MaxOpenConns": 128,
// 		"MaxIdleConns": 32,
// 		"InsideName": "main"
// 	}],
