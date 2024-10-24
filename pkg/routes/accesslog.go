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

package routes

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

func AccessLog(app *gin.Engine) {
	app.Use(func(c *gin.Context) {
		timeStart := time.Now()
		c.Writer.Header().Set("Zinc", meta.Version)

		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		c.Request.Body = rdr2

		c.Next()

		took := time.Since(timeStart) / time.Millisecond

		log.Info().
			Str("method", c.Request.Method).
			Int("code", c.Writer.Status()).
			Int("took", int(took)).
			Str("body", readBody(rdr1)).
			Msg(c.Request.RequestURI)
	})
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
