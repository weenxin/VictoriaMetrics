package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmselect/netstorage"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/storage"
	"log"
	"strings"
	"time"
)


var (
	connect *sql.DB = nil
	ErrorNotFoundMetricName = errors.New("does not find a metric name in filter")
	ErrorMetricsNameNotMatch = errors.New("metric name should be {table}.{aggeration}.{field}, ex: perflog.sum.count")
	ErrorUnsportedAggerations = errors.New("unsupported aggeration function, " +
		"current: "+ strings.Join(aggerationFuctions,",") )

	aggerationFuctions = []string{AggCount,AggSum}
)

const (
	TokenTable = iota
	TokenAaggregation
	TokenField
	TokenCount

	AggCount = "count"
	AggSum = "sum"
)



func InitDatabase(){
	var err error
	connect, err = newClickhouseClient()
	if err != nil {
		panic("initialize clickhouse client error : %v"+ err.Error())
	}
}

func newClickhouseClient() (*sql.DB, error) {
	if *address == "" {
		return nil, ErrHostEmpty
	}
	url := fmt.Sprintf("tcp://%v?username=%v&password=%v&database=%v", *address, *clickhouseUserName, *clickhousePassword, *clickhouseDatabase)

	connect, err := sql.Open("clickhouse", url)
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return nil,err
	}

	return connect, nil
}

func getMetricName(sq *storage.SearchQuery) (string,error){
	for _, tags := range sq.TagFilterss {
		for _,tag := range tags{
			if string(tag.Key) == "__name__" || string(tag.Key) == ""{
				return string(tag.Value),nil
			}
		}
	}

	logger.Infof("does not find a __name__ fileter: %v",sq.String())
	return "",ErrorNotFoundMetricName
}



func getTokens(metricName string) ( table, aggregation, filed string, err error){
	tokens := strings.SplitN(metricName,".",3)
	if len(tokens) < TokenCount{
		return "","","",ErrorMetricsNameNotMatch
	}
	table = tokens[TokenTable]
	aggregation = tokens[TokenAaggregation]
	filed = tokens[TokenField]

	found := false;
	for _,agg := range aggerationFuctions {
		if agg == aggregation{
			found = true;
			break
		}
	}

	if !found{
		return "","","",ErrorUnsportedAggerations
	}

	return table,aggregation,filed,nil

}

//获取当前需要聚合的维度信息，
func getQFieldsuery(filed,aggregation string ) string {
	fileds := []string{*clickhouseTimesampFiled}
	fileds = append(fileds,*clickhouseTags...)

	aggField := ""
	switch aggregation {
	case AggCount:
		aggField = fmt.Sprintf("count( %v )",filed)
	case AggSum:
		aggField = fmt.Sprintf("sum( %v )",filed)
	default:
		panic("unknown aggration")//should return a error before call this funtion
	}

	fileds = append(fileds,aggField)

	return strings.Join(fileds,",")
}


func getFiledOperationQuery(filter storage.TagFilter)string{
	filed := string(filter.Key)
	operation := "="
	if filter.IsRegexp {
		operation = " LIKE "
		if filter.IsNegative {
			operation = " NOT LIKE "
		}
	}else {
		operation = " = "
		if filter.IsNegative {
			operation = " != "
		}
	}

	return fmt.Sprintf("%v %v '%v'",filed,operation,string(filter.Value))
}

func getGroupbyQuery() string {
	fileds := []string{*clickhouseTimesampFiled}
	fileds = append(fileds,*clickhouseTags...)


	return strings.Join(fileds,",")
}

func getWhereQuery(sq *storage.SearchQuery) string{
	timeQuery := fmt.Sprintf("%v > %v AND %v < %v",
		*clickhouseTimesampFiled,sq.MinTimestamp/1000, *clickhouseTimesampFiled,sq.MaxTimestamp/1000)

	fieldFilter := []string{timeQuery}
	for _,filters := range sq.TagFilterss {
		for _, filter := range filters{
			if string(filter.Key) != "" && string(filter.Key) != "__name__"{
				fieldFilter = append(fieldFilter,getFiledOperationQuery(filter))
			}

		}
	}

	if *clickhouseDateField != ""{
		fieldFilter = append(fieldFilter,fmt.Sprintf(" %v >= '%v' AND %v <= '%v' " ,
			*clickhouseDateField,
			time.Unix(sq.MinTimestamp/1000,0).Format("2006-01-02"),
			*clickhouseDateField,
			time.Unix(sq.MaxTimestamp/1000,0).Format("2006-01-02")))
	}

	return strings.Join(fieldFilter," AND ")
}


func formatQuery(table,aggregation,filed string, sq *storage.SearchQuery)string {

	sql :=  fmt.Sprintf("SELECT %v  FROM %v WHERE %v GROUP BY %v ORDER BY %v",
		getQFieldsuery(filed,aggregation),
		table,
		getWhereQuery(sq),
		getGroupbyQuery(),
		*clickhouseTimesampFiled,
		)

	return sql

}
func fomatMetricName(name string, columns []string, values []interface{}) storage.MetricName{
	var metricsName storage.MetricName
	metricsName.MetricGroup = ([]byte)(name)
	for i, column := range columns {
		metricsName.Tags = append(metricsName.Tags,storage.Tag{
			([]byte) (column), ([]byte)(*(values[i].(*string)))})
	}

	return metricsName
}

func metricNameEqueal(first *storage.MetricName, second *storage.MetricName)bool{
	if string(first.MetricGroup) != string(second.MetricGroup){
		return false
	}
	for i , _ := range first.Tags {
		if string(first.Tags[i].Key) != string(second.Tags[i].Key){
			return false
		}
		if string(first.Tags[i].Value) != string(second.Tags[i].Value){
			return false
		}
	}

	return true
}

func addMetrics(timeSeries []*Result,meticName string, columns []string,
	values []interface{}) []*Result{

	formatMetricName := fomatMetricName(meticName,columns[1:len(columns)-1],values[1:len(values)-1])
	timestamp := (values[0].(*time.Time)).Unix() * 1000
	value := *(values[len(values)-1].(*float64))


	found := false
	for i, timeseriel := range timeSeries{
		if metricNameEqueal(&timeseriel.MetricName,&formatMetricName){
			timeSeries[i].Timestamps = append(timeSeries[i].Timestamps,timestamp)
			timeSeries[i].Values = append(timeSeries[i].Values,value)
			found = true
		}
	}
	if !found{
		result := Result{}
		result.MetricName = formatMetricName
		result.Values = make([]float64,0,1024)
		result.Timestamps = make([]int64,0, 1024)
		result.Values = append(result.Values,value)
		result.Timestamps = append(result.Timestamps,timestamp)
		timeSeries = append(timeSeries,&result)
	}

	return timeSeries


}

func query(sq *storage.SearchQuery, fetchData bool, deadline netstorage.Deadline)([]*Result, error){
	if connect == nil {
		InitDatabase()
	}

	metricName, err := getMetricName(sq)
	if err != nil {
		return []*Result{},err
	}
	table,aggregation,filed, err := getTokens(metricName)
	if err != nil {
		return []*Result{},err
	}

	query := formatQuery(table,aggregation,filed,sq)
	fmt.Println(query)



	rows, err := connect.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return []*Result{},err
	}

	coloumnLen := len(columns)
	fileds := make([]interface{},0 ,coloumnLen)

	fileds = append(fileds, new(time.Time))
	for i := 1 ; i < coloumnLen-1 ; i ++ {
		fileds = append(fileds, new(string))
	}
	fileds = append(fileds,new(float64))

	var series []*Result

	for rows.Next() {

		if err := rows.Scan(fileds...); err != nil {
			log.Fatal(err)
		}else{
			series = addMetrics(series,metricName,columns,fileds)
		}
	}

	if err := rows.Err(); err != nil {
		return []*Result{}, err
	}

	return series , nil

}