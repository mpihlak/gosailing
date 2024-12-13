package datasource

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sort"
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
	CourseOverGround   float64
	SpeedOverGround    float64
	ApparentWindSpeed  float64
	ApparentWindAngle  float64
}

type NavigationDataProvider interface {
	Next() (NavigationDataPoint, bool)
}

type ReplayNavigationDataProvider struct {
	startTime *time.Time
	endTime   *time.Time
	fieldMap  map[string]int
	records   [][]string
	pos       int
}

func NewReplayNavigationDataProvider(reader io.Reader, startTime, endTime *time.Time) (*ReplayNavigationDataProvider, error) {
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
		startTime: startTime,
		endTime:   endTime,
		fieldMap:  fieldMap,
		records:   records[1:],
	}, nil
}

func (r *ReplayNavigationDataProvider) assignFieldValue(record []string, key string, target *float64) bool {
	fieldPos, ok := r.fieldMap[key]
	if !ok {
		log.Printf("field %s not found", key)
		return false
	}

	val, err := strconv.ParseFloat(record[fieldPos], 64)
	if err != nil {
		log.Printf("error parsing field %s: %v", key, err)
		return false
	}

	*target = val

	return true
}

func (r *ReplayNavigationDataProvider) isTimeInRange(timestamp time.Time) bool {
	if r.startTime != nil && timestamp.Before(*r.startTime) {
		return false
	}
	if r.endTime != nil && timestamp.After(*r.endTime) {
		return false
	}
	return true
}

func (r *ReplayNavigationDataProvider) Next() (NavigationDataPoint, bool) {
	timeFieldPos, ok := r.fieldMap["time"]
	if !ok {
		return NavigationDataPoint{}, false
	}

	for r.pos < len(r.records) {
		record := r.records[r.pos]
		r.pos++

		if timeFieldPos >= len(record) {
			log.Printf("timeFieldPos: %d, record: %v", timeFieldPos, record)
			return NavigationDataPoint{}, false
		}

		parsedTime, err := time.Parse(time.RFC3339, record[timeFieldPos])
		if err != nil {
			log.Printf("error parsing time: %v", err)
			return NavigationDataPoint{}, false
		}

		if !r.isTimeInRange(parsedTime) {
			continue
		}

		result := NavigationDataPoint{
			Timestamp: parsedTime,
		}

		// Example data:
		// time,dpt,mtw,awa,aws,cog,hdg,rot,sog,stw,pitch,yaw,roll,lng,lat,tws,twa,twd,vmg,dist,cum_dist
		// 2024-09-11T17:27:52+03:00,22.5,288.65,31.999692858056477,17.5932,87.18125810710607,78.09987705428252,1.5999846429028237,5.851439999999999,5.50152,-3.0997016716577535,-101.79741146089336,-20.597832734953094,24.798672,59.488683,13.00591745831071,45.792515253115,132.97377336022106,4.079967736770225,0.0015914414543093099,0.0015914414543093099
		//

		ok = r.assignFieldValue(record, "hdg", &result.Heading) &&
			r.assignFieldValue(record, "twa", &result.TrueWindAngle) &&
			r.assignFieldValue(record, "tws", &result.TrueWindSpeed) &&
			r.assignFieldValue(record, "twd", &result.TrueWindDirection) &&
			r.assignFieldValue(record, "aws", &result.ApparentWindSpeed) &&
			r.assignFieldValue(record, "awa", &result.ApparentWindAngle) &&
			r.assignFieldValue(record, "lat", &result.Latitude) &&
			r.assignFieldValue(record, "lng", &result.Longitude) &&
			r.assignFieldValue(record, "cum_dist", &result.CumulativeDistance) &&
			r.assignFieldValue(record, "cog", &result.CourseOverGround) &&
			r.assignFieldValue(record, "sog", &result.SpeedOverGround)

		// Ignore errors from sometimes missing stw field
		r.assignFieldValue(record, "stw", &result.SpeedThroughWater)

		if ok {
			return result, true
		}
	}

	return NavigationDataPoint{}, false
}

// GetAllPoints returns all navigation data points in the replay data
func (r *ReplayNavigationDataProvider) GetAllPoints() []NavigationDataPoint {
	var points []NavigationDataPoint
	// Store current position so we can restore it
	currentPos := r.pos
	// Reset position to start
	r.pos = 0

	for {
		d, ok := r.Next()
		if !ok {
			break
		}
		points = append(points, d)
	}

	// Restore original position
	r.pos = currentPos
	return points
}

// GetBounds returns the minimum and maximum latitude and longitude values from a slice of navigation points
func GetBounds(points []NavigationDataPoint) (minLat, maxLat, minLng, maxLng float64) {
	if len(points) == 0 {
		return 0, 0, 0, 0
	}

	minLat = points[0].Latitude
	maxLat = points[0].Latitude
	minLng = points[0].Longitude
	maxLng = points[0].Longitude

	for _, p := range points[1:] {
		if p.Latitude < minLat {
			minLat = p.Latitude
		}
		if p.Latitude > maxLat {
			maxLat = p.Latitude
		}
		if p.Longitude < minLng {
			minLng = p.Longitude
		}
		if p.Longitude > maxLng {
			maxLng = p.Longitude
		}
	}

	return minLat, maxLat, minLng, maxLng
}

func MedianWindDirection(points []NavigationDataPoint) float64 {
	windDirections := make([]float64, len(points))
	for i, p := range points {
		windDirections[i] = p.TrueWindDirection
	}

	sort.Float64s(windDirections)
	n := len(windDirections)
	medianWind := windDirections[n/2]
	if len(windDirections)%2 == 0 {
		medianWind = (windDirections[n/2] + windDirections[n/2+1]) / 2
	}
	return medianWind
}
