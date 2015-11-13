package state_test

import ( 
	"fmt"

	"go.permanent.de/anaLog/state"
)

func ExampleAtos() {
	fmt.Println(state.Started == state.Atos("Started"))	
	fmt.Println(state.Atos("asd"))
	// Output: 
	// true
	// Unknown
}

func Example() {
	fmt.Println(state.Failed == state.Atos("Running"))
	fmt.Println(state.OK)
	// Output: 
	// false
	// OK
}