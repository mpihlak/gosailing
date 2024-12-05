package datasource

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type NavigationDataPoint struct {
	Timestamp          time.Time
	Heading            float64
	SpeedThroughWater  float64
	TrueWindAngle      float64
	TrueWindSpeed      float64
	TrueWindDirection  float64
	Latitude           float64
	Longitude          float64
	CumulativeDistance float64
}

type NavigationDataProvider interface {
	Next() (NavigationDataPoint, bool)
}

type ReplayNavigationDataProvider struct {
	fieldMap map[string]int
	records  [][]string
	pos      int
}

func NewReplayNavigationDataProvider(reader io.Reader) (*ReplayNavigationDataProvider, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("no navigation data")
	}

	fieldMap := make(map[string]int)
	for pos, key := range records[0] {
		fieldMap[key] = pos
	}

	return &ReplayNavigationDataProvider{
		fieldMap: fieldMap,
		records:  records[1:],
	}, nil
}

func (r *ReplayNavigationDataProvider) assignFieldValue(record []string, key string, target *float64) bool {
	fieldPos, ok := r.fieldMap[key]
	if !ok {
		return false
	}

	val, err := strconv.ParseFloat(record[fieldPos], 64)
	if err != nil {
		return false
	}

	*target = val

	return true
}

func (r *ReplayNavigationDataProvider) Next() (NavigationDataPoint, bool) {
	var result NavigationDataPoint
	if r.pos >= len(r.records) {
		return result, false
	}

	fieldPos, ok := r.fieldMap["time"]
	if !ok {
		return result, false
	}
	record := r.records[r.pos]
	if fieldPos >= len(record) {
		return result, false
	}
	parsedTime, err := time.Parse(time.RFC3339, record[fieldPos])
	if err != nil {
		return result, false
	}

	result.Timestamp = parsedTime

	// Example data:
	// time,dpt,mtw,awa,aws,cog,hdg,rot,sog,stw,pitch,yaw,roll,lng,lat,tws,twa,twd,vmg,dist,cum_dist
	// 2024-09-11T17:27:52+03:00,22.5,288.65,31.999692858056477,17.5932,87.18125810710607,78.09987705428252,1.5999846429028237,5.851439999999999,5.50152,-3.0997016716577535,-101.79741146089336,-20.597832734953094,24.798672,59.488683,13.00591745831071,45.792515253115,132.97377336022106,4.079967736770225,0.0015914414543093099,0.0015914414543093099
	//

	ok = r.assignFieldValue(record, "hdg", &result.Heading) &&
		r.assignFieldValue(record, "stw", &result.SpeedThroughWater) &&
		r.assignFieldValue(record, "twa", &result.TrueWindAngle) &&
		r.assignFieldValue(record, "tws", &result.TrueWindSpeed) &&
		r.assignFieldValue(record, "twd", &result.TrueWindDirection) &&
		r.assignFieldValue(record, "lat", &result.Latitude) &&
		r.assignFieldValue(record, "lng", &result.Longitude) &&
		r.assignFieldValue(record, "cum_dist", &result.CumulativeDistance)

	r.pos++

	return result, ok
}
