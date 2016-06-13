package main

import (
	"fmt"

	"github.com/bo0mer/yamt/netstat"
)

func main() {
	fmt.Println(netstat.ReadIfStats())
}
