package main

import (
	"errors"
	//"fmt"
	"net/http"
	"text/template"
)

var SplineDataArray = `
			[
	            {
	                name: 'Tokyo',
	                data: [7.0, 6.9, 9.5, 14.5, 18.2, 21.5, 25.2, 26.5, 23.3, 18.3, 13.9, 9.6]
	            }, 
	            {
	                name: 'New York',
	                data: [-0.2, 0.8, 5.7, 11.3, 17.0, 22.0, 24.8, 24.1, 20.1, 14.1, 8.6, 2.5]
	            }, 
	            {
	                name: 'Berlin',
	                data: [-0.9, 0.6, 3.5, 8.4, 13.5, 17.0, 18.6, 17.9, 14.3, 9.0, 3.9, 1.0]
	            }, 
	            {
	                name: 'London',
	                data: [3.9, 4.2, 5.7, 8.5, 11.9, 15.2, 17.0, 16.6, 14.2, 10.3, 6.6, 4.8]
	            }
	        ]
`

var PieDataArray = `
                    [
                        ['Firefox',   45.0],
                        ['IE',       26.8],
                        ['Chrome',  12.8],
                        ['Safari',    8.5],
                        ['Opera',     6.2],
                        ['Others',   0.7]
                    ]
`

var ArgsSpline = map[string]string{
	"ChartType":    "spline",
	"Title":        "Monthly Average Temperature",
	"SubTitle":     "Source: WorldClimate.com",
	"YAxisText":    "Temperature (°C)",
	"XAxisNumbers": "['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12']",
	"ValueSuffix":  "°C",
	"DataArray":    SplineDataArray,
}

var ArgsColumn = map[string]string{

	"ChartType":    "column",
	"Title":        "Monthly Average Temperature",
	"SubTitle":     "Source: WorldClimate.com",
	"YAxisText":    "Temperature (°C)",
	"XAxisNumbers": "['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12']",
	"ValueSuffix":  "°C",
	"DataArray":    SplineDataArray,
}

var ArgsPie = map[string]string{

	"ChartType": "pie",
	"Title":     "Browser market shares at a specific website, 2014",
	"SubTitle":  "Source: website.com",
	"SerieName": "Browser shares",
	"DataArray": PieDataArray,
}

var (
	ChartHandlers = make(map[string]ChartIf)
	ChartFiles    []string
	Index         int
)

func handler(w http.ResponseWriter, r *http.Request) {
	var Args map[string]string
	var tmpl string

	Args = ArgsColumn
	Args = ArgsSpline
	Args = ArgsPie
	tmpl = TemplatePieHtml

	path := r.URL.Path[1:]
	if path == "spline" {
		tmplagrs, err := Parse(ChartFiles[Index])
		Index++
		Index = Index % len(ChartFiles)
		if err == nil {
			Args = tmplagrs.args
			tmpl = TemplateSplineHtml
		}
	}

	if t, err := template.New("foo").Parse(tmpl); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		if err = t.ExecuteTemplate(w, "T", Args); err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	var err error
	ChartFiles, err = LookupChartFiles(".")
	if err != nil {
		return err
	}
	if len(ChartFiles) == 0 {
		return errors.New("No chart data.")
	}

	ChartHandlers["spline"] = new(SplineChart)

	return http.ListenAndServe(addr, nil)
}
