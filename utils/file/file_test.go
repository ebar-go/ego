package file

import "testing"

func TestExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "exist",
			args: args{
				path:"E:\\Workspace\\go\\src\\github.com\\ebar-go\\ego\\README.md",
			},
			want: true,
		},
		{
			name: "notExist",
			args: args{
				path:"E:\\Workspace\\go\\src\\github.com\\ebar-go\\ego\\README01.md",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exist(tt.args.path); got != tt.want {
				t.Errorf("Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: "C:\\Users\\cx\\AppData\\Local\\Temp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrentDir(); got != tt.want {
				t.Errorf("GetCurrentDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrentPath(); got != tt.want {
				t.Errorf("GetCurrentPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetExecuteDir(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: "E:\\Workspace\\go\\src\\github.com\\ebar-go\\ego\\utils\\file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExecuteDir(); got != tt.want {
				t.Errorf("GetExecuteDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMkdir(t *testing.T) {
	type args struct {
		dir     string
		parents bool
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
			if err := Mkdir(tt.args.dir, tt.args.parents); (err != nil) != tt.wantErr {
				t.Errorf("Mkdir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}