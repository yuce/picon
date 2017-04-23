package picon

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	pilosa "github.com/pilosa/go-client-pilosa"
)

type promptInfo struct {
	address string
	index   string
}

func printResponse(response *pilosa.QueryResponse) {
	if !response.Success {
		printError(errors.New(response.ErrorMessage))
		return
	}
	results := response.Results()
	for i, result := range results {
		printResult(i, len(results), result)
	}
}

func printError(err error) {
	fmt.Println(colorString(fgRed, err.Error()))
}

func printWarning(msg string) {
	fmt.Println(colorString(fgRed, msg))
}

func colorString(color Ansi, msg string) string {
	return fmt.Sprintf("%s%s%s", color, msg, attrReset)
}

func printResult(index int, count int, result *pilosa.QueryResult) {
	headerFmt := fmt.Sprintf("[%%%dd] --------", int(math.Ceil(float64(count)/10.0)))
	lines := []string{fmt.Sprintf(headerFmt, index)}
	canPrint := false
	switch {
	case result.Bitmap != nil:
		if len(attributesToString(result.Bitmap.Attributes)) > 0 {
			lines = append(lines,
				fmt.Sprintf("\tAttributes: %s", attributesToString(result.Bitmap.Attributes)))
			canPrint = true
		}
		if len(bitsToString(result.Bitmap.Bits)) > 0 {
			lines = append(lines,
				fmt.Sprintf("\tBits      : %s", bitsToString(result.Bitmap.Bits)))
			canPrint = true
		}
	case result.CountItems != nil && len(result.CountItems) > 0:
		for _, item := range result.CountItems {
			lines = append(lines, fmt.Sprintf("\tCount(%d) = %d\n", item.ID, item.Count))
			canPrint = true
		}
	case result.Count > 0:
		lines = append(lines, fmt.Sprintf("\tCount: %d\n", result.Count))
		canPrint = true
	}
	if canPrint {
		fmt.Println(strings.Join(lines, "\n"))
	}
}

func attributesToString(attrs map[string]interface{}) string {
	parts := make([]string, 0, len(attrs))
	for k, v := range attrs {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ", ")
}

func bitsToString(bits []uint64) string {
	parts := make([]string, 0, len(bits))
	for _, v := range bits {
		parts = append(parts, strconv.Itoa(int(v)))
	}
	return strings.Join(parts, ", ")
}

func autoSessionName() string {
	return time.Now().Format("2006-01-02_15-04-05")
}
