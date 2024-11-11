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

package token

import (
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/token"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewSynonymTokenFilter(options interface{}) (analysis.TokenFilter, error) {
	synonyms, err := zutils.GetStringSliceFromMap(options, "synonyms")
	if err != nil {
		return nil, errors.New(errors.ErrorTypeParsingException, "[token_filter] keyword option [synonyms] should be an array of string")
	}

	dict := analysis.NewTokenMap()
	for _, synonym := range synonyms {
		if strings.Contains(synonym, "=>") {
			for _, syn := range strings.Split(synonym, "=>") {
				for _, s := range strings.Split(syn, ",") {
					dict.AddToken(s)
				}
			}
		} else {
			for _, s := range strings.Split(synonym, ",") {
				dict.AddToken(s)
			}
		}
	}
	return token.NewKeyWordMarkerFilter(dict), nil
}
