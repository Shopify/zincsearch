package index

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zincsearch/zincsearch/pkg/core"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
	"github.com/zincsearch/zincsearch/test/utils"
)

func TestCatIndices(t *testing.T) {
	t.Run("prepare", func(t *testing.T) {
		index, err := core.NewIndex("TestCatIndices.index_1", "disk", 2)
		assert.NoError(t, err)
		assert.NotNil(t, index)

		err = core.StoreIndex(index)
		assert.NoError(t, err)
	})

	t.Run("catIndices", func(t *testing.T) {
		formats := []string{"json", "table"}
		verboseArr := []string{"false", "true"}
		for _, format := range formats {
			for _, verbose := range verboseArr {
				c, w := utils.NewGinContext()
				params := map[string]string{
					"page_num":  "1",
					"page_size": "20",
					"format":    format,
					"v":         verbose,
					"h":         "health,status,index,uuid,pri,rep,docs.count,docs.deleted,store.size,pri.store.size",
				}
				utils.SetGinRequestURL(c, "/es/_cat/indices", params)
				CatIndices(c)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.NotNil(t, w.Body)

				if format == "json" {
					var resp []CatIndicesResponseItem
					err := json.Unmarshal(w.Body.Bytes(), &resp)
					assert.NoError(t, err)
					assert.NotNil(t, resp)
					assert.Equal(t, len(resp), 1)
					assert.Equal(t, resp[0].Index, "TestCatIndices.index_1")
				} else {
					assert.Contains(t, w.Body.String(), "TestCatIndices.index_1")
				}
			}
		}
	})

	t.Run("cleanup", func(t *testing.T) {
		err := core.DeleteIndex("TestCatIndices.index_1")
		assert.NoError(t, err)
	})
}
