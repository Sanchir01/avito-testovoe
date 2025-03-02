package mocks

import (
	"context"
	"github.com/Sanchir01/avito-testovoe/internal/feature/product"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestHandlerProducts_GetAllProducts(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test get all products",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productProvider := NewHandlerProducts(t)
			productProvider.On("GetAllProducts", mock.Anything).Return([]*product.DataBaseProduct{}, nil)
			_m := &HandlerProducts{
				Mock: productProvider.Mock,
			}
			_, err := _m.GetAllProducts(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestHandlerProducts_GetProductByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name string

		args args

		wantErr bool
	}{
		{
			name: "test get product by id",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productProvider := NewHandlerProducts(t)
			productProvider.On("GetProductByID", mock.Anything, mock.Anything).Return(&product.DataBaseProduct{}, nil).Once()
			_m := &HandlerProducts{
				Mock: productProvider.Mock,
			}
			_, err := _m.GetProductByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
