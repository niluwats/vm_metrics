// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package insights

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/niluwats/vm_metrics/dbhandle"
	"github.com/niluwats/vm_metrics/internal/config"
	"github.com/niluwats/vm_metrics/internal/iam"
)

// ListMetricDefinitions returns the list of metrics available for the specified resource in the form "Localized Name (metric name)".
func ListMetricDefinitions(resourceURI string) ([]string, error) {
	a, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		return nil, err
	}

	metricsDefClient := insights.NewMetricDefinitionsClient(config.SubscriptionID())
	metricsDefClient.Authorizer = a
	metricsDefClient.AddToUserAgent(config.UserAgent())
	result, err := metricsDefClient.List(context.Background(), resourceURI, "")
	if err != nil {
		return nil, err
	}
	metrics := make([]string, len(*result.Value))
	for i := range *result.Value {
		metrics[i] = fmt.Sprintf("%s (%s)", *(*result.Value)[i].Name.LocalizedValue, *(*result.Value)[i].Name.Value)
	}
	return metrics, nil
}

// GetMetricsData returns the specified metric data points for the specified resource ID spanning the last five minutes.
func GetMetricsData(ctx context.Context, resourceID string, metrics []string) ([]string, error) {
	a, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		return nil, err
	}
	metricsClient := insights.NewMetricsClient(config.SubscriptionID())
	metricsClient.Authorizer = a
	metricsClient.AddToUserAgent(config.UserAgent())

	endTime := time.Now().UTC()
	startTime := endTime.Add(time.Duration(-2) * time.Minute)
	timespan := fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	resp, err := metricsClient.List(context.Background(), resourceID, timespan, nil, strings.Join(metrics, ","), "Average", nil, "", "", insights.Data, "")

	if err != nil {
		return nil, err
	}
	var metricData []string
	var MetricToSave dbhandle.Metrics

	for _, v := range *resp.Value {
		for _, t := range *v.Timeseries {
			for _, mv := range *t.Data {
				avg := 0.0
				if mv.Average != nil {
					avg = byteToMb(*mv.Average)
				}
				metricData = append(metricData, fmt.Sprintf("%s @ %s - avg: %f", *v.Name.LocalizedValue, *mv.TimeStamp, avg))
				if *v.Name.LocalizedValue == "Network In Total" {
					MetricToSave.NetworkIn = map[string]interface{}{
						"metric_name": "Network In Total",
						"time":        *mv.TimeStamp,
						"average":     avg,
					}
				} else if *v.Name.LocalizedValue == "Network Out Total" {
					MetricToSave.NetworkOut = map[string]interface{}{
						"metric_name": "Network Out Total",
						"time":        *mv.TimeStamp,
						"average":     avg,
					}
				} else if *v.Name.LocalizedValue == "Percentage CPU" {
					MetricToSave.PercentCpu = map[string]interface{}{
						"metric_name": "Percentage CPU",
						"time":        *mv.TimeStamp,
						"average":     avg,
					}
				} else {
					MetricToSave.AvailableMemBytes = map[string]interface{}{
						"metric_name": "Available Memory Bytes",
						"time":        *mv.TimeStamp,
						"average":     avg,
					}
				}
				break
			}
		}
	}
	dbhandle.SaveMetrics(MetricToSave)
	return metricData, nil
}
func byteToMb(i float64) float64 {
	return i / (1024 * 1024)
}
