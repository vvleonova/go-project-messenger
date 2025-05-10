package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_HealthCheck(t *testing.T) {
	tests := []struct {
		name        string
		mockStorage *storageMock
		wantErr     bool
	}{
		{
			name: "ping successfull",
			mockStorage: &storageMock{
				PingFunc: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "ping fails",
			mockStorage: &storageMock{
				PingFunc: func() error {
					return errors.New("database is unreachable")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание сервиса с моком
			svc := &Service{
				storage: tt.mockStorage,
			}

			// проверка функции
			err := svc.HealthCheck()
			if tt.wantErr {
				assert.Error(t, err, "expected an error, got nil")
			} else {
				assert.NoError(t, err, "expected no error, got %v", err)
			}
		})
	}
}
