package stage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudzero/cloudzero-agent-validator/pkg/config"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/diagnostic/stage"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/status"
)

func makeReport() status.Accessor {
	return status.NewAccessor(&status.ClusterStatus{})
}

func TestChecker_CheckOK(t *testing.T) {
	tcases := []struct {
		name    string
		stageID status.StatusType
	}{
		{
			name:    "init started",
			stageID: status.StatusType_STATUS_TYPE_INIT_STARTED,
		},
		{
			name:    "init stopped",
			stageID: status.StatusType_STATUS_TYPE_INIT_STARTED,
		},
		{
			name:    "pod stopped",
			stageID: status.StatusType_STATUS_TYPE_POD_STARTED,
		},
		{
			name:    "pod stopped",
			stageID: status.StatusType_STATUS_TYPE_POD_STOPPING,
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Settings{}
			provider := stage.NewProvider(context.Background(), cfg, tc.stageID)
			accessor := makeReport()
			assert.NoError(t, provider.Check(context.Background(), nil, accessor, cfg))
			accessor.ReadFromReport(func(s *status.ClusterStatus) {
				assert.Equal(t, tc.stageID, s.State)
			})
		})
	}
}
