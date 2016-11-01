package main

// #include <unistd.h>
//
// void loop() {
//		for (;;) {
//			sleep(1);
//		}
// }
import "C"
import "fmt"
import "time"

func loop() {
	fmt.Println("loop")
	C.loop()
}

func foo() {
	time.Sleep(time.Second)
	fmt.Println("foo")
	go foo()
}

func main() {
	for i := 0; i < 100; i++ {
		go loop()
	}
	go foo()
	loop()
}

