package main

import (
	"log"
	"strconv"
	"sync"
)

func main() {
	log.Println("compute should return 97, got ", compute())
}

func compute() int {
	count := map[string]int{}
	countChannel := make(chan int, 10)
	resultChan := make(chan bool, 97)

	go func() {
		for countItem := range countChannel {
			count[strconv.Itoa(countItem)] = countItem
			// if countItem == 73 {
			// 	time.Sleep(10 * time.Second)
			// }
			resultChan <- true
		}
	}()

	var lock sync.WaitGroup
	for i := 0; i < 97; i++ {
		lock.Add(1)
		go func(i int) {
			defer lock.Done()
			log.Println(i)
			countChannel <- i
		}(i)
	}

	log.Println("waiting on lock...")
	lock.Wait()
	log.Println("closing channel")
	close(countChannel)
	// for n := 0; n < 97; n++ {
	// 	<-resultChan
	// }
	// intslice := []int{}
	// for _, v := range count {
	// 	v := v
	// 	intslice = append(intslice, v)
	// }
	// sort.IntSlice(intslice).Sort()
	// for item := range intslice {
	// 	log.Println(item)
	// }
	return len(count)
}
