package generator

import (
	"testing"
)

func Test(t *testing.T) {
	type fields struct {
		orgName      string
		endPoint     string
		outDir       string
		peerConut    int
		beginPort    int
		chaincodes   []string
		batchTimeout string
		batchSize    *BatchSize
		cpuLimit     float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				orgName:      "org1",
				endPoint:     "flxdu.cn",
				outDir:       "./fabricGen",
				peerConut:    6,
				beginPort:    7050,
				chaincodes:   []string{"cc1", "cc2"},
				batchTimeout: "2s",
				batchSize:    NewBatchSize("10", "98 MB", "512 KB"),
				cpuLimit:     1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewConfigtx(tt.fields.orgName, tt.fields.endPoint, tt.fields.outDir, tt.fields.peerConut, tt.fields.beginPort, tt.fields.chaincodes, tt.fields.batchTimeout, tt.fields.batchSize, tt.fields.cpuLimit)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigtx() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err = c.Gen(); (err != nil) != tt.wantErr {
				t.Errorf("Gen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
