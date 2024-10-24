/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package meta

import (
	"strconv"
)

// go build -a -ldflags "-X github.com/zincsearch/zincsearch/pkg/meta.AuthEnable=false" -o zincsearch
var (
	// Prefix used for ES routes
	ESRoutePrefix = "/es"

	// Enable GUI route
	GUIEnable = "true"

	// Enable Auth
	AuthEnable = "true"
)

func IsGUIEnabled() bool {
	v, err := strconv.ParseBool(GUIEnable)
	if err != nil {
		return true
	}
	return v
}

func IsAuthEnabled() bool {
	v, err := strconv.ParseBool(AuthEnable)
	if err != nil {
		return true
	}
	return v
}
