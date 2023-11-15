package collect

import (
	"context"
	"testing"

	troubleshootv1beta2 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"
	corev1 "k8s.io/api/core/v1"
)

func Test_veleroCommandExec(t *testing.T) {
	type args struct {
		ctx             context.Context
		progressChan    chan<- interface{}
		c               *CollectVelero
		veleroCollector *troubleshootv1beta2.Velero
		pod             *corev1.Pod
		command         VeleroCommand
		output          CollectorResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := veleroCommandExec(tt.args.ctx, tt.args.progressChan, tt.args.c, tt.args.veleroCollector, tt.args.pod, tt.args.command, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("veleroCommandExec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
