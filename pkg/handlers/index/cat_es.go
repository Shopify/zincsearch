package index

import (
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zincsearch/zincsearch/pkg/core"
)

// @Id CatIndices
// @Summary Cat indexes
// @security BasicAuth
// @Tags    Index
// @Param   format query  string  false  "output format"
// @Param   v      query  bool    false  "verbose"
// @Param   h      query  string  false  "columns to display"
// @Produce json/text
// @Success 200 {object} IndexListResponse
// @Router /es/_cat/indices/{target} [get]
func CatIndices(c *gin.Context) {
	format := c.DefaultQuery("format", "table")
	verbose, _ := strconv.ParseBool(c.DefaultQuery("v", "false"))
	columns := strings.Split(c.DefaultQuery("h", strings.Join(catIndicesColumns, ",")), ",")

	for _, column := range columns {
		switch column {
		case "i":
			column = "index"
		}

		if !slices.Contains(catIndicesColumns, column) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid column: %s", column)})
			return
		}
	}

	if format != "json" && format != "table" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid format: %s", format)})
		return
	}

	items := core.ZINC_INDEX_LIST.ListStat()

	targetIndex := c.Param("target")

	var response []CatIndicesResponseItem
	for _, item := range items {
		if len(targetIndex) > 0 {
			if match, err := filepath.Match(targetIndex, item.GetName()); !match || err != nil {
				continue
			}
		}

		response = append(response, CatIndicesResponseItem{
			Health:       "green",
			Status:       "open",
			Index:        item.GetName(),
			UUID:         "N/A",
			Pri:          "1",
			Rep:          "0",
			DocsCount:    strconv.Itoa(int(item.GetStats().DocNum)),
			DocsDeleted:  "N/A",
			StoreSize:    formatBytes(int(item.GetStats().StorageSize)),
			PriStoreSize: formatBytes(int(item.GetStats().StorageSize)),
		})
	}

	if format == "json" {
		c.JSON(http.StatusOK, response)
	} else {
		c.String(http.StatusOK, buildIndicesTable(response, columns, verbose))
	}
}

func buildIndicesTable(data []CatIndicesResponseItem, columns []string, verbose bool) string {
	var width = map[string]int{}
	for _, column := range columns {
		width[column] = len(column)
	}
	for _, item := range data {
		for _, column := range columns {
			switch column {
			case "index":
				if len(item.Index) > width[column] {
					width[column] = len(item.Index)
				}
			}
		}
	}

	var res []string

	if verbose {
		line := []string{}
		for _, column := range columns {
			line = append(line, fmt.Sprintf("%s%s", column, strings.Repeat(" ", width[column]-len(column))))
		}
		res = append(res, strings.Join(line, " "))
	}
	for _, item := range data {
		line := []string{}
		for _, column := range columns {
			switch column {
			case "health":
				line = append(line, fmt.Sprintf("%s%s", item.Health, strings.Repeat(" ", width[column]-len(item.Health))))
			case "status":
				line = append(line, fmt.Sprintf("%s%s", item.Status, strings.Repeat(" ", width[column]-len(item.Status))))
			case "index":
				line = append(line, fmt.Sprintf("%s%s", item.Index, strings.Repeat(" ", width[column]-len(item.Index))))
			case "uuid":
				line = append(line, fmt.Sprintf("%s%s", item.UUID, strings.Repeat(" ", width[column]-len(item.UUID))))
			case "pri":
				line = append(line, fmt.Sprintf("%s%s", item.Pri, strings.Repeat(" ", width[column]-len(item.Pri))))
			case "rep":
				line = append(line, fmt.Sprintf("%s%s", item.Rep, strings.Repeat(" ", width[column]-len(item.Rep))))
			case "docs.count":
				line = append(line, fmt.Sprintf("%s%s", item.DocsCount, strings.Repeat(" ", width[column]-len(item.DocsCount))))
			case "docs.deleted":
				line = append(line, fmt.Sprintf("%s%s", item.DocsDeleted, strings.Repeat(" ", width[column]-len(item.DocsDeleted))))
			case "store.size":
				line = append(line, fmt.Sprintf("%s%s", item.StoreSize, strings.Repeat(" ", width[column]-len(item.StoreSize))))
			case "pri.store.size":
				line = append(line, fmt.Sprintf("%s%s", item.PriStoreSize, strings.Repeat(" ", width[column]-len(item.PriStoreSize))))
			}
		}
		res = append(res, strings.Join(line, " "))
	}
	return strings.Join(res, "\n")
}

var catIndicesColumns = []string{
	"health",
	"status",
	"index",
	"uuid",
	"pri",
	"rep",
	"docs.count",
	"docs.deleted",
	"store.size",
	"pri.store.size",
}

type CatIndicesResponseItem struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

// @Id CatAliases
// @Summary Cat aliases
// @security BasicAuth
// @Tags    Index
// @Param   format query  string  false  "output format"
// @Param   v      query  bool    false  "verbose"
// @Param   h      query  string  false  "columns to display"
// @Produce json/text
// @Success 200 {object} IndexListResponse
// @Router /es/_cat/aliases/{target} [get]
func CatAliases(c *gin.Context) {
	format := c.DefaultQuery("format", "table")
	verbose, _ := strconv.ParseBool(c.DefaultQuery("v", "false"))
	columns := strings.Split(c.DefaultQuery("h", strings.Join(catAliasesColumns, ",")), ",")

	for _, column := range columns {
		if !slices.Contains(catAliasesColumns, column) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid column: %s", column)})
			return
		}
	}

	if format != "json" && format != "table" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid format: %s", format)})
		return
	}

	targetAlias := c.Param("target")
	aliases := core.ZINC_INDEX_ALIAS_LIST.Aliases

	var response []CatAliasesResponseItem
	for alias, indexes := range aliases {
		if len(targetAlias) > 0 {
			if match, err := filepath.Match(targetAlias, alias); !match || err != nil {
				continue
			}
		}

		for _, index := range indexes {
			response = append(response, CatAliasesResponseItem{
				Alias:         alias,
				Index:         index,
				Filter:        "-",
				RoutingIndex:  "-",
				RoutingSearch: "-",
				IsWriteIndex:  "-",
			})
		}
	}

	if format == "json" {
		c.JSON(http.StatusOK, response)
	} else {
		c.String(http.StatusOK, buildAliasesTable(response, columns, verbose))
	}
}

func buildAliasesTable(data []CatAliasesResponseItem, columns []string, verbose bool) string {
	var width = map[string]int{}
	for _, column := range columns {
		width[column] = len(column)
	}
	for _, item := range data {
		for _, column := range columns {
			switch column {
			case "alias":
				if len(item.Alias) > width[column] {
					width[column] = len(item.Alias)
				}
			case "index":
				if len(item.Index) > width[column] {
					width[column] = len(item.Index)
				}
			}
		}
	}

	var res []string

	if verbose {
		line := []string{}
		for _, column := range columns {
			line = append(line, fmt.Sprintf("%s%s", column, strings.Repeat(" ", width[column]-len(column))))
		}
		res = append(res, strings.Join(line, " "))
	}
	for _, item := range data {
		line := []string{}
		for _, column := range columns {
			switch column {
			case "alias":
				line = append(line, fmt.Sprintf("%s%s", item.Alias, strings.Repeat(" ", width[column]-len(item.Alias))))
			case "index":
				line = append(line, fmt.Sprintf("%s%s", item.Index, strings.Repeat(" ", width[column]-len(item.Index))))
			case "filter":
				line = append(line, fmt.Sprintf("%s%s", item.Filter, strings.Repeat(" ", width[column]-len(item.Filter))))
			case "routing.index":
				line = append(line, fmt.Sprintf("%s%s", item.RoutingIndex, strings.Repeat(" ", width[column]-len(item.RoutingIndex))))
			case "routing.search":
				line = append(line, fmt.Sprintf("%s%s", item.RoutingSearch, strings.Repeat(" ", width[column]-len(item.RoutingSearch))))
			case "is_write_index":
				line = append(line, fmt.Sprintf("%s%s", item.IsWriteIndex, strings.Repeat(" ", width[column]-len(item.IsWriteIndex))))
			}
		}
		res = append(res, strings.Join(line, " "))
	}
	return strings.Join(res, "\n")
}

var catAliasesColumns = []string{
	"alias",
	"index",
	"filter",
	"routing.index",
	"routing.search",
	"is_write_index",
}

type CatAliasesResponseItem struct {
	Alias         string `json:"alias"`
	Index         string `json:"index"`
	Filter        string `json:"filter"`
	RoutingIndex  string `json:"routing.index"`
	RoutingSearch string `json:"routing.search"`
	IsWriteIndex  string `json:"is_write_index"`
}

func formatBytes(bytes int) string {
	const (
		_  = iota
		kb = 1 << (10 * iota)
		mb
		gb
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1fgb", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1fmb", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1fkb", float64(bytes)/float64(kb))
	}
	return fmt.Sprintf("%db", bytes)
}
