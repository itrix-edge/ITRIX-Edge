//Changelog
// 2017/8/30
// 1.Remove sqlite db and deleting scheduling
// 2.Simplify create and delete algorithm
//2017/9/14
//1.Adding log path configuration
//2017/11/30
//Fix XFS on NFS export
//http://blog.erben.sk/2015/01/25/xfs-filesystem-nfs-export-mount-fails-with-stale-nfs-handle/

package main

import (
	"crypto/md5"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Build = "xxxxx"

const Version = "v1.3"

var dataFolder string

var expireDuration time.Duration

var db *sql.DB
var accountLock string

//Logger
var log = logging.MustGetLogger("main")

type Message struct {
	Message string `json:"message"`
}

type JobInfo struct {
	Account string `json:"account"`
}

func main() {
	//Parsing flag
	portPtr := flag.String("port", "8888", "Port for storage API server")
	dataFolderPtr := flag.String("data-folder", "/nfs/", "Data folder that shared with NFS/FTP")
	logPathPtr := flag.String("logpath", "/var/log", "Log path")
	dataFolder = *dataFolderPtr
	debugPtr := flag.Bool("debug", false, "Run in debug mode")
	flag.Parse()

	//init log
	initLogSetting(*logPathPtr, logging.DEBUG)
	log.Infof("=======Initialize API server %v %v=======", Version, Build)
	log.Debugf("Port: %v, Data Folder: %v, debug Mode: %v", *portPtr, *dataFolderPtr, *debugPtr)
	if !*debugPtr {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	//POST (account/password)
	r.POST("/storage", CreateHandler)
	r.GET("/storage/:account", GetHandler)
	//DELETE (R,D)
	r.DELETE("/storage/:account", DeleteHandler)

	log.Fatal(r.Run(":" + *portPtr))
}

//init log
func initLogSetting(logPath string, level logging.Level) {
	//Setup Console format
	var consoleFormat = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{level:.4s} > %{shortfile} [%{shortfunc}] %{message}%{color:reset}`,
	)

	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormat)
	consoleBackendLeveled := logging.AddModuleLevel(consoleBackendFormatter)
	consoleBackendLeveled.SetLevel(level, "")

	//Setup file format
	var fileFormat = logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05.000} %{level:.4s} > %{shortfile} [%{shortfunc}] %{message}`,
	)

	fileName := logPath + "/" + filepath.Base(os.Args[0]) + ".log"

	//create log rolling
	f := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
	fileBackend := logging.NewLogBackend(f, "", 0)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)
	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatter)
	fileBackendLeveled.SetLevel(level, "")

	// Set the backends to be used.
	logging.SetBackend(consoleBackendLeveled, fileBackendLeveled)
}

//get MD5 password
func GenFtpPasswd(account string) string {
	data := []byte(account + "ITRIDNN")
	password := fmt.Sprintf("%x", md5.Sum(data))[:6]
	log.Debugf("Generate Ftp password: %v", password)
	return password
}

//get MD5 folder name
func GenNFSFolder(account string) string {
	data := []byte(account + "ITRI+DNN+NFS Folder")
	folder := fmt.Sprintf("%x", md5.Sum(data))
	log.Debugf("Generate NFS folder: %v", folder)
	return folder
}

//recycle storage
func RecycleStorage(account string) {
	log.Infof("Recycle storage: Account %v", account)
	accountLock = account
	log.Infof("Remove FTP, NFS, and folder %v", account)
	//remove ftp account
	//sed -i '/^mac1:.*/d' /nfs/passwd
	cmd := exec.Command("/bin/bash", "-c", "sed -i '/^"+account+".*/d' "+dataFolder+"passwd")
	err := cmd.Run()
	if err != nil {
		log.Errorf("delete ftp account "+account+" failed: %v", err)
		return
	}

	//remove nfs server
	folder := GenNFSFolder(account)
	cmdString := "sed -i '/^" + strings.Replace(dataFolder+folder, "/", `\/`, -1) + ".*/d' " + dataFolder + "exports"
	log.Debug(cmdString)
	cmd = exec.Command("/bin/bash", "-c", cmdString)
	err = cmd.Run()
	if err != nil {
		log.Errorf("delete NFS account "+account+" failed: %v", err)
		return
	}
	cmd = exec.Command("/bin/bash", "-c", "exportfs -ra")
	err = cmd.Run()
	if err != nil {
		log.Errorf("exportfs config failed: %v", err)
		return
	}

	//delete user folder
	err = os.RemoveAll(dataFolder + folder)
	if err != nil {
		log.Errorf("delete folder "+folder+" failed: %v", err)
		return
	}
	err = os.RemoveAll(dataFolder + account)
	if err != nil {
		log.Errorf("delete folder "+account+" failed: %v", err)
		return
	}
	accountLock = ""
	log.Infof("--------Recycle account: Account %v-----DONE----", account)
}

//SingleHandler
func CreateHandler(c *gin.Context) {

	log.Info("=====Create Account=====")

	// Decode
	log.Info("Decoding JSON")
	var jobInfo JobInfo
	decodeErr := c.BindJSON(&jobInfo)
	if decodeErr != nil {
		log.Errorf("Decode failed")
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Json formate error"})
		return
	}

	//Check empty
	log.Info("Check JSON data")
	if jobInfo.Account == "" {
		log.Errorf("Missing parameter error: v%", jobInfo)
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Missing parameter error"})
		return
	}
	for strings.Compare(accountLock, jobInfo.Account) == 0 {
		log.Debugf("Wait for unlock: %v", accountLock)
		time.Sleep(time.Second * 1)
	}
	//check account is create or not
	folder := GenNFSFolder(jobInfo.Account)
	log.Infof("Create Account: %v", jobInfo.Account)
	if _, err := os.Stat(dataFolder + folder); os.IsNotExist(err) {

		//Create Folder
		createFolder(dataFolder + folder)

		//Create NFS configuration
		log.Info("Create NFS configuration")
		cmd := exec.Command("/bin/bash", "-c", "A=$((`grep fsid "+dataFolder+"exports | awk -F',|=' {'print $3'} | sort -n | tail -n 1` +1));B=$(grep fsid "+dataFolder+"exports | awk -F',|=' {'print $3'} | sort -n | awk 'p && p != $1 { for( i = p; i < $1; i++ ) print i; } {p = $1 + 1 }' | head -1);if [ $B ]; then A=$B;fi;echo \""+dataFolder+folder+" *(rw,fsid=$A,no_subtree_check,sync,all_squash,anonuid=0,anongid=0)\" >> "+dataFolder+"exports;exportfs -ra")
		err := cmd.Start()
		if err != nil {
			log.Errorf("Create NFS config failed: %v", err)
			return
		}
		// cmd = exec.Command("/bin/bash", "-c", "exportfs -ra")
		// err = cmd.Start()
		// if err != nil {
		// 	log.Errorf("exportfs config failed: %v", err)
		// 	return
		// }
		cmd = exec.Command("/bin/bash", "-c", "ln -s "+dataFolder+folder+" "+dataFolder+jobInfo.Account)
		err = cmd.Start()
		if err != nil {
			log.Errorf("createn ln failed: %v", err)
			return
		}

		//Create FTP configuration
		log.Info("Create FTP configuration")
		cmd = exec.Command("/bin/bash", "-c", "echo "+jobInfo.Account+":$(openssl passwd -1 "+GenFtpPasswd(jobInfo.Account)+") >> "+dataFolder+"passwd; chmod 777 -R "+dataFolder+jobInfo.Account)
		err = cmd.Start()
		if err != nil {
			log.Errorf("Create FTP config failed: %v", err)
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Error Account existed"})
	}
}

func GetHandler(c *gin.Context) {
	log.Info("=====Get Account=====")

	//account
	account := c.Params.ByName("account")
	log.Debugf("Account: " + account)
	if account == "" {
		log.Error("Error account is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Error account is empty"})
		return
	}
	if _, err := os.Stat(dataFolder + account); os.IsNotExist(err) {
		log.Error("Error account doesn't exist")
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Error account doesn't exist"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "Account exist"})
	}

}

func DeleteHandler(c *gin.Context) {
	log.Info("=====DELETE Account=====")

	//account
	account := c.Params.ByName("account")
	log.Debugf("Account: " + account)
	if account == "" {
		log.Error("Error account is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Error account is empty"})
		return
	}
	folder := GenNFSFolder(account)
	if _, err := os.Stat(dataFolder + folder); os.IsNotExist(err) {
		log.Error("Error account doesn't exist")
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Error account doesn't exist"})
	} else {
		RecycleStorage(account)
	}

}

//create folder
func createFolder(folderName string) {
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		log.Infof("Create new folder %v", folderName)
		os.MkdirAll(folderName, 0777)
	}
}
