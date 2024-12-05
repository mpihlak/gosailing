package datasource

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReplayNavigationDataProvider(t *testing.T) {
	require := require.New(t)

	testData := `time,dpt,mtw,awa,aws,cog,hdg,rot,sog,stw,pitch,yaw,roll,lng,lat,tws,twa,twd,vmg,dist,cum_dist
2024-09-11T17:27:52+03:00,22.5,288.65,31.999692858056477,17.5932,87.18125810710607,78,1.5999846429028237,5.851439999999999,5.5,-3.0997016716577535,-101.79741146089336,-20.597832734953094,24.798672,59.488683,13.00,45.79,132.25,4.079967736770225,0.0015914414543093099,55
`
	buf := bytes.NewBuffer([]byte(testData))
	ds, err := NewReplayNavigationDataProvider(buf)
	require.NoError(err)
	require.NotNil(ds)

	d, ok := ds.Next()
	require.True(ok)

	require.Equal("2024-09-11T17:27:52+03:00", d.Timestamp.Format(time.RFC3339))
	require.Equal(float64(78), d.Heading)
	require.Equal(float64(5.5), d.SpeedThroughWater)
	require.Equal(float64(45.79), d.TrueWindAngle)
	require.Equal(float64(13.0), d.TrueWindSpeed)
	require.Equal(float64(132.25), d.TrueWindDirection)
	require.Equal(float64(59.488683), d.Latitude)
	require.Equal(float64(24.798672), d.Longitude)
	require.Equal(float64(55), d.CumulativeDistance)

	d, ok = ds.Next()
	require.False(ok)
}
