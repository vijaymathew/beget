// Copyright 2018, 2019 Vijay Mathew Pandyalakal<vijay.the.lisper@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"fmt"
	"github.com/vijaymathew/beget"
)

func main() {
	resp, err := beget.Get(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("%v\n", resp)
	}
}
