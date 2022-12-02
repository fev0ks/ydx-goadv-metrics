package pages

import (
	"sort"
	"strconv"
	"strings"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

var (
	metricsTableTitles = []string{"#", "type", "name", "value"}
)

func GetMetricsPage(metrics map[string]*model.Metric) string {
	hb := GetHTMLBuilder()
	hb.
		Add(TableStile).
		Add(OBody).
		//AddHeader(fmt.Sprintf("Metrics result at %v", time.Now().Format(time.RFC1123))).
		AddHeader("Metrics result").
		Add(OTable)
	addTableTitles(hb, metricsTableTitles)
	metricsList := convertToList(metrics)
	sortMetrics(metricsList)
	for i, metric := range metricsList {
		addMetricDataLine(hb, i+1, metric)
	}
	page := hb.
		Add(CTable).
		Add(CBody).
		GetHTMLPage()
	return page
}

func addTableTitles(hb *htmlBuilder, titles []string) {
	hb.Add(OLine)
	for _, title := range titles {
		hb.Add(OTitle)
		hb.Add(title)
		hb.Add(CTitle)
	}
	hb.Add(CLine)
}

func addMetricDataLine(hb *htmlBuilder, number int, metric *model.Metric) {
	hb.
		Add(OLine).
		Add(OColumn).
		Add(strconv.Itoa(number)).
		Add(CColumn).
		Add(OColumn).
		Add(string(metric.MType)).
		Add(CColumn).
		Add(OColumn).
		Add(metric.Name).
		Add(CColumn).Add(OColumn).
		Add(metric.GetValue()).
		Add(CColumn).
		Add(CLine)
}

func sortMetrics(metrics []*model.Metric) {
	sort.SliceStable(metrics, func(i, j int) bool {
		return strings.Compare(string(metrics[i].MType), string(metrics[j].MType)) < 0 ||
			strings.Compare(string(metrics[i].MType), string(metrics[j].MType)) == 0 &&
				strings.Compare(metrics[i].Name, metrics[j].Name) < 0
	})
}

func convertToList(metrics map[string]*model.Metric) []*model.Metric {
	list := make([]*model.Metric, 0, len(metrics))
	for _, value := range metrics {
		list = append(list, value)
	}
	return list
}
