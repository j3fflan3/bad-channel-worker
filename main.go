package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

const customerThreadPool int = 10
const numOfCustomers int = 999
const (
	succeed = "\u2713"
	failed  = "\u2717"
)

type ComputeResults struct {
	Compute1  int
	MCompute1 map[string]int
	Compute2  int
	ACompute2 []int
	Compute3  int
	ACompute3 []int
	Compute4  int
	ACompute4 []int
	Compute5  int
	ACompute5 []int
	Compute6  int
	ACompute6 []UpdateEvent
	Compute7  int
	ACompute7 [][]int
}

type TestResults struct {
	Compute1 string
	Compute2 string
	Compute3 string
	Compute4 string
	Compute5 string
	Compute6 string
	Compute7 string
}

func main() {
	var results []ComputeResults
	for i := 0; i < 10; i++ {
		var result ComputeResults
		result.Compute1, result.MCompute1 = compute1()
		var consumerGroup sync.WaitGroup
		consumerGroup.Add(customerThreadPool)
		///time.Sleep(1 * time.Second)
		//log.Println(len(result.MCompute1))
		result.Compute2, result.ACompute2 = compute2(&consumerGroup)
		result.Compute3, result.ACompute3 = compute3()
		result.Compute4, result.ACompute4 = compute4()
		result.Compute5, result.ACompute5 = compute5()
		result.Compute6, result.ACompute6 = compute6()
		result.Compute7, result.ACompute7 = compute7()
		//time.Sleep(1 * time.Second)
		results = append(results, result)
		resetArrays()
	}
	var expected ComputeResults
	expected.Compute1 = numOfCustomers
	expected.Compute2 = numOfCustomers
	expected.Compute3 = numOfCustomers
	expected.Compute4 = numOfCustomers
	expected.Compute5 = numOfCustomers
	expected.Compute6 = numOfCustomers
	expected.Compute7 = numOfCustomers
	succeeded := fmt.Sprintf("Succeeded! %v", succeed)
	failure := fmt.Sprintf("Failed! %v", failed)
	tr := TestResults{
		Compute1: succeeded,
		Compute2: succeeded,
		Compute3: succeeded,
		Compute4: succeeded,
		Compute5: succeeded,
		Compute6: succeeded,
		Compute7: succeeded,
	}
	var c1, c2, c3, c4, c5, c6, c7 []int
	for _, test := range results {
		val1 := test.Compute1
		c1 = append(c1, val1)
		if test.Compute1 != expected.Compute1 && tr.Compute1 == succeeded {
			tr.Compute1 = failure
			//log.Println(test.Compute1)
		}
		val2 := test.Compute2
		c2 = append(c2, val2)
		if test.Compute2 != expected.Compute2 && tr.Compute2 == succeeded {
			tr.Compute2 = failure
			//log.Println(test.Compute2)
		}
		val3 := test.Compute3
		c3 = append(c3, val3)
		if test.Compute3 != expected.Compute3 && tr.Compute3 == succeeded {
			tr.Compute3 = failure
			//log.Println(test.Compute3)
		}
		val4 := test.Compute4
		c4 = append(c4, val4)
		if test.Compute4 != expected.Compute4 && tr.Compute4 == succeeded {
			tr.Compute4 = failure
			//log.Println(test.Compute4)
		}
		val5 := test.Compute5
		c5 = append(c5, val5)
		if test.Compute5 != expected.Compute5 && tr.Compute5 == succeeded {
			tr.Compute5 = failure
			//log.Println(test.Compute4)
		}
		val6 := test.Compute6
		c6 = append(c6, val6)
		if test.Compute6 != expected.Compute6 && tr.Compute6 == succeeded {
			tr.Compute6 = failure
			//log.Println(test.Compute4)
		}
		val7 := test.Compute7
		c7 = append(c7, val7)
		if test.Compute7 != expected.Compute7 && tr.Compute7 == succeeded {
			tr.Compute7 = failure
			//log.Println(test.Compute4)
		}
	}
	log.Println("Ran each test 10x")
	log.Printf("Compute1 %v\tExpected all %v, got %v\n", tr.Compute1, expected.Compute1, c1)
	log.Printf("Compute2 %v\tExpected all %v, got %v\n", tr.Compute2, expected.Compute2, c2)
	log.Printf("Compute3 %v\tExpected all %v, got %v\n", tr.Compute3, expected.Compute3, c3)
	log.Printf("Compute4 %v\tExpected all %v, got %v\n", tr.Compute4, expected.Compute4, c4)
	log.Printf("Compute5 %v\tExpected all %v, got %v\n", tr.Compute5, expected.Compute5, c5)
	log.Printf("Compute6 %v\tExpected all %v, got %v\n", tr.Compute6, expected.Compute6, c6)
	log.Printf("Compute7 %v\tExpected all %v, got %v\n", tr.Compute7, expected.Compute7, c7)
}

func compute1() (int, map[string]int) {
	count := map[string]int{}
	countChannel := make(chan int, numOfCustomers)
	resultChan := make(chan bool, numOfCustomers)

	go func() {
		for countItem := range countChannel {
			count[strconv.Itoa(countItem)] = countItem
			resultChan <- true
		}
	}()

	var lock sync.WaitGroup
	for i := 0; i < numOfCustomers; i++ {
		lock.Add(1)
		go func(i int) {
			defer lock.Done()
			// If you uncomment the line below, you will get inconsistent results
			//log.Println(i)
			countChannel <- i
		}(i)
	}

	lock.Wait()
	close(countChannel)
	// for n := 0; n < numOfCustomers; n++ {
	// 	<-resultChan
	// }
	return len(count), count
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
	aUpdateEvents = []UpdateEvent{}
}
func computeWorkerNoLock(id int, jobs <-chan int, results chan<- bool) {
	for job := range jobs {
		//log.Printf("worker %v, job %v", id, job)
		ACustomers.Customers = append(ACustomers.Customers, job)
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

func compute5() (int, []int) {
	// Our customer array we want to populate
	var aCustomers Customers
	// The job and results channels
	jobs := make(chan int, numOfCustomers)
	results := make(chan bool, numOfCustomers)

	// Spawn worker routines
	for n := 0; n < customerThreadPool; n++ {
		go func(id int, jobs <-chan int, results chan<- bool) {
			for job := range jobs {
				//log.Printf("worker %v, job %v", id, job)
				aCustomers.Mu.Lock()
				aCustomers.Customers = append(aCustomers.Customers, job)
				aCustomers.Mu.Unlock()
				results <- true
			}

		}(n, jobs, results)
	}

	// Load jobs buffer
	for i := 0; i < numOfCustomers; i++ {
		n := i
		jobs <- n
	}
	// Close jobs chan
	close(jobs)

	// loop and block until all results come back
	for j := 0; j < numOfCustomers; j++ {
		<-results
	}
	// return our results
	return len(aCustomers.Customers), aCustomers.Customers
}

type UpdateEvent struct {
	UpdateId         uint64    `json:"updId"`
	BatchId          uint64    `json:"batch"`
	GraphId          uint64    `json:"-"`
	ObjectId         string    `json:"objId"`
	ObjectType       string    `json:"objType"`
	DoNotUse         bool      `json:"dnu"`
	EntitledDate     time.Time `json:"entDt"`
	EntitledAt       uint64    `json:"-"`
	EntitledChildren []string  `json:"entitled_children"`
}

var aUpdateEvents []UpdateEvent
var mu sync.Mutex

// compute6 mimics "ToParquet"
func compute6() (int, []UpdateEvent) {
	eventChannel := make(chan UpdateEvent, 1000000)
	var wg sync.WaitGroup
	wg.Add(2)
	go createParquetFiles(&wg, "", eventChannel)
	go createParquetFiles(&wg, "", eventChannel)
	startPumpingLogEvents(eventChannel)
	close(eventChannel)
	wg.Wait()
	return len(aUpdateEvents), aUpdateEvents
}

func createParquetFiles(wg *sync.WaitGroup, path string, eventChannel chan UpdateEvent) {
	defer wg.Done()
	eventsInFile := 0
	eventForFB := UpdateEvent{}
	for event := range eventChannel {
		eventForFB.ObjectType = event.ObjectType
		eventForFB.ObjectId = event.ObjectId
		eventForFB.UpdateId = uint64(event.UpdateId)
		eventForFB.DoNotUse = event.DoNotUse
		eventForFB.EntitledAt = uint64(event.EntitledAt)
		mu.Lock()
		aUpdateEvents = append(aUpdateEvents, event)
		mu.Unlock()
		eventsInFile++
	}
	//log.Printf("createParquetFiles::eventsInFile %v\n", eventsInFile)
}

func startPumpingLogEvents(eventChannel chan UpdateEvent) {
	count := 0
	events := make([]UpdateEvent, 1024*1024)
	for {
		if count == numOfCustomers {
			break
		}
		count = stream(events)
		for i := 0; i < count; i++ {
			eventChannel <- events[i]
		}
	}
}

func stream(events []UpdateEvent) int {
	total := 0
	for objId := 0; objId < numOfCustomers; objId++ {
		event := UpdateEvent{
			ObjectType: "object type",
			ObjectId:   fmt.Sprint(objId),
			UpdateId:   uint64(objId),
			DoNotUse:   false,
			EntitledAt: uint64(objId),
		}
		events = append(events, event)
		total = objId + 1
	}
	return total
}

const s3consumers int = 100

var fileMu sync.Mutex

// Mimics s3Downloader.DownloadBatch
func compute7() (int, [][]int) {
	files := make([][]int, s3consumers)
	fileInfo := make(chan int, s3consumers)
	var wg sync.WaitGroup
	wg.Add(s3consumers)
	for i := 0; i < s3consumers; i++ {
		go func(workerId int) {
			defer wg.Done()
			for info := range fileInfo {
				//fileMu.Lock()
				files[workerId] = append(files[workerId], info)
				//fileMu.Unlock()
			}
		}(i)
	}

	for i := 0; i < numOfCustomers; i++ {
		n := i
		fileInfo <- n
	}
	close(fileInfo)
	wg.Wait()
	count := 0
	for i := 0; i < len(files); i++ {
		count += len(files[i])
	}
	return count, files
}
