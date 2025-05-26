package ecs

import (
	st "github.com/vistormu/xpeto/internal/structures"
)

type SystemFilter func(context *Context) bool

type SystemManager struct {
	loaded  *st.HashSet[System]
	filters map[System]SystemFilter
	systems []System
}

func NewSystemManager() *SystemManager {
	return &SystemManager{
		loaded:  st.NewHashSet[System](),
		filters: make(map[System]SystemFilter),
		systems: make([]System, 0),
	}
}

func (s *SystemManager) Register(system System, filter SystemFilter) {
	s.filters[system] = filter
	s.systems = append(s.systems, system)
}

func (s *SystemManager) Load(context *Context, system System) {
	if s.loaded.Contains(system) {
		return
	}

	filter, ok := s.filters[system]
	if !ok {
		return
	}

	if filter(context) {
		system.OnLoad(context)
		s.loaded.Add(system)
	}
}

func (s *SystemManager) Filter(context *Context, system System) bool {
	if filter, ok := s.filters[system]; ok {
		return filter(context)
	}
	return false
}

func (s *SystemManager) Systems() []System {
	return s.systems
}
