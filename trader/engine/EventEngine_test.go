package engine

import (
	"fmt"
	"gonpy/trader"
	"reflect"
	"testing"
)

func TestEventEngine(t *testing.T) {
	type testCase struct {
		EventEngine *EventEngine
		val         int
		want        interface{}
	}

	testGroup := map[string]testCase{
		"case1": {
			EventEngine: &EventEngine{
				Name:     "case1EventEngine",
				Interval: 1,
				Queue:    make(chan *trader.Event),
				Active:   false,
				Thread:   "",
				Timer:    "",
				Handlers: map[string][]HandlerFunc{
					"eTimer": {ETimerHandlerFunc, ETimerHandlerFunc},
				},
			},
			val:  6,
			want: false,
		},
	}

	for name, tc := range testGroup {
		t.Run(name, func(t *testing.T) {
			// got := ListNodeToSlice(removeElements(tc.head, tc.val))
			got := tc.EventEngine.Active
			tc.EventEngine.start()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("testGroup, want:%v, got:%v", tc.want, got)
			} else {
				t.Logf("testGroup, want:%v, got:%v", tc.want, got)
			}
		})
	}
}

func ETimerHandlerFunc(event *trader.Event) {
	fmt.Println(event)
}
