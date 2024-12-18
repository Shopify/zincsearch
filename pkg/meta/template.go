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

import "time"

type Template struct {
	Name          string         `json:"name"`
	IndexTemplate *IndexTemplate `json:"index_template"`
}

type IndexTemplate struct {
	IndexPatterns []string         `json:"index_patterns"`
	Priority      int              `json:"priority,omitempty"` // highest priority is chosen
	Template      TemplateTemplate `json:"template"`
	CreatedAt     time.Time        `json:"created_at,omitempty"`
	UpdatedAt     time.Time        `json:"updated_at,omitempty"`
}

type TemplateTemplate struct {
	Settings *IndexSettings `json:"settings,omitempty"`
	Mappings *Mappings      `json:"mappings,omitempty"`
}
