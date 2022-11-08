package main

import (
	"log"
	"strconv"
	"sync"
)

const customerThreadPool int = 10
const numOfCustomers int = 999

func main() {
	log.Println("compute should return 97, got ", compute())
	var consumerGroup sync.WaitGroup
	consumerGroup.Add(customerThreadPool)
	custLen, _ := compute2(&consumerGroup)
	log.Printf("compute2 should return %v, got %v\n", numOfCustomers, custLen)
	// for cust := range custs {
	// 	log.Printf("%v\n", cust)
	// }
	custLen2, _ := compute3()
	log.Printf("compute3 should return %v, got %v\n", numOfCustomers, custLen2)
	// for cust := range custs2 {
	// 	log.Printf("%v\n", cust)
	// }
	custLen3, _ := compute4()
	log.Printf("compute4 should return %v, got %v\n", numOfCustomers, custLen3)
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
			// If you uncomment the line below, you will get inconsistent results
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

func compute2(consumerGroup *sync.WaitGroup) (int, []int) {
	custLock := struct {
		Customers []int
		Mu        sync.Mutex
	}{}
	customerChannel := make(chan int, 1000)
	results := make(chan bool, numOfCustomers)
	for i := 0; i < numOfCustomers; i++ {
		customerChannel <- i
	}
	close(customerChannel)
	for n := 0; n < customerThreadPool; n++ {
		go func() {
			defer consumerGroup.Done()
			for cust := range customerChannel {
				/* IF YOU UNCOMMENT THESE LINES, IT WILL SEEM TO WORK - DANGER! */
				// for j := 0; j < 21; j++ {
				// 	count := 0
				// 	for count < 2 {
				// 		count++
				// 		needMorePass := true
				// 		for needMorePass {
				// 			needMorePass = false
				// 			var shardgroup sync.WaitGroup
				// 			shardgroup.Add(97)
				// 			for parentShard0 := 0; parentShard0 < 97; parentShard0++ {
				// 				go func(parentShard int) {
				// 					defer shardgroup.Done()
				// 					// doing nothing more than looping
				// 				}(parentShard0)
				// 			}
				// 			shardgroup.Wait()
				// 		}
				// 	}
				// }
				c := cust
				custLock.Customers = append(custLock.Customers, c)
				results <- true
			}
		}()
	}
	for i := 0; i < numOfCustomers; i++ {
		<-results
	}
	consumerGroup.Wait()
	return len(custLock.Customers), custLock.Customers
}

type Customers struct {
	Customers []int
	Mu        sync.Mutex
}

var ACustomers Customers
var ACustomers2 Customers

func computeWorkerNoLock(id int, jobs <-chan int, results chan<- bool) {
	for job := range jobs {
		//log.Printf("worker %v, job %v", id, job)
		//ACustomers.Mu.Lock()
		ACustomers.Customers = append(ACustomers.Customers, job)
		//ACustomers.Mu.Unlock()
		results <- true
	}
}
func computeWorker(id int, jobs <-chan int, results chan<- bool) {
	for job := range jobs {
		//log.Printf("worker %v, job %v", id, job)
		ACustomers2.Mu.Lock()
		ACustomers2.Customers = append(ACustomers2.Customers, job)
		ACustomers2.Mu.Unlock()
		results <- true
	}
}

func compute3() (int, []int) {
	customerChannel := make(chan int, numOfCustomers)
	results := make(chan bool, numOfCustomers)
	for n := 0; n < customerThreadPool; n++ {
		go computeWorkerNoLock(n, customerChannel, results)
	}
	for i := 0; i < numOfCustomers; i++ {
		n := i
		customerChannel <- n
	}
	close(customerChannel)

	for j := 0; j < numOfCustomers; j++ {
		<-results
	}
	return len(ACustomers.Customers), ACustomers.Customers
}

func compute4() (int, []int) {
	customerChannel := make(chan int, numOfCustomers)
	results := make(chan bool, numOfCustomers)
	for n := 0; n < customerThreadPool; n++ {
		go computeWorker(n, customerChannel, results)
	}
	for i := 0; i < numOfCustomers; i++ {
		n := i
		customerChannel <- n
	}
	close(customerChannel)

	for j := 0; j < numOfCustomers; j++ {
		<-results
	}
	return len(ACustomers2.Customers), ACustomers2.Customers
}
