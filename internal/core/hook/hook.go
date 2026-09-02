package hook

import (
	"context"
	"sort"
	"sync"
)

// ActionHandler is a callback for side-effect events (WordPress actions).
type ActionHandler func(ctx context.Context, payload any) error

// FilterHandler is a callback for modifying/transforming data (WordPress filters).
type FilterHandler func(ctx context.Context, data any) (any, error)

type actionEntry struct {
	priority int
	handler  ActionHandler
}

type filterEntry struct {
	priority int
	handler  FilterHandler
}

// HookManager manages WordPress-style actions and filters.
type HookManager interface {
	AddAction(name string, priority int, handler ActionHandler)
	DoAction(ctx context.Context, name string, payload any) error

	AddFilter(name string, priority int, handler FilterHandler)
	ApplyFilter(ctx context.Context, name string, data any) (any, error)
}

type hookManager struct {
	mu      sync.RWMutex
	actions map[string][]actionEntry
	filters map[string][]filterEntry
}

// NewHookManager creates a new HookManager.
func NewHookManager() HookManager {
	return &hookManager{
		actions: make(map[string][]actionEntry),
		filters: make(map[string][]filterEntry),
	}
}

// AddAction adds a listener to an action hook with a specific priority (lower runs earlier).
func (h *hookManager) AddAction(name string, priority int, handler ActionHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entries := append(h.actions[name], actionEntry{priority: priority, handler: handler})
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].priority < entries[j].priority
	})
	h.actions[name] = entries
}

// DoAction executes all action handlers registered for the given hook name.
func (h *hookManager) DoAction(ctx context.Context, name string, payload any) error {
	h.mu.RLock()
	entries, exists := h.actions[name]
	if !exists || len(entries) == 0 {
		h.mu.RUnlock()
		return nil
	}
	handlers := make([]ActionHandler, len(entries))
	for i, e := range entries {
		handlers[i] = e.handler
	}
	h.mu.RUnlock()

	for _, fn := range handlers {
		if err := fn(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}

// AddFilter adds a filter hook callback with a specific priority (lower runs earlier).
func (h *hookManager) AddFilter(name string, priority int, handler FilterHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entries := append(h.filters[name], filterEntry{priority: priority, handler: handler})
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].priority < entries[j].priority
	})
	h.filters[name] = entries
}

// ApplyFilter runs data through all filter handlers registered for the given hook name.
func (h *hookManager) ApplyFilter(ctx context.Context, name string, data any) (any, error) {
	h.mu.RLock()
	entries, exists := h.filters[name]
	if !exists || len(entries) == 0 {
		h.mu.RUnlock()
		return data, nil
	}
	handlers := make([]FilterHandler, len(entries))
	for i, e := range entries {
		handlers[i] = e.handler
	}
	h.mu.RUnlock()

	current := data
	var err error
	for _, fn := range handlers {
		current, err = fn(ctx, current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
