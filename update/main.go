package main

import "fmt"

func main() {
	// var wg sync.WaitGroup
	// wg.Add(1)
	// // go someServer("192.248.150.228", "root", "2P#f!KE=rBB+4tQ8", "43.159.29.14:30758")
	// // C:\users\ярослав\appdata\local\temp\go-build4198649741\b001\exe\nknovh.exe
	// go someServer("78.141.237.2", "root", "9S%aADiazhW[)4F!", "43.155.29.14:30758")
	// wg.Wait()
	//chekNameNodes()

	newDir := chekNameNodesAdd()
	if newDir == "" {
		fmt.Println("good")
	}
	fmt.Println(newDir)
	// args := chekBusyIp("192.248.150.228")
	// fmt.Println(args)
}
