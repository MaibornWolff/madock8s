package pkg

import (
	"regexp"
	"strings"
)

func ParsePrometheusResult(prometheusresult string) []MetricsMetaBlock {
	results := strings.Split(prometheusresult, "\n")
	helpentries := findEntriesByExpression(results, "#\\sHELP\\s(.+?)\\s(.+)$")
	typeentries := findEntriesByExpression(results, "#\\sTYPE\\s(.+?)\\s(.+?)$")
	blocks := createMetricMetaBlocks(helpentries, typeentries)
	return blocks
}

func createMetricMetaBlocks(helpentries map[string]string, typeentries map[string]string) []MetricsMetaBlock {
	var blocks []MetricsMetaBlock

	for key, helptext := range helpentries {
		mtype := typeentries[key]

		blocks = append(blocks, MetricsMetaBlock{
			Name: key,
			Help: helptext,
			Type: mtype,
		})
	}

	return blocks
}

func findEntriesByExpression(prometheusresults []string, exp string) map[string]string {
	entries := make(map[string]string, 0)

	ex := regexp.MustCompile(exp)

	for _, result := range prometheusresults {
		if ex.MatchString(result) {
			match := ex.FindStringSubmatch(result)
			entries[match[1]] = match[2]
		}
	}

	return entries
}

type MetricsMetaBlock struct {
	Name string
	Help string
	Type string
}
