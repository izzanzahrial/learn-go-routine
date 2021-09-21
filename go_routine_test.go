package learn_go_routine_test

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Go routine
// You cant catch return value in go routine

func RunHelloWorld() {
	fmt.Println("Hello world")
}

func TestHelloWorldGoroutine(t *testing.T) {
	go RunHelloWorld() // <-- "go" is the goroutine part, will run async
	fmt.Println("Done")

	time.Sleep(1 * time.Second)
}

// Test many go routine
func DisplayNumber(number int) {
	fmt.Println("Number :", number)
}

func TestDisplayNumber(t *testing.T) {
	for i := 0; i < 100000; i++ {
		go DisplayNumber(i)
	}
	time.Sleep(10 * time.Second)
}

// Channel (syncronous)
// a place where go routines can communicate syncronously
// since in go routine you can't return value, by using channel go routine can return value
// in channel there's deliver and reciever
// the channel will block the data from the deliver until there is a reciever
// channel is like async await
// if you dont close the channel when you dont use it, it could become memory leak

func TestChannel(t *testing.T) {
	channel := make(chan string)

	channel <- "Izzan" // to deliver data to channel
	name := <-channel  // to recieve data from channel

	fmt.Println(name)
	// fmt.Println(<-channel) ,another way to call channel directly
	close(channel) // to close the channel
	// defer close(channel) ,another way to close channel
	// by using defer, it will close the channel even tho there is an error in it

	// channel in goroutine
	go func() {
		time.Sleep(2 * time.Second)
		channel <- "Ahmad Izzan Zahrial"
		fmt.Println("The data has been delivered")
	}()

	channelData := <-channel
	fmt.Println(channelData)
	time.Sleep(5 * time.Second)
}

// Channel as parameter
// in channel you dont have to use pointer, since channel will not copy the value of data,
// but will copy the address(pass by reference), unlike the other parameter(pass by value)

func ChannelAsParameter(channel chan string) { // channel as parameter
	time.Sleep(2 * time.Second)
	channel <- "Ahmad Izzan Zahrial"
}

func TestChannelAsParameter(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go ChannelAsParameter(channel)

	dataChannel := <-channel
	fmt.Println(dataChannel)

	time.Sleep(5 * time.Second)
}

// Channel in and out
// You can make channel specifically for in(recieving) or out(delivering)

func ChannelIn(channel chan<- string) { // chan<- ,channel recieving data
	time.Sleep(2 * time.Second)
	channel <- "Ahmad Izzan Zahrial"
}

func ChannelOut(channel <-chan string) { // <-chan ,channel delivering data
	dataChannel := <-channel
	fmt.Println(dataChannel)
}

func TestChannelInOut(t *testing.T) {
	channel := make(chan string)

	go ChannelIn(channel)
	go ChannelOut(channel)

	time.Sleep(5 * time.Second)
	close(channel)
}

// Buffered Channel
// Since in channel if there are more than one data to be recieve, the other have to wait
// By using buffer, you can make a queue

func TestBufferedChannel(t *testing.T) {
	channel := make(chan string, 3) // "3" , buffer size
	defer close(channel)

	go func() {
		channel <- "Ahmad"
		channel <- "Izzan"
		channel <- "Zahrial"
	}()

	go func() {
		fmt.Println(<-channel)
		fmt.Println(<-channel)
		fmt.Println(<-channel)
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("Done")
}

// Range Channel
func TestRangeChannel(t *testing.T) {
	channel := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			channel <- strconv.Itoa(i)
		}
		close(channel)
	}()

	for data := range channel {
		fmt.Println("Data :", data)
	}

	fmt.Println("Done")
}

func ChannelResponse(channel chan<- string) {
	channel <- "camel"
}

// Select Channel
func TestSelectChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)
	defer close(channel1)
	defer close(channel2)

	go ChannelResponse(channel1)
	go ChannelResponse(channel2)

	counter := 0
	for {
		select {
		case data := <-channel1:
			fmt.Println("Data from channel 1:", data)
			counter++
		case data := <-channel2:
			fmt.Println("Data from channel 2:", data)
			counter++
		}
		if counter == 2 {
			break
		}
	}
}

// Default Select
func TestDefaultSelectChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)
	defer close(channel1)
	defer close(channel2)

	go ChannelResponse(channel1)
	go ChannelResponse(channel2)

	counter := 0
	for {
		select {
		case data := <-channel1:
			fmt.Println("Data from channel 1:", data)
			counter++
		case data := <-channel2:
			fmt.Println("Data from channel 2:", data)
			counter++
		default:
			fmt.Println("Processing...")
		}
		if counter == 2 {
			break
		}
	}
}

// Race Condition
// is a problem that occurs when go routine run parallely(because multiple thread)
// if goroutine share the same variable and they run parallely
// they will race with each other, thus race condition happen
// Race condition simulation
func TestRaceCondition(t *testing.T) {
	var x = 0
	for i := 1; i < 1000; i++ { // create 1000 go routine
		go func() { // each go routine loop 100 times
			for j := 1; j <= 100; j++ { // using x as increment variable
				x = x + 1 // the total should be 1000 * 100 = 100.000
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Total :", x) // total != 100.000
}

// sync.Mutex (Mutual Exclusion)
// sync.Mutex can be used to lock and unlock data, by using this, only one goroutine can access the data (like a queue)
// hence you can use sync.Mutex to fix race condition
// the cons : make the function a bit slower
// https://pkg.go.dev/sync#Mutex

func TestMutex(t *testing.T) {
	var x = 0
	var mutex sync.Mutex
	for i := 1; i < 1000; i++ {
		go func() {
			for j := 1; j <= 100; j++ {
				mutex.Lock() // lock the mutex, only 1 can access
				x = x + 1
				mutex.Unlock() // unlock the mutex, so the other can access it
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Total :", x) // total != 100.000
}

// RWMutex (Read Write Mutex)
// Mutex for read and write data, so the read and write doesn't have to fight for the lock
// because this mutex have 2 kind of lock, write lock and read lock
// https://pkg.go.dev/sync#RWMutex
type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) { // Write mutex
	account.RWMutex.Lock()
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BankAccount) GetBalance() int { // Read mutex (Rlock, RUnlock)
	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

func TestReadWriteMutex(t *testing.T) {
	account := BankAccount{}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				account.AddBalance(1)
				fmt.Println(account.GetBalance())
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Total balance:", account.Balance)
}

// Deadlock
// a situasion where a set of processes are blocked because each process is holding a resource
// and waiting for another resource acquired by some other process
// Deadlock simulation

type UserBalance struct {
	Mutex   sync.Mutex
	Name    string
	Balance int
}

func (user *UserBalance) Lock() {
	user.Mutex.Lock()
}

func (user *UserBalance) Unlock() {
	user.Mutex.Unlock()
}

func (user *UserBalance) Change(amount int) {
	user.Balance = user.Balance - amount
}

func Transfer(user1 *UserBalance, user2 *UserBalance, amount int) {
	user1.Lock()
	fmt.Println("Lock", user1.Name)
	user1.Change(-amount)

	time.Sleep(1 * time.Second)

	user2.Lock()
	fmt.Println("Lock", user2.Name)
	user2.Change(amount)

	time.Sleep(1 * time.Second)

	user1.Unlock()
	user2.Unlock()
}

func TestDeadLock(t *testing.T) {
	user1 := UserBalance{
		Name:    "Izzan",
		Balance: 4000,
	}
	user2 := UserBalance{
		Name:    "Zahrial",
		Balance: 4000,
	}

	// Because user1 and user2 trasnfer to each other at the same time, it's causing the deadlock
	// it's happen when the transfer1 trying to lock the user2, but the user2 already locked by transfer2
	// and also transfer2 trying to lock the user1, but the user1 already locked by transfer1
	// they wait each other indefinitely
	go Transfer(&user1, &user2, 1000)
	go Transfer(&user2, &user1, 1000)

	time.Sleep(5 * time.Second)

	fmt.Println("User:", user1.Name, ", Balance:", user1.Balance)
	fmt.Println("User:", user2.Name, ", Balance:", user2.Balance)
}

// WaitGroup
// function to wait gorountine until done
// https://pkg.go.dev/sync#WaitGroup
func RunAsync(group *sync.WaitGroup) {
	defer group.Done()

	group.Add(1)

	fmt.Println("Hello")
	time.Sleep(1 * time.Second)
}

func TestWaitGroup(t *testing.T) {
	group := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go RunAsync(group)
	}

	group.Wait()
	fmt.Println("Done")
}

// Once
// struct that make sure only once the function run
// https://pkg.go.dev/sync#Once
func PrintOnce() {
	fmt.Println("The only one")
}

func TestOnce(t *testing.T) {
	once := sync.Once{}
	group := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go func() {
			group.Add(1)
			once.Do(PrintOnce) // do this function only once
			group.Done()
		}()
	}

	group.Wait()
}

// Pool (Pool Pattern)
// Design pattern that used to store data, use the data, and store it back after it's done
// Reuseable
// example : database connection, so we dont have to keep making a new connection
// https://pkg.go.dev/sync#Pool

func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() interface{} { // a way to make default value in pool
			return "New"
		},
	}
	group := sync.WaitGroup{}

	pool.Put("Ahmad")
	pool.Put("Izzan")
	pool.Put("Zahrial")

	for i := 0; i < 10; i++ {
		go func() {
			group.Add(1)
			data := pool.Get()
			fmt.Println(data)
			pool.Put(data)
			group.Done()
		}()
	}

	group.Wait()
}

// sync.Map
// https://pkg.go.dev/sync#Map
func AddToMap(data *sync.Map, value int, group *sync.WaitGroup) {
	defer group.Done()

	group.Add(1)
	data.Store(value, value)
}

func TestMap(t *testing.T) {
	data := &sync.Map{}
	group := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go AddToMap(data, 1, group)
	}

	group.Wait()

	data.Range(func(key, value interface{}) bool {
		fmt.Println(key, ":", value)
		return true
	})
}

// Cond (condition)
// wait with statement
// Wait(), Signal(), BroadCast()
// https://pkg.go.dev/sync#Cond
var locker = &sync.Mutex{}
var cond = sync.NewCond(locker)
var group = &sync.WaitGroup{}

func WaitCondition(value int) {
	defer group.Done()
	group.Add(1)

	cond.L.Lock()
	cond.Wait() // wait until is it ok to proceed

	fmt.Println("Done", value)
	cond.L.Unlock()
}

func TestCond(t *testing.T) {
	for i := 0; i < 10; i++ {
		go WaitCondition(1)
	}

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			cond.Signal() // signaling the condition that is ok to run
		}
	}()

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	cond.Broadcast() // if you're using broadcast, it will signal all the goroutine that is ok to run
	// }()

	group.Wait()
}

// Atomic
// Package atomic provides low-level atomic memory primitives useful for implementing synchronization algorithms.
// in simple word, a way to use simple type data with no race condition
// https://pkg.go.dev/sync/atomic

func TestAtomic(t *testing.T) {
	var num int64 = 0
	group := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		go func() {
			group.Add(1)
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&num, 1)
			}
			group.Done()
		}()
	}
	group.Wait()
	fmt.Println("Number:", num)
}

// Timer
// https://pkg.go.dev/time#Timer

func TestTimer(t *testing.T) {
	timer := time.NewTimer(5 * time.Second)
	fmt.Println(time.Now())

	time := <-timer.C // passing the value from timer.Channel to time, the value should be 5 since we wait 5 second
	fmt.Println(time)
}

// time.After()
// to pass only the channel, without the timer data
func TestAfterTimer(t *testing.T) {
	channel := time.After(5 * time.Second)
	fmt.Println(time.Now())

	time := <-channel
	fmt.Println(time)
}

// time.AfterFunc()
// timer to run function
func TestAfterFunc(t *testing.T) {
	group := sync.WaitGroup{}
	group.Add(1)

	time.AfterFunc(1*time.Second, func() { // <-- run this function after 1 second
		fmt.Println("Run this")
		group.Done()
	})

	group.Wait()
}

// time.Ticker (channel)
// repetition per time
// https://pkg.go.dev/time#Ticker
func TestTicker(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		time.Sleep(5 * time.Second)
		ticker.Stop()
	}()

	for tick := range ticker.C {
		fmt.Println(tick)
	}
}

// time.Tick()
// repetition per time but only the channel data, without the timer data

func TestTick(t *testing.T) {
	channel := time.Tick(1 * time.Second)

	for time := range channel {
		fmt.Println(time)
	}
}

// GOMAXPROCS
// num of thread
func TestGetThread(t *testing.T) {
	totalCPU := runtime.NumCPU()
	fmt.Println("Total CPU:", totalCPU)

	totalThread := runtime.GOMAXPROCS(-1)
	fmt.Println("Total Thread:", totalThread)

	totalGoroutine := runtime.NumGoroutine()
	fmt.Println("Total Goroutine:", totalGoroutine)
}
