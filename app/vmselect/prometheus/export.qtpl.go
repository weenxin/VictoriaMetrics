// Code generated by qtc from "export.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line app/vmselect/prometheus/export.qtpl:1
package prometheus

//line app/vmselect/prometheus/export.qtpl:1
import (
	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmselect/netstorage"
	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmselect/netstorage/clickhouse"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/storage"
	"github.com/valyala/quicktemplate"
)

//line app/vmselect/prometheus/export.qtpl:9
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line app/vmselect/prometheus/export.qtpl:9
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line app/vmselect/prometheus/export.qtpl:9
func StreamExportPrometheusLine(qw422016 *qt422016.Writer, rs *netstorage.Result) {
//line app/vmselect/prometheus/export.qtpl:10
	if len(rs.Timestamps) == 0 {
//line app/vmselect/prometheus/export.qtpl:10
		return
//line app/vmselect/prometheus/export.qtpl:10
	}
//line app/vmselect/prometheus/export.qtpl:11
	bb := quicktemplate.AcquireByteBuffer()

//line app/vmselect/prometheus/export.qtpl:12
	writeprometheusMetricName(bb, &rs.MetricName)

//line app/vmselect/prometheus/export.qtpl:13
	for i, ts := range rs.Timestamps {
//line app/vmselect/prometheus/export.qtpl:14
		qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:14
		qw422016.N().S(` `)
//line app/vmselect/prometheus/export.qtpl:15
		qw422016.N().F(rs.Values[i])
//line app/vmselect/prometheus/export.qtpl:15
		qw422016.N().S(` `)
//line app/vmselect/prometheus/export.qtpl:16
		qw422016.N().DL(ts)
//line app/vmselect/prometheus/export.qtpl:16
		qw422016.N().S(`
`)
//line app/vmselect/prometheus/export.qtpl:17
	}
//line app/vmselect/prometheus/export.qtpl:18
	quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:19
}

//line app/vmselect/prometheus/export.qtpl:19
func WriteExportPrometheusLine(qq422016 qtio422016.Writer, rs *clickhouse.Result) {
//line app/vmselect/prometheus/export.qtpl:19
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:19
	StreamExportPrometheusLine(qw422016, rs)
//line app/vmselect/prometheus/export.qtpl:19
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:19
}

//line app/vmselect/prometheus/export.qtpl:19
func ExportPrometheusLine(rs *clickhouse.Result) string {
//line app/vmselect/prometheus/export.qtpl:19
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:19
	WriteExportPrometheusLine(qb422016, rs)
//line app/vmselect/prometheus/export.qtpl:19
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:19
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:19
	return qs422016
//line app/vmselect/prometheus/export.qtpl:19
}

//line app/vmselect/prometheus/export.qtpl:21
func StreamExportJSONLine(qw422016 *qt422016.Writer, rs *netstorage.Result) {
//line app/vmselect/prometheus/export.qtpl:22
	if len(rs.Timestamps) == 0 {
//line app/vmselect/prometheus/export.qtpl:22
		return
//line app/vmselect/prometheus/export.qtpl:22
	}
//line app/vmselect/prometheus/export.qtpl:22
	qw422016.N().S(`{"metric":`)
//line app/vmselect/prometheus/export.qtpl:24
	streammetricNameObject(qw422016, &rs.MetricName)
//line app/vmselect/prometheus/export.qtpl:24
	qw422016.N().S(`,"values":[`)
//line app/vmselect/prometheus/export.qtpl:26
	if len(rs.Values) > 0 {
//line app/vmselect/prometheus/export.qtpl:27
		values := rs.Values

//line app/vmselect/prometheus/export.qtpl:28
		qw422016.N().F(values[0])
//line app/vmselect/prometheus/export.qtpl:29
		values = values[1:]

//line app/vmselect/prometheus/export.qtpl:30
		for _, v := range values {
//line app/vmselect/prometheus/export.qtpl:30
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:31
			qw422016.N().F(v)
//line app/vmselect/prometheus/export.qtpl:32
		}
//line app/vmselect/prometheus/export.qtpl:33
	}
//line app/vmselect/prometheus/export.qtpl:33
	qw422016.N().S(`],"timestamps":[`)
//line app/vmselect/prometheus/export.qtpl:36
	if len(rs.Timestamps) > 0 {
//line app/vmselect/prometheus/export.qtpl:37
		timestamps := rs.Timestamps

//line app/vmselect/prometheus/export.qtpl:38
		qw422016.N().DL(timestamps[0])
//line app/vmselect/prometheus/export.qtpl:39
		timestamps = timestamps[1:]

//line app/vmselect/prometheus/export.qtpl:40
		for _, ts := range timestamps {
//line app/vmselect/prometheus/export.qtpl:40
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:41
			qw422016.N().DL(ts)
//line app/vmselect/prometheus/export.qtpl:42
		}
//line app/vmselect/prometheus/export.qtpl:43
	}
//line app/vmselect/prometheus/export.qtpl:43
	qw422016.N().S(`]}`)
//line app/vmselect/prometheus/export.qtpl:45
	qw422016.N().S(`
`)
//line app/vmselect/prometheus/export.qtpl:46
}

//line app/vmselect/prometheus/export.qtpl:46
func WriteExportJSONLine(qq422016 qtio422016.Writer, rs *clickhouse.Result) {
//line app/vmselect/prometheus/export.qtpl:46
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:46
	StreamExportJSONLine(qw422016, rs)
//line app/vmselect/prometheus/export.qtpl:46
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:46
}

//line app/vmselect/prometheus/export.qtpl:46
func ExportJSONLine(rs *netstorage.Result) string {
//line app/vmselect/prometheus/export.qtpl:46
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:46
	WriteExportJSONLine(qb422016, rs)
//line app/vmselect/prometheus/export.qtpl:46
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:46
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:46
	return qs422016
//line app/vmselect/prometheus/export.qtpl:46
}

//line app/vmselect/prometheus/export.qtpl:48
func StreamExportPromAPILine(qw422016 *qt422016.Writer, rs *netstorage.Result) {
//line app/vmselect/prometheus/export.qtpl:48
	qw422016.N().S(`{"metric":`)
//line app/vmselect/prometheus/export.qtpl:50
	streammetricNameObject(qw422016, &rs.MetricName)
//line app/vmselect/prometheus/export.qtpl:50
	qw422016.N().S(`,"values":`)
//line app/vmselect/prometheus/export.qtpl:51
	streamvaluesWithTimestamps(qw422016, rs.Values, rs.Timestamps)
//line app/vmselect/prometheus/export.qtpl:51
	qw422016.N().S(`}`)
//line app/vmselect/prometheus/export.qtpl:53
}

//line app/vmselect/prometheus/export.qtpl:53
func WriteExportPromAPILine(qq422016 qtio422016.Writer, rs *clickhouse.Result) {
//line app/vmselect/prometheus/export.qtpl:53
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:53
	StreamExportPromAPILine(qw422016, rs)
//line app/vmselect/prometheus/export.qtpl:53
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:53
}

//line app/vmselect/prometheus/export.qtpl:53
func ExportPromAPILine(rs *netstorage.Result) string {
//line app/vmselect/prometheus/export.qtpl:53
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:53
	WriteExportPromAPILine(qb422016, rs)
//line app/vmselect/prometheus/export.qtpl:53
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:53
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:53
	return qs422016
//line app/vmselect/prometheus/export.qtpl:53
}

//line app/vmselect/prometheus/export.qtpl:55
func StreamExportPromAPIResponse(qw422016 *qt422016.Writer, resultsCh <-chan *quicktemplate.ByteBuffer) {
//line app/vmselect/prometheus/export.qtpl:55
	qw422016.N().S(`{"status":"success","data":{"resultType":"matrix","result":[`)
//line app/vmselect/prometheus/export.qtpl:61
	bb, ok := <-resultsCh

//line app/vmselect/prometheus/export.qtpl:62
	if ok {
//line app/vmselect/prometheus/export.qtpl:63
		qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:64
		quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:65
		for bb := range resultsCh {
//line app/vmselect/prometheus/export.qtpl:65
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:66
			qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:67
			quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:68
		}
//line app/vmselect/prometheus/export.qtpl:69
	}
//line app/vmselect/prometheus/export.qtpl:69
	qw422016.N().S(`]}}`)
//line app/vmselect/prometheus/export.qtpl:73
}

//line app/vmselect/prometheus/export.qtpl:73
func WriteExportPromAPIResponse(qq422016 qtio422016.Writer, resultsCh <-chan *quicktemplate.ByteBuffer) {
//line app/vmselect/prometheus/export.qtpl:73
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:73
	StreamExportPromAPIResponse(qw422016, resultsCh)
//line app/vmselect/prometheus/export.qtpl:73
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:73
}

//line app/vmselect/prometheus/export.qtpl:73
func ExportPromAPIResponse(resultsCh <-chan *quicktemplate.ByteBuffer) string {
//line app/vmselect/prometheus/export.qtpl:73
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:73
	WriteExportPromAPIResponse(qb422016, resultsCh)
//line app/vmselect/prometheus/export.qtpl:73
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:73
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:73
	return qs422016
//line app/vmselect/prometheus/export.qtpl:73
}

//line app/vmselect/prometheus/export.qtpl:75
func StreamExportStdResponse(qw422016 *qt422016.Writer, resultsCh <-chan *quicktemplate.ByteBuffer) {
//line app/vmselect/prometheus/export.qtpl:76
	for bb := range resultsCh {
//line app/vmselect/prometheus/export.qtpl:77
		qw422016.N().Z(bb.B)
//line app/vmselect/prometheus/export.qtpl:78
		quicktemplate.ReleaseByteBuffer(bb)

//line app/vmselect/prometheus/export.qtpl:79
	}
//line app/vmselect/prometheus/export.qtpl:80
}

//line app/vmselect/prometheus/export.qtpl:80
func WriteExportStdResponse(qq422016 qtio422016.Writer, resultsCh <-chan *quicktemplate.ByteBuffer) {
//line app/vmselect/prometheus/export.qtpl:80
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:80
	StreamExportStdResponse(qw422016, resultsCh)
//line app/vmselect/prometheus/export.qtpl:80
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:80
}

//line app/vmselect/prometheus/export.qtpl:80
func ExportStdResponse(resultsCh <-chan *quicktemplate.ByteBuffer) string {
//line app/vmselect/prometheus/export.qtpl:80
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:80
	WriteExportStdResponse(qb422016, resultsCh)
//line app/vmselect/prometheus/export.qtpl:80
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:80
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:80
	return qs422016
//line app/vmselect/prometheus/export.qtpl:80
}

//line app/vmselect/prometheus/export.qtpl:82
func streamprometheusMetricName(qw422016 *qt422016.Writer, mn *storage.MetricName) {
//line app/vmselect/prometheus/export.qtpl:83
	qw422016.N().Z(mn.MetricGroup)
//line app/vmselect/prometheus/export.qtpl:84
	if len(mn.Tags) > 0 {
//line app/vmselect/prometheus/export.qtpl:84
		qw422016.N().S(`{`)
//line app/vmselect/prometheus/export.qtpl:86
		tags := mn.Tags

//line app/vmselect/prometheus/export.qtpl:87
		qw422016.N().Z(tags[0].Key)
//line app/vmselect/prometheus/export.qtpl:87
		qw422016.N().S(`=`)
//line app/vmselect/prometheus/export.qtpl:87
		qw422016.N().QZ(tags[0].Value)
//line app/vmselect/prometheus/export.qtpl:88
		tags = tags[1:]

//line app/vmselect/prometheus/export.qtpl:89
		for i := range tags {
//line app/vmselect/prometheus/export.qtpl:90
			tag := &tags[i]

//line app/vmselect/prometheus/export.qtpl:90
			qw422016.N().S(`,`)
//line app/vmselect/prometheus/export.qtpl:91
			qw422016.N().Z(tag.Key)
//line app/vmselect/prometheus/export.qtpl:91
			qw422016.N().S(`=`)
//line app/vmselect/prometheus/export.qtpl:91
			qw422016.N().QZ(tag.Value)
//line app/vmselect/prometheus/export.qtpl:92
		}
//line app/vmselect/prometheus/export.qtpl:92
		qw422016.N().S(`}`)
//line app/vmselect/prometheus/export.qtpl:94
	}
//line app/vmselect/prometheus/export.qtpl:95
}

//line app/vmselect/prometheus/export.qtpl:95
func writeprometheusMetricName(qq422016 qtio422016.Writer, mn *storage.MetricName) {
//line app/vmselect/prometheus/export.qtpl:95
	qw422016 := qt422016.AcquireWriter(qq422016)
//line app/vmselect/prometheus/export.qtpl:95
	streamprometheusMetricName(qw422016, mn)
//line app/vmselect/prometheus/export.qtpl:95
	qt422016.ReleaseWriter(qw422016)
//line app/vmselect/prometheus/export.qtpl:95
}

//line app/vmselect/prometheus/export.qtpl:95
func prometheusMetricName(mn *storage.MetricName) string {
//line app/vmselect/prometheus/export.qtpl:95
	qb422016 := qt422016.AcquireByteBuffer()
//line app/vmselect/prometheus/export.qtpl:95
	writeprometheusMetricName(qb422016, mn)
//line app/vmselect/prometheus/export.qtpl:95
	qs422016 := string(qb422016.B)
//line app/vmselect/prometheus/export.qtpl:95
	qt422016.ReleaseByteBuffer(qb422016)
//line app/vmselect/prometheus/export.qtpl:95
	return qs422016
//line app/vmselect/prometheus/export.qtpl:95
}
