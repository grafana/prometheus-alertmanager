package enrichment

import (
	"net/http"

	"github.com/go-kit/log"
	//	"github.com/go-kit/log/level"

	commoncfg "github.com/prometheus/common/config"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/alertmanager/template"
)

type Enricher struct {
	// Map of receiver name to enrichments
	enrichments map[string][]*Enrichment
}

func NewEnricher(enrs []config.Enrichment, logger log.Logger) (*Enricher, error) {
	enrichments := make(map[string][]*Enrichment)

	for _, enr := range enrs {
		enrichment, err := NewEnrichment(enr, logger)
		if err != nil {
			return nil, err
		}

		for _, rcv := range enr.Receivers {
			enrichments[rcv] = append(enrichments[rcv], enrichment)
		}
	}

	return &Enricher{
		enrichments: enrichments,
	}, nil
}

func (e *Enricher) GetEnrichments(rcv string) []*Enrichment {
	return e.enrichments[rcv]
}

type Enrichment struct {
	logger   log.Logger
	matchers labels.Matchers
	client   *http.Client
}

func NewEnrichment(conf *config.Enrichment, logger log.Logger) (*Enrichment, error) {
	client, err := commoncfg.NewClientFromConfig(*conf.HTTPConfig, "enrichment", httpOpts...)
	if err != nil {
		return nil, err
	}

	return &Enrichment{
		logger:   log.With(logger, "enrichment", conf.Name),
		matchers: append(labels.Matchers{}, conf.Matchers...),
		client:   client,
	}
}

func (e *Enrichment) Apply(ctx context.Context, alerts ...*types.Alert) error {
	// TODO: Template isn't needed by this function but we need to pass something.
	data := notify.GetTemplateData(ctx, &template.Template{}, alerts, e.logger)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return false, err
	}

	var url string
	if n.conf.URL != nil {
		url = n.conf.URL.String()
	} else {
		content, err := os.ReadFile(n.conf.URLFile)
		if err != nil {
			return false, fmt.Errorf("read url_file: %w", err)
		}
		url = strings.TrimSpace(string(content))
	}

	if n.conf.Timeout > 0 {
		postCtx, cancel := context.WithTimeoutCause(ctx, n.conf.Timeout, fmt.Errorf("configured webhook timeout reached (%s)", n.conf.Timeout))
		defer cancel()
		ctx = postCtx
	}

	resp, err := notify.PostJSON(ctx, n.client, url, &buf)
	if err != nil {
		if ctx.Err() != nil {
			err = fmt.Errorf("%w: %w", err, context.Cause(ctx))
		}
		return true, notify.RedactURL(err)
	}
	defer notify.Drain(resp)

	// TODO: Parse and merge result.

	return nil
}
