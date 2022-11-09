package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

const customerThreadPool int = 10
const numOfCustomers int = 999
const (
	succeed = "\u2713"
	failed  = "\u2717"
)

type ComputeResults struct {
	Compute1 int
	Compute2 int
	Compute3 int
	Compute4 int
}

type TestResults struct {
	Compute1 string
	Compute2 string
	Compute3 string
	Compute4 string
}

func main() {
	var results []ComputeResults
	for i := 0; i < 10; i++ {
		var result ComputeResults
		result.Compute1 = compute1()
		var consumerGroup sync.WaitGroup
		consumerGroup.Add(customerThreadPool)
		result.Compute2, _ = compute2(&consumerGroup)
		result.Compute3, _ = compute3()
		result.Compute4, _ = compute4()
		results = append(results, result)
		resetArrays()
	}
	var expected ComputeResults
	expected.Compute1 = numOfCustomers
	expected.Compute2 = numOfCustomers
	expected.Compute3 = numOfCustomers
	expected.Compute4 = numOfCustomers
	succeeded := fmt.Sprintf("Succeeded! %v", succeed)
	failure := fmt.Sprintf("Failed! %v", failed)
	tr := TestResults{Compute1: succeeded, Compute2: succeeded, Compute3: succeeded, Compute4: succeeded}
	for _, test := range results {
		if test.Compute1 != expected.Compute1 && tr.Compute1 == succeeded {
			tr.Compute1 = failure
			//log.Println(test.Compute1)
		}
		if test.Compute2 != expected.Compute2 && tr.Compute2 == succeeded {
			tr.Compute2 = failure
			//log.Println(test.Compute2)
		}
		if test.Compute3 != expected.Compute3 && tr.Compute3 == succeeded {
			tr.Compute3 = failure
			//log.Println(test.Compute3)
		}
		if test.Compute4 != expected.Compute4 && tr.Compute4 == succeeded {
			tr.Compute4 = failure
			//log.Println(test.Compute4)
		}
	}
	log.Println("Ran each test 10x")
	log.Printf("Compute1 %v\n", tr.Compute1)
	log.Printf("Compute2 %v\n", tr.Compute2)
	log.Printf("Compute3 %v\n", tr.Compute3)
	log.Printf("Compute4 %v\n", tr.Compute4)

}

func compute1() int {
	count := map[string]int{}
	countChannel := make(chan int, numOfCustomers)
	resultChan := make(chan bool, numOfCustomers)

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
			// log.Println(i)
			countChannel <- i
		}(i)
	}

	//log.Println("waiting on lock...")
	lock.Wait()
	//log.Println("closing channel")
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

func resetArrays() {
	ACustomers = Customers{}
	ACustomers2 = Customers{}
}
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
