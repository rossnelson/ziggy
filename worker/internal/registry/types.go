package registry

// Definition describes a workflow for self-registration.
type Definition struct {
	Name      string
	Workflow  interface{}
	IDPattern func(owner string) string
	NewInput  func(owner, ziggyID, tz string) any
	AutoStart bool
}

// ActivityDef describes an activity for self-registration.
type ActivityDef struct {
	Name     string
	Activity interface{}
}

// Config holds the configuration for the registry.
type Config struct {
	HostPort      string
	Namespace     string
	TaskQueue     string
	Owner         string
	Timezone      string
	StartWorkflow bool
}
