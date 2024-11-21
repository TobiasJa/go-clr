//go:build windows
// +build windows

package clr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeArrayCreate(t *testing.T) {
	type args struct {
		vt        VT
		cDims     uint32
		rgsabound *SafeArrayBound
	}
	type want struct {
		cDims      uint16
		cbElements uint32
		cLocks     uint32
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "Create SafeArray",
			args: args{
				vt:    VT_UI1,
				cDims: uint32(1),
				rgsabound: &SafeArrayBound{
					cElements: uint32(10000),
					lLbound:   int32(0),
				},
			},
			want: want{
				cDims:      uint16(1),
				cbElements: uint32(1),
				cLocks:     uint32(0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeArrayCreate(tt.args.vt, tt.args.cDims, tt.args.rgsabound)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeArrayCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want.cDims, got.cDims) {
				t.Errorf("SafeArrayCreate() cDims not equal: got.cDims = %v, tt.want.cDims %v", got.cDims, tt.want.cDims)
			}
			if !assert.Equal(t, tt.want.cbElements, got.cbElements) {
				t.Errorf("SafeArrayCreate() cbElements not equal: got.cbElements = %v, tt.want.cbElements %v", got.cbElements, tt.want.cbElements)
			}
			if !assert.Equal(t, tt.want.cLocks, got.cLocks) {
				t.Errorf("SafeArrayCreate() cLocks not equal: got.cLocks = %v, tt.want.cLocks %v", got.cLocks, tt.want.cLocks)
			}
			SafeArrayDestroy(got)
		})
	}
}
