package client

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

// templateCache shares GET /blueprints/<id> between the env0_template_project_assignment
// reads of one provider run. N assignments of one template otherwise send N identical GETs,
// all on the same rate-limit key.
type templateCache struct {
	mu        sync.Mutex
	templates map[string]Template
	// generation is bumped on every invalidation, so a fetch that started before it can't
	// store a template that is already stale.
	generation uint64
	// group collapses concurrent misses for one id (terraform refreshes in parallel) into one GET.
	group singleflight.Group
}

func newTemplateCache() *templateCache {
	return &templateCache{templates: map[string]Template{}}
}

func (c *templateCache) get(id string, fetch func(string) (Template, error)) (Template, error) {
	c.mu.Lock()

	if template, ok := c.templates[id]; ok {
		c.mu.Unlock()

		return template, nil
	}

	generation := c.generation
	c.mu.Unlock()

	result, err, _ := c.group.Do(id, func() (any, error) {
		template, err := fetch(id)
		if err != nil {
			return Template{}, err
		}

		c.mu.Lock()
		if c.generation == generation {
			c.templates[id] = template
		}
		c.mu.Unlock()

		return template, nil
	})

	return result.(Template), err
}

func (c *templateCache) invalidate(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.templates, id)
	c.generation++
	// Callers that arrive after the invalidation start a new GET instead of joining one in flight.
	c.group.Forget(id)
}
