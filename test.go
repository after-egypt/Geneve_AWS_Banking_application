package main

import (
	"fmt"
	//"math/rand"
	//"sync"
	//"time"
	"encoding/binary"
	"encoding/hex"
)

/*
func process(ch chan int, wg sync.WaitGroup) {
	i := 0
	for {
		ch <- i
		if rand.Float64() >= .5 {
			time.Sleep(1 * time.Second)
		}
		wg.Add(1)
		wg.Done()
		i++

	}
}*/
func main() {
	b := make([]byte, 3)
	binary.BigEndian.PutUint16(b[1:], 0xABCD)
	fmt.Println(hex.EncodeToString(b))
}
