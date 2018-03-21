package qbittorrent_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/bearx3f/btwebapi/qbittorrent"
)

func TestNewClient(t *testing.T) {
	type args struct {
		baseURL  string
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantQb  *qbittorrent.Client
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQb, err := qbittorrent.NewClient(tt.args.baseURL, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotQb, tt.wantQb) {
				t.Errorf("NewClient() = %v, want %v", gotQb, tt.wantQb)
			}
		})
	}
}

func TestClient_Logout(t *testing.T) {
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.Logout(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_APIVersion(t *testing.T) {
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.qb.APIVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.APIVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.APIVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_APIMinimumVersion(t *testing.T) {
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.qb.APIMinimumVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.APIMinimumVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.APIMinimumVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListTorrent(t *testing.T) {
	tests := []struct {
		name         string
		qb           *qbittorrent.Client
		wantTorrents *[]qbittorrent.Torrent
		wantErr      bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTorrents, err := tt.qb.ListTorrent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListTorrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTorrents, tt.wantTorrents) {
				t.Errorf("Client.ListTorrent() = %v, want %v", gotTorrents, tt.wantTorrents)
			}
		})
	}
}

func TestClient_Start(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.Start(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_StartAll(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.StartAll(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.StartAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Pause(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.Pause(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.Pause() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_PauseAll(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.PauseAll(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.PauseAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Remove(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.Remove(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_RemoveData(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.RemoveData(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Recheck(t *testing.T) {
	type args struct {
		hash []string
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.Recheck(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Client.Recheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UploadFile(t *testing.T) {
	type args struct {
		file *os.File
	}
	tests := []struct {
		name    string
		qb      *qbittorrent.Client
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qb.UploadFile(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
