/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
)

func main() {
	d := dependency{
		id:             os.Args[1],
		versionPattern: os.Args[2],
	}

	if err := d.update(os.Args[3], os.Args[4], os.Args[5]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to update %s %s\n", d.id, d.versionPattern)
		os.Exit(101)
	}
}
