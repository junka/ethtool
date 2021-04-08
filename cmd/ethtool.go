package main

import (
	ethtool "ethtool/pkg"
	"fmt"
)

func main() {
	err := ethtool.Parse_args()
	if err != nil {
		fmt.Printf("For more information run ethtool -h\n")
		return
	}
	// fmt.Println(a)
	ethtool.Do_actions()
	// ethtool.Show_usage()
}
