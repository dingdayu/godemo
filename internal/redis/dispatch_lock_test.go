package redis

import (
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
)

func TestMain(m *testing.M) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client = redis.NewClient(&redis.Options{
		Addr: s.Addr()})
	m.Run()
}

func TestNewDispatchLock(t *testing.T) {
	type args struct {
		client *redis.Client
	}
	tests := []struct {
		name string
		args args
		want *DispatchLock
	}{
		{name: "NewDispatchLock", args: args{client: GetClient()}, want: &DispatchLock{client: GetClient()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDispatchLock(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDispatchLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDispatchLock_Lock(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		id         int
		expiration time.Duration
	}

	arg := args{id: 1, expiration: time.Minute}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{name: "lock", fields: fields{client: GetClient()}, args: arg, want: true, wantErr: false},
		{name: "unlock", fields: fields{client: GetClient()}, args: arg, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DispatchLock{
				client: tt.fields.client,
			}
			got, err := s.Lock(tt.args.id, tt.args.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lock() got = %v, want %v", got, tt.want)
			}
		})
	}
}
