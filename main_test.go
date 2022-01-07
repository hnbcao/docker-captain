package main

import "testing"

func Test_buildAuth(t *testing.T) {
	type args struct {
		userName string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "hub.d.cisdigital.cn",
			args: args{
				userName: "admin",
				password: "fOciCYZ4EolU1uHR",
			},
			want:    "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJmT2NpQ1laNEVvbFUxdUhSIn0=",
			wantErr: false,
		},
		{
			name: "harbor.test.cisdigital.cn",
			args: args{
				userName: "admin",
				password: "cisdigital-12345",
			},
			want:    "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJjaXNkaWdpdGFsLTEyMzQ1In0=",
			wantErr: false,
		},
		{
			name: "harbor.dev.cisdigital.cn",
			args: args{
				userName: "admin",
				password: "DGLa8rhyzCWnnDQf",
			},
			want:    "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJER0xhOHJoeXpDV25uRFFmIn0=",
			wantErr: false,
		},
		{
			name: "harbor.ng.cloud.nisco.cn",
			args: args{
				userName: "admin",
				password: "Harbor12345",
			},
			want:    "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJIYXJib3IxMjM0NSJ9",
			wantErr: false,
		},
		{
			name: "harbor.cisdigital.cn",
			args: args{
				userName: "admin",
				password: "3fmeGDOoqdI5D7Qz",
			},
			want:    "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIzZm1lR0RPb3FkSTVEN1F6In0=",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildAuth(tt.args.userName, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildAuth() got = %v, want %v", got, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}
