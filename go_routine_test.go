package learn_go_routine_test

import (
	"fmt"
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
