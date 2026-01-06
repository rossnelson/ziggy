package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Registration struct {
	Name     string
	Type     string // "workflow" or "activity"
	Function interface{}
}


type WorkflowStatus struct {
	Status    string
	CloseTime time.Time
}

type Registry struct {
	mu            sync.RWMutex
	client        client.Client
	registrations []Registration
	config        Config
}

var (
	instance     *Registry
	once         sync.Once
	workflowDefs []Definition
	activityDefs []ActivityDef
	defMu        sync.RWMutex
)

func Get() *Registry {
	once.Do(func() {
		instance = &Registry{
			registrations: make([]Registration, 0),
		}
	})
	return instance
}

func NewRegistry() *Registry {
	return &Registry{
		registrations: make([]Registration, 0),
	}
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

// RegisterWorkflow registers a workflow definition for self-registration via init().
func RegisterWorkflow(def Definition) {
	defMu.Lock()
	defer defMu.Unlock()
	workflowDefs = append(workflowDefs, def)
}

// RegisterActivity registers an activity definition for self-registration via init().
func RegisterActivity(def ActivityDef) {
	defMu.Lock()
	defer defMu.Unlock()
	activityDefs = append(activityDefs, def)
}

// GetWorkflowDefs returns all registered workflow definitions.
func GetWorkflowDefs() []Definition {
	defMu.RLock()
	defer defMu.RUnlock()
	return workflowDefs
}

// GetActivityDefs returns all registered activity definitions.
func GetActivityDefs() []ActivityDef {
	defMu.RLock()
	defer defMu.RUnlock()
	return activityDefs
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
	cfg := r.config
	registrations := r.registrations
	r.mu.RUnlock()

	w := worker.New(c, cfg.TaskQueue, worker.Options{})

	// Register legacy registrations (from AddWorkflow/AddActivity)
	for _, reg := range registrations {
		switch reg.Type {
		case "workflow":
			w.RegisterWorkflowWithOptions(reg.Function, workflow.RegisterOptions{
				Name: reg.Name,
			})
		case "activity":
			w.RegisterActivityWithOptions(reg.Function, activity.RegisterOptions{
				Name: reg.Name,
			})
		}
	}

	// Register new self-registered workflows
	for _, def := range GetWorkflowDefs() {
		w.RegisterWorkflowWithOptions(def.Workflow, workflow.RegisterOptions{
			Name: def.Name,
		})
	}

	// Register new self-registered activities
	for _, def := range GetActivityDefs() {
		w.RegisterActivityWithOptions(def.Activity, activity.RegisterOptions{
			Name: def.Name,
		})
	}

	// Start auto-start workflows
	if cfg.StartWorkflow {
		if err := r.ensureWorkflowsRunning(ctx, cfg); err != nil {
			return fmt.Errorf("ensure workflows: %w", err)
		}
	}

	return w.Run(worker.InterruptCh())
}

func (r *Registry) ensureWorkflowsRunning(ctx context.Context, cfg Config) error {
	defs := GetWorkflowDefs()
	if len(defs) == 0 {
		return nil
	}

	// Get the primary workflow ID (Ziggy) for dependent workflows
	var ziggyID string
	for _, def := range defs {
		if def.Name == "ZiggyWorkflow" && def.IDPattern != nil {
			ziggyID = def.IDPattern(cfg.Owner)
			break
		}
	}

	for _, def := range defs {
		if !def.AutoStart {
			continue
		}
		if def.IDPattern == nil {
			continue
		}

		id := def.IDPattern(cfg.Owner)
		var input interface{}
		if def.NewInput != nil {
			input = def.NewInput(cfg.Owner, ziggyID, cfg.Timezone)
		}

		if err := r.ensureWorkflow(ctx, id, def.Name, input); err != nil {
			return fmt.Errorf("workflow %s: %w", def.Name, err)
		}
	}

	return nil
}

func (r *Registry) ensureWorkflow(ctx context.Context, workflowID, workflowName string, input interface{}) error {
	status, err := r.DescribeWorkflow(ctx, workflowID)
	if err == nil && status.Status == "WORKFLOW_EXECUTION_STATUS_RUNNING" {
		fmt.Printf("Workflow %s already running\n", workflowID)
		return nil
	}

	fmt.Printf("Starting workflow %s (%s)\n", workflowID, workflowName)
	_, err = r.ExecuteWorkflow(ctx, workflowID, workflowName, input)
	return err
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

func (r *Registry) DescribeWorkflow(ctx context.Context, workflowID string) (*WorkflowStatus, error) {
	r.mu.RLock()
	if r.client == nil {
		r.mu.RUnlock()
		return nil, fmt.Errorf("registry not initialized")
	}
	c := r.client
	r.mu.RUnlock()

	desc, err := c.DescribeWorkflowExecution(ctx, workflowID, "")
	if err != nil {
		return nil, err
	}

	status := &WorkflowStatus{
		Status: desc.WorkflowExecutionInfo.Status.String(),
	}

	if desc.WorkflowExecutionInfo.CloseTime != nil {
		status.CloseTime = desc.WorkflowExecutionInfo.CloseTime.AsTime()
	}

	return status, nil
}

func (r *Registry) TerminateWorkflow(ctx context.Context, workflowID, reason string) error {
	r.mu.RLock()
	if r.client == nil {
		r.mu.RUnlock()
		return fmt.Errorf("registry not initialized")
	}
	c := r.client
	r.mu.RUnlock()

	return c.TerminateWorkflow(ctx, workflowID, "", reason)
}

func (r *Registry) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		r.client.Close()
		r.client = nil
	}
}
