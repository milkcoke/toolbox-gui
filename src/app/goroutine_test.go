package app

import (
	"log"
	"sync"
	"testing"
)

func printRandom(num int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(num)
}

func Test_Print_Random(t *testing.T) {

	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go printRandom(i, wg)
	}

	wg.Wait()
}
