package temporal

import (
	"context"
	"fmt"
	"sync"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Registration struct {
	Name     string
	Type     string // "workflow" or "activity"
	Function interface{}
}

type Config struct {
	HostPort  string
	Namespace string
	TaskQueue string
}

type Registry struct {
	mu            sync.RWMutex
	client        client.Client
	registrations []Registration
	config        Config
}

var (
	instance *Registry
	once     sync.Once
)

func Get() *Registry {
	once.Do(func() {
		instance = &Registry{
			registrations: make([]Registration, 0),
		}
	})
	return instance
}

func AddWorkflow(name string, workflow interface{}) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = append(r.registrations, Registration{
		Name:     name,
		Type:     "workflow",
		Function: workflow,
	})
}

func AddActivity(name string, activity interface{}) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = append(r.registrations, Registration{
		Name:     name,
		Type:     "activity",
		Function: activity,
	})
}

func (r *Registry) Initialize(cfg Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		return nil
	}

	r.config = cfg

	c, err := client.Dial(client.Options{
		HostPort:  cfg.HostPort,
		Namespace: cfg.Namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to create temporal client: %w", err)
	}

	r.client = c
	return nil
}

func (r *Registry) StartWorker(ctx context.Context) error {
	r.mu.RLock()
	if r.client == nil {
		r.mu.RUnlock()
		return fmt.Errorf("registry not initialized")
	}
	c := r.client
	taskQueue := r.config.TaskQueue
	registrations := r.registrations
	r.mu.RUnlock()

	w := worker.New(c, taskQueue, worker.Options{})

	for _, reg := range registrations {
		switch reg.Type {
		case "workflow":
			w.RegisterWorkflow(reg.Function)
		case "activity":
			w.RegisterActivity(reg.Function)
		}
	}

	return w.Run(worker.InterruptCh())
}

func (r *Registry) ExecuteWorkflow(ctx context.Context, workflowID string, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	r.mu.RLock()
	if r.client == nil {
		r.mu.RUnlock()
		return nil, fmt.Errorf("registry not initialized")
	}
	c := r.client
	taskQueue := r.config.TaskQueue
	r.mu.RUnlock()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueue,
	}

	return c.ExecuteWorkflow(ctx, options, workflow, args...)
}

func (r *Registry) SignalWorkflow(ctx context.Context, workflowID, signalName string, arg interface{}) error {
	r.mu.RLock()
	if r.client == nil {
		r.mu.RUnlock()
		return fmt.Errorf("registry not initialized")
	}
	c := r.client
	r.mu.RUnlock()

	return c.SignalWorkflow(ctx, workflowID, "", signalName, arg)
}

func (r *Registry) QueryWorkflow(ctx context.Context, workflowID, queryType string, args ...interface{}) (interface{}, error) {
	r.mu.RLock()
	if r.client == nil {
		r.mu.RUnlock()
		return nil, fmt.Errorf("registry not initialized")
	}
	c := r.client
	r.mu.RUnlock()

	response, err := c.QueryWorkflow(ctx, workflowID, "", queryType, args...)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := response.Get(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Registry) GetClient() client.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}

func (r *Registry) Shutdown() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		r.client.Close()
		r.client = nil
	}
}
