package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_HealthCheck(t *testing.T) {
	tests := []struct {
		name     string
		mockPing func() error
		wantErr  bool
	}{
		{
			name: "ping successfull",
			mockPing: func() error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "ping fails",
			mockPing: func() error {
				return errors.New("database is unreachable")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание мока с нужным поведением
			mockStorage := &storageMock{
				PingFunc: tt.mockPing,
			}

			// создание сервиса с моком
			svc := &Service{
				storage: mockStorage,
			}

			err := svc.HealthCheck()
			if tt.wantErr {
				assert.Error(t, err, "expected an error, got nil")
			} else {
				assert.NoError(t, err, "expected no error, got %v", err)
			}
		})
	}

}
