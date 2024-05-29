package qhn

import (
	"context"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/lunagic/quiet-hacker-news/internal/hackernews"
	"github.com/lunagic/quiet-hacker-news/internal/resources"
)

func New(hn *hackernews.Client) (*Service, error) {
	indexTemplate, err := template.New("index").Parse(resources.IndexTemplate)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:    hn,
		template:  indexTemplate,
		startTime: time.Now(),
		mutex:     &sync.Mutex{},
		cache:     []hackernews.Item{},
	}, nil
}

type Service struct {
	client    *hackernews.Client
	template  *template.Template
	startTime time.Time
	mutex     *sync.Mutex
	cache     []hackernews.Item
}

func (s *Service) Background(ctx context.Context) error {
	items, err := s.client.TopStories(ctx, 30)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache = items

	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFileFS(w, r, resources.Public, "public"+r.URL.Path)

		return
	}

	s.rootHandler(w, r)
}

type TemplatePayload struct {
	Items  []hackernews.Item
	Uptime string
}

func (s *Service) rootHandler(w http.ResponseWriter, _ *http.Request) {
	s.template.Execute(w, TemplatePayload{
		Uptime: time.Since(s.startTime).Round(time.Second).String(),
		Items:  s.cache,
	})
}
