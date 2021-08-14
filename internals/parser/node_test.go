package parser

import (
	"reflect"
	"testing"
)

func TestProjectManager_GetPackageNameFromString(t *testing.T) {
	type fields struct {
		ModFile string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		struct {
			name    string
			fields  fields
			args    args
			want    string
			wantErr bool
		}{
			name: "regexp get go module",
			fields: fields{
				ModFile: "",
			},
			args: args{
				s: "module github.com/zrb-inc/spring",
			},
			want:    "github.com/zrb-inc/spring",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProjectManager{
				ModFile: tt.fields.ModFile,
			}
			got, err := p.GetPackageNameFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectManager.GetPackageNameFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProjectManager.GetPackageNameFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectManager_GetPackageName(t *testing.T) {
	type fields struct {
		ModFile string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		struct {
			name    string
			fields  fields
			want    string
			wantErr bool
		}{
			name: "",
			fields: fields{
				ModFile: "./testdata/go.mod",
			},
			want:    "github.com/zrb-inc/spring",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProjectManager{
				ModFile: tt.fields.ModFile,
			}
			got, err := p.GetPackageName()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectManager.GetPackageName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProjectManager.GetPackageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_TravelRoot(t *testing.T) {
	type fields struct {
		Fs Fs
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Node
		wantErr bool
	}{
		struct {
			name    string
			fields  fields
			args    args
			want    *Node
			wantErr bool
		}{
			name: "test travel dir",
			fields: fields{
				Fs: Fs{},
			},
			args: args{
				path: "./testdata",
			},
			want:    &Node{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				Fs: tt.fields.Fs,
			}
			got, err := b.TravelRoot(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.TravelRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.TravelRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
