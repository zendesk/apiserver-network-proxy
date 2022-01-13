package server

import (
	"context"

	"k8s.io/klog/v2"
	pkgagent "sigs.k8s.io/apiserver-network-proxy/pkg/agent"
)

type AgentIdBackendManager struct {
	*DefaultBackendStorage
}

var _ BackendManager = &AgentIdBackendManager{}

// NewAgentIdBackendManager returns a AgentIdBackendManager.
func NewAgentIdBackendManager() *AgentIdBackendManager {
	return &AgentIdBackendManager{
		DefaultBackendStorage: NewDefaultBackendStorage(
			[]pkgagent.IdentifierType{pkgagent.UID})}
}

func (dibm *AgentIdBackendManager) Backend(ctx context.Context) (Backend, error) {
	dibm.mu.RLock()
	defer dibm.mu.RUnlock()
	if len(dibm.backends) == 0 {
		return nil, &ErrNotFound{}
	}

	agentID := ctx.Value(agentId).(string)
	if agentID != "" {
		bes, exist := dibm.backends[agentID]
		if exist && len(bes) > 0 {
			klog.V(5).InfoS("Get the backend through the AgentIdBackendManager", "agentId", agentId)
			return dibm.backends[agentID][0], nil
		}
	}
	klog.V(5).Infof("no agent found for identifier: %v", agentID)
	return nil, &ErrNotFound{}
}
