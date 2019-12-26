package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// var myM *testing.M

func TestMain(m *testing.M) {
	clearUp()
	setup()
	code := m.Run()
	// myM = m
	// clearUp()
	os.Exit(code)
}

//  func Test_create(t *testing.T) {
//  	clearUp()

// }

func setup() {
	fmt.Println("=============Setup==================")
	go main()
	time.Sleep(time.Duration(5) * time.Second)
	go startExpireChecking()
}
func clearUp() {
	fmt.Println("=============Cleanup==================")
	cmd := exec.Command("/bin/sh", "-c", "yes | cp -f dnn.sqlite /nfs/;echo 'mac:$1$QIquJ0/h$aAzQvA.be4l3MGliMOlZF.' > /nfs/passwd;echo '/nfs/mac *(rw,no_subtree_check,sync,all_squash,anonuid=0,anongid=0)' > /nfs/exports")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func createAccount(machineid, account, passwd string) string {
	json := `{"machineId" : "%v", "gpuType" : "v100", "imgTag" : "simple:201706", "account" : "%v", "pwd" : "%v"}`
	jsonStr := []byte(fmt.Sprintf(json, machineid, account, passwd))
	req, _ := http.NewRequest("POST", "/storage", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	HandleStorage(w, req)

	cmd := exec.Command("/bin/sh", "-c", "grep", account, "/nfs/passwd")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("FTP account not found = %v", account)
		return "fail"
	}
	cmd = exec.Command("/bin/sh", "-c", "grep", account, "/nfs/exports")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("NFS account not found = %v", account)
		return "fail"
	}
	if existed, _ := exists("/nfs/" + account); !existed {
		fmt.Printf("Folder not found = %v", account)
		return "fail"
	} else {
		return "OK"
	}
}

func delete(machineid, account string) {
	r := mux.NewRouter()
	r.HandleFunc("/storage/{machineId}", HandleStorage)
	req, _ := http.NewRequest("DELETE", "/storage/"+machineid, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func checkDelete(account string) string {
	cmd := exec.Command("sh", "check.sh", account)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return strings.TrimSpace(out.String())
}

func startExpireChecking() {
	CheckExpire()
	time.Sleep(time.Duration(10) * time.Second)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Test_CreateAndDelete(t *testing.T) {
	type args struct {
		machineId string
		account   string
		password  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Create and Delete with the same account",
			args: args{machineId: "m1", account: "950223", password: "123456"},
			want: "OK",
		},
		{
			name: "Create and Delete with the same account",
			args: args{machineId: "m2", account: "950223", password: "123456"},
			want: "OK",
		},
		{
			name: "Create and Delete with the same account",
			args: args{machineId: "m3", account: "950223", password: "123456"},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createAccount(tt.args.machineId, tt.args.account, tt.args.password); got != tt.want {
				t.Errorf("createAccount() = %v, want %v", got, tt.want)
			}
			delete(tt.args.machineId, tt.args.account)
			time.Sleep(11 * time.Second)

		})
	}
	if got := checkDelete("950223"); got != "OK" {
		t.Errorf("checkDelete() = %v, want %v", got, "OK")
	}
}

func Test_CreateMachineWithTheSameAccount(t *testing.T) {
	type args struct {
		machineId string
		account   string
		password  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Create Machine with the same account",
			args: args{machineId: "m1", account: "950223", password: "123456"},
			want: "OK",
		},
		{
			name: "Create Machine with the same account",
			args: args{machineId: "m2", account: "950223", password: "123456"},
			want: "OK",
		},
		{
			name: "Create Machine with the same account",
			args: args{machineId: "m3", account: "950223", password: "123456"},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createAccount(tt.args.machineId, tt.args.account, tt.args.password); got != tt.want {
				t.Errorf("HandleStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_Delete(t *testing.T) {
	type args struct {
		machineId string
		account   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Delete Machine",
			args: args{machineId: "m1", account: "950223"},
			want: "fail",
		},
		{
			name: "Delete Machine",
			args: args{machineId: "m2", account: "950223"},
			want: "fail",
		},
		{
			name: "Delete Machine",
			args: args{machineId: "m3", account: "950223"},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := delete(tt.args.machineId, tt.args.account); got != tt.want {
				t.Errorf("HandleStorage() = %v, want %v", got, tt.want)
			}
			//time.Sleep(10 * time.Second)
		})
	}
}

func Test_CreateMachineWithDifferentAccount(t *testing.T) {
	type args struct {
		machineId string
		account   string
		password  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Create Machine with the same account",
			args: args{machineId: "m1", account: "950223", password: "123456"},
			want: "OK",
		},
		{
			name: "Create Machine with the same account",
			args: args{machineId: "m2", account: "950224", password: "123456"},
			want: "OK",
		},
		{
			name: "Create Machine with the same account",
			args: args{machineId: "m3", account: "950225", password: "123456"},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createAccount(tt.args.machineId, tt.args.account, tt.args.password); got != tt.want {
				t.Errorf("HandleStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Delete2(t *testing.T) {
	type args struct {
		machineId string
		account   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Delete Machine",
			args: args{machineId: "m1", account: "950223"},
			want: "OK",
		},
		{
			name: "Delete Machine",
			args: args{machineId: "m2", account: "950224"},
			want: "OK",
		},
		{
			name: "Delete Machine",
			args: args{machineId: "m3", account: "950225"},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := delete(tt.args.machineId, tt.args.account); got != tt.want {
				t.Errorf("HandleStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GenFtpPasswd(t *testing.T) {
	type args struct {
		account string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Gen ftp password",
			args: args{account: "950223"},
			want: "e9f9f7",
		},
		{
			name: "Gen ftp password",
			args: args{account: "A50503"},
			want: "0d54f2",
		},
		{
			name: "Gen ftp password",
			args: args{account: "A40503"},
			want: "a05ef2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenFtpPasswd(tt.args.account); got != tt.want {
				t.Errorf("GenFtpPasswd() = %v, want %v", got, tt.want)
			}
		})
	}
}
