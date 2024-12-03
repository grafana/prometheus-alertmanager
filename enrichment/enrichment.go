// Copyright 2024 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package enrichment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	commoncfg "github.com/prometheus/common/config"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)

type Enrichments struct {
	enrichments []*Enrichment
}

func NewEnrichments(enrs []config.Enrichment) (*Enrichments, error) {
	enrichments := make([]*Enrichment, 0, len(enrs))

	for _, enr := range enrs {
		enrichment, err := NewEnrichment(enr)
		if err != nil {
			return nil, err
		}

		enrichments = append(enrichments, enrichment)
	}

	return &Enrichments{
		enrichments: enrichments,
	}, nil
}

func (e *Enrichments) Apply(ctx context.Context, l log.Logger, alerts ...*types.Alert) {
	var (
		success = 0
		failed  = 0
	)

	// TODO: These could/should be done async. Need to decide if to allow dependent enrichments.
	for i, enr := range e.enrichments {
		if err := enr.Apply(ctx, l, alerts...); err != nil {
			// Attempt to apply all enrichments, one doesn't need to affect the others.
			level.Error(l).Log("msg", "Enrichment failed", "i", i, "err", err)
			failed++
		} else {
			success++
		}
	}

	level.Debug(l).Log("msg", "Enrichments applied", "success", success, "failed", failed)
}

type Enrichment struct {
	conf   config.Enrichment
	client *http.Client
}

func NewEnrichment(conf config.Enrichment) (*Enrichment, error) {
	client, err := commoncfg.NewClientFromConfig(*conf.HTTPConfig, "enrichment")
	if err != nil {
		return nil, err
	}

	return &Enrichment{
		conf:   conf,
		client: client,
	}, nil
}

func (e *Enrichment) Apply(ctx context.Context, l log.Logger, alerts ...*types.Alert) error {
	// TODO: Template isn't needed by this function but we need to pass something.
	data := notify.GetTemplateData(ctx, &template.Template{}, alerts, l)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}

	url := e.conf.URL.String()

	if e.conf.Timeout > 0 {
		postCtx, cancel := context.WithTimeoutCause(ctx, e.conf.Timeout, fmt.Errorf("configured enrichment timeout reached (%s)", e.conf.Timeout))
		defer cancel()
		ctx = postCtx
	}

	resp, err := notify.PostJSON(ctx, e.client, url, &buf)
	if err != nil {
		if ctx.Err() != nil {
			err = fmt.Errorf("%w: %w", err, context.Cause(ctx))
		}
		return notify.RedactURL(err)
	}
	defer resp.Body.Close()

	var result template.Data
	err = json.NewDecoder(resp.Body).Decode(&result)

	// TODO: Do something with the result.
	// TODO: Don't log the URL unredacted.
	level.Info(l).Log("msg", "Enrichment result",
		"url", url,
		"groupLabels", result.GroupLabels,
		"commonLabels", result.CommonLabels,
		"commonAnnotations", result.CommonLabels)

	return nil
}
