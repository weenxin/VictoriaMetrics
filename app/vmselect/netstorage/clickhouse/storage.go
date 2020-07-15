package clickhouse

import (
	"errors"
	"flag"
	"fmt"
	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmselect/netstorage"
	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmstorage"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/flagutil"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/storage"
	"github.com/VictoriaMetrics/metrics"
	"runtime"
	"sort"
	"sync"
)

var (
	maxTagKeysPerSearch   = flag.Int("clickhouse.maxTagKeys", 100e3,
		"The maximum number of tag keys returned per search")
	maxTagValuesPerSearch = flag.Int("search.maxTagValues", 100e3, "The maximum number of tag values returned per search")
	maxMetricsPerSearch   = flag.Int("search.maxUniqueTimeseries", 300e3, "The maximum number of unique time series each search can scan")

	address = flag.String("clickhouse.address","public-chdb183.idczw.hb1.kwaidc.com:9000","clickhouse address like ")
	clickhouseUserName = flag.String("clickhouse.user", "default", "clickhouse cluster user name ")
	clickhousePassword = flag.String("clickhouse.password", "", "clickhouse cluster password ")
	clickhouseDatabase = flag.String("clickhouse.database", "openfalcon", "clickhouse database to access")
	clickhouseTimesampFiled = flag.String("clickhouse.timestamp", "timestamp", "timestamp column name in table")
 	clickhouseTags = flagutil.NewArray("clickhouse.tag", "fields to keep --clickhouse.tag=subtag --clickhouse." +
 		"tag=extra1")

	clickhoseFileds = flagutil.NewArray("clickhouse.filed", "fields to keep --clickhouse.filed=count --clickhouse." +
		"tag=sum")

)

var (
	ErrHostEmpty = errors.New("clickhouse host can not be empty")
)


// Result is a single timeseries result.
//
// ProcessSearchQuery returns Result slice.
type Result struct {
	// The name of the metric.
	MetricName storage.MetricName

	// Values are sorted by Timestamps.
	Values     []float64
	Timestamps []int64

	// Marshaled MetricName. Used only for results sorting
	// in app/vmselect/promql
	MetricNameMarshaled []byte
}

func (r *Result) reset() {
	r.MetricName.Reset()
	r.Values = r.Values[:0]
	r.Timestamps = r.Timestamps[:0]
	r.MetricNameMarshaled = r.MetricNameMarshaled[:0]
}

// Results holds results returned from ProcessSearchQuery.
type Results struct {
	tr        storage.TimeRange
	fetchData bool
	deadline  netstorage.Deadline

	results []*Result
	//sr               *storage.Search
}

// Len returns the number of results in rss.
func (rss *Results) Len() int {
	return len(rss.results)
}

// Cancel cancels rss work.
func (rss *Results) Cancel() {
	rss.mustClose()
}

func (rss *Results) mustClose() {
	//putStorageSearch(rss.sr)
	//rss.sr = nil
}


// RunParallel runs in parallel f for all the results from rss.
//
// f shouldn't hold references to rs after returning.
// workerID is the id of the worker goroutine that calls f.
//
// rss becomes unusable after the call to RunParallel.
func (rss *Results) RunParallel(f func(rs *Result, workerID uint)) error {
	defer rss.mustClose()

	for _, rs := range rss.results{
		f(rs,0)
	}

	return nil
}

var perQueryRowsProcessed = metrics.NewHistogram(`vm_per_query_rows_processed_count`)
var perQuerySeriesProcessed = metrics.NewHistogram(`vm_per_query_series_processed_count`)

var gomaxprocs = runtime.GOMAXPROCS(-1)

type packedTimeseries struct {
	metricName string

	// Values are sorted by Timestamps.
	Values     []float64
	Timestamps []int64

	//brs        []storage.BlockRef
}


var dedupsDuringSelect = metrics.NewCounter(`vm_deduplicated_samples_total{type="select"}`)


// DeleteSeries deletes time series matching the given tagFilterss.
func DeleteSeries(sq *storage.SearchQuery) (int, error) {
	tfss, err := setupTfss(sq.TagFilterss)
	if err != nil {
		return 0, err
	}
	return vmstorage.DeleteMetrics(tfss)
}

// GetLabels returns labels until the given deadline.
func GetLabels(deadline netstorage.Deadline) ([]string, error) {
	labels, err := vmstorage.SearchTagKeys(*maxTagKeysPerSearch)
	if err != nil {
		return nil, fmt.Errorf("error during labels search: %w", err)
	}

	// Substitute "" with "__name__"
	for i := range labels {
		if labels[i] == "" {
			labels[i] = "__name__"
		}
	}

	// Sort labels like Prometheus does
	sort.Strings(labels)

	return labels, nil
}

// GetLabelValues returns label values for the given labelName
// until the given deadline.
func GetLabelValues(labelName string, deadline netstorage.Deadline) ([]string, error) {
	if labelName == "__name__" {
		labelName = ""
	}

	// Search for tag values
	labelValues, err := vmstorage.SearchTagValues([]byte(labelName), *maxTagValuesPerSearch)
	if err != nil {
		return nil, fmt.Errorf("error during label values search for labelName=%q: %w", labelName, err)
	}

	// Sort labelValues like Prometheus does
	sort.Strings(labelValues)

	return labelValues, nil
}

// GetLabelEntries returns all the label entries until the given deadline.
func GetLabelEntries(deadline netstorage.Deadline) ([]storage.TagEntry, error) {
	labelEntries, err := vmstorage.SearchTagEntries(*maxTagKeysPerSearch, *maxTagValuesPerSearch)
	if err != nil {
		return nil, fmt.Errorf("error during label entries request: %w", err)
	}

	// Substitute "" with "__name__"
	for i := range labelEntries {
		e := &labelEntries[i]
		if e.Key == "" {
			e.Key = "__name__"
		}
	}

	// Sort labelEntries by the number of label values in each entry.
	sort.Slice(labelEntries, func(i, j int) bool {
		a, b := labelEntries[i].Values, labelEntries[j].Values
		if len(a) != len(b) {
			return len(a) > len(b)
		}
		return labelEntries[i].Key > labelEntries[j].Key
	})

	return labelEntries, nil
}

// GetTSDBStatusForDate returns tsdb status according to https://prometheus.io/docs/prometheus/latest/querying/api/#tsdb-stats
func GetTSDBStatusForDate(deadline netstorage.Deadline, date uint64, topN int) (*storage.TSDBStatus, error) {
	status, err := vmstorage.GetTSDBStatusForDate(date, topN)
	if err != nil {
		return nil, fmt.Errorf("error during tsdb status request: %w", err)
	}
	return status, nil
}

// GetSeriesCount returns the number of unique series.
func GetSeriesCount(deadline netstorage.Deadline) (uint64, error) {
	n, err := vmstorage.GetSeriesCount()
	if err != nil {
		return 0, fmt.Errorf("error during series count request: %w", err)
	}
	return n, nil
}

func getStorageSearch() *storage.Search {
	v := ssPool.Get()
	if v == nil {
		return &storage.Search{}
	}
	return v.(*storage.Search)
}

func putStorageSearch(sr *storage.Search) {
	sr.MustClose()
	ssPool.Put(sr)
}

var ssPool sync.Pool

// ProcessSearchQuery performs sq on storage nodes until the given deadline.
//
// Results.RunParallel or Results.Cancel must be called on the returned Results.
func ProcessSearchQuery(sq *storage.SearchQuery, fetchData bool, deadline netstorage.Deadline) (*Results, error) {
	// Setup search.
	tr := storage.TimeRange{
		MinTimestamp: sq.MinTimestamp,
		MaxTimestamp: sq.MaxTimestamp,
	}
	if err := vmstorage.CheckTimeRange(tr); err != nil {
		return nil, err
	}
	var series []*Result
	if fetchData{
		var err error
		series,err = query(sq,fetchData,deadline)
		if err != nil {
			return nil,err
		}
	}

	var rss Results
	rss.tr = tr
	rss.fetchData = fetchData
	rss.deadline = deadline
	rss.results = series
	return &rss, nil
}

func setupTfss(tagFilterss [][]storage.TagFilter) ([]*storage.TagFilters, error) {
	tfss := make([]*storage.TagFilters, 0, len(tagFilterss))
	for _, tagFilters := range tagFilterss {
		tfs := storage.NewTagFilters()
		for i := range tagFilters {
			tf := &tagFilters[i]
			if err := tfs.Add(tf.Key, tf.Value, tf.IsNegative, tf.IsRegexp); err != nil {
				return nil, fmt.Errorf("cannot parse tag filter %s: %w", tf, err)
			}
		}
		tfss = append(tfss, tfs)
		tfss = append(tfss, tfs.Finalize()...)
	}
	return tfss, nil
}


