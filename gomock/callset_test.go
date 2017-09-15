// Copyright 2011 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gomock

import "testing"

func TestCallSetAdd(t *testing.T) {
	methodVal := "TestMethod"
	var receiverVal interface{} = "TestReceiver"
	cs := make(callSet)

	numCalls := 10
	for i := 0; i < numCalls; i++ {
		cs.Add(&Call{receiver: receiverVal, method: methodVal})
	}

	if len(cs) != 1 {
		t.Errorf("expected only one reciever in callSet")
	}
	if numActualMethods := len(cs[receiverVal]); numActualMethods != 1 {
		t.Errorf("expected only method on the reciever in callSet, found %d", numActualMethods)
	}
	if numActualCalls := len(cs[receiverVal][methodVal]); numActualCalls != numCalls {
		t.Errorf("expected all %d calls in callSet, found %d", numCalls, numActualCalls)
	}
}

func TestCallSetRemove(t *testing.T) {
	methodVal := "TestMethod"
	var receiverVal interface{} = "TestReceiver"

	cs := make(callSet)
	ourCalls := []*Call{}

	numCalls := 10
	for i := 0; i < numCalls; i++ {
		// NOTE: abuse the `numCalls` value to convey initial ordering of mocked calls
		generatedCall := &Call{receiver: receiverVal, method: methodVal, numCalls: i}
		cs.Add(generatedCall)
		ourCalls = append(ourCalls, generatedCall)
	}

	// validateOrder validates that the calls in the array are ordered as they were added
	validateOrder := func(calls []*Call) {
		// lastNum tracks the last `numCalls` (call order) value seen
		lastNum := -1
		for _, c := range calls {
			if lastNum >= c.numCalls {
				t.Errorf("found call %d after call %d", c.numCalls, lastNum)
			}
			lastNum = c.numCalls
		}
	}

	for _, c := range ourCalls {
		validateOrder(cs[receiverVal][methodVal])
		cs.Remove(c)
	}
}
