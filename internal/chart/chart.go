package chart

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/internal/database"
)

type ChartDataPoint struct {
	Date  string
	Value int
}

func GetChartData(db *sqlx.DB, days int) ([]ChartDataPoint, error) {
	data := buildDataPointDateRange(days)

	for i, d := range data {

		querydate, err := parseQueryDate(d.Date)
		if err != nil {
			return nil, err
		}

		count, err := database.GetClaimCountByDate(db, querydate)
		if err != nil {
			return nil, err
		}

		data[i].Value = count
	}

	// Reverse the order of the array
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data, nil
}

func buildDataPointDateRange(days int) []ChartDataPoint {
	var dataPoints []ChartDataPoint

	t := time.Now()
	t.Date()

	dataPoints = append(dataPoints, ChartDataPoint{Date: formatChartDate(t), Value: 0})

	for i := 1; i < days; i++ {
		v := t.Add(time.Duration(-24*i) * time.Hour)

		dataPoints = append(dataPoints, ChartDataPoint{Date: formatChartDate(v), Value: 0})
	}
	return dataPoints
}

func formatChartDate(date time.Time) string {
	return date.Format("02/01/2006")
}

func parseQueryDate(date string) (string, error) {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		return "", err
	}

	return t.Format("Mon, 02 Jan 2006"), nil
}
