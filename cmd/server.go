package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/FXAZfung/image-board/cmd/flags"
	"github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/FXAZfung/image-board/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// ServerCmd represents the server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server at the specified address",
	Long: `Start the server at the specified address
the address is defined in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()
		if config.Conf.DelayedStart != 0 {
			utils.Log.Infof("delayed start for %d seconds", config.Conf.DelayedStart)
			time.Sleep(time.Duration(config.Conf.DelayedStart) * time.Second)
		}
		if !flags.Debug && !flags.Dev {
			gin.SetMode(gin.ReleaseMode)
		}
		r := gin.New()
		r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
		server.Init(r)
		var httpSrv, httpsSrv, unixSrv *http.Server
		if config.Conf.Scheme.HttpPort != -1 {
			httpBase := fmt.Sprintf("%s:%d", config.Conf.Scheme.Address, config.Conf.Scheme.HttpPort)
			utils.Log.Infof("start HTTP server @ %s", httpBase)
			httpSrv = &http.Server{Addr: httpBase, Handler: r}
			go func() {
				err := httpSrv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					utils.Log.Fatalf("failed to start http: %s", err.Error())
				}
			}()
		}
		if config.Conf.Scheme.HttpsPort != -1 {
			httpsBase := fmt.Sprintf("%s:%d", config.Conf.Scheme.Address, config.Conf.Scheme.HttpsPort)
			utils.Log.Infof("start HTTPS server @ %s", httpsBase)
			httpsSrv = &http.Server{Addr: httpsBase, Handler: r}
			go func() {
				err := httpsSrv.ListenAndServeTLS(config.Conf.Scheme.CertFile, config.Conf.Scheme.KeyFile)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					utils.Log.Fatalf("failed to start https: %s", err.Error())
				}
			}()
		}
		if config.Conf.Scheme.UnixFile != "" {
			utils.Log.Infof("start unix server @ %s", config.Conf.Scheme.UnixFile)
			unixSrv = &http.Server{Handler: r}
			go func() {
				listener, err := net.Listen("unix", config.Conf.Scheme.UnixFile)
				if err != nil {
					utils.Log.Fatalf("failed to listen unix: %+v", err)
				}
				// set socket file permission
				mode, err := strconv.ParseUint(config.Conf.Scheme.UnixFilePerm, 8, 32)
				if err != nil {
					utils.Log.Errorf("failed to parse socket file permission: %+v", err)
				} else {
					err = os.Chmod(config.Conf.Scheme.UnixFile, os.FileMode(mode))
					if err != nil {
						utils.Log.Errorf("failed to chmod socket file: %+v", err)
					}
				}
				err = unixSrv.Serve(listener)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					utils.Log.Fatalf("failed to start unix: %s", err.Error())
				}
			}()
		}
		//if config.Conf.S3.Port != -1 && config.Conf.S3.Enable {
		//	s3r := gin.New()
		//	s3r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
		//	server.InitS3(s3r)
		//	s3Base := fmt.Sprintf("%s:%d", config.Conf.Scheme.Address, config.Conf.S3.Port)
		//	utils.Log.Infof("start S3 server @ %s", s3Base)
		//	go func() {
		//		var err error
		//		if config.Conf.S3.SSL {
		//			httpsSrv = &http.Server{Addr: s3Base, Handler: s3r}
		//			err = httpsSrv.ListenAndServeTLS(config.Conf.Scheme.CertFile, config.Conf.Scheme.KeyFile)
		//		}
		//		if !config.Conf.S3.SSL {
		//			httpSrv = &http.Server{Addr: s3Base, Handler: s3r}
		//			err = httpSrv.ListenAndServe()
		//		}
		//		if err != nil && !errors.Is(err, http.ErrServerClosed) {
		//			utils.Log.Fatalf("failed to start s3 server: %s", err.Error())
		//		}
		//	}()
		//}
		//var ftpDriver *server.FtpMainDriver
		//var ftpServer *ftpserver.FtpServer
		//if config.Conf.FTP.Listen != "" && config.Conf.FTP.Enable {
		//	var err error
		//	ftpDriver, err = server.NewMainDriver()
		//	if err != nil {
		//		utils.Log.Fatalf("failed to start ftp driver: %s", err.Error())
		//	} else {
		//		utils.Log.Infof("start ftp server on %s", conf.Conf.FTP.Listen)
		//		go func() {
		//			ftpServer = ftpserver.NewFtpServer(ftpDriver)
		//			err = ftpServer.ListenAndServe()
		//			if err != nil {
		//				utils.Log.Fatalf("problem ftp server listening: %s", err.Error())
		//			}
		//		}()
		//	}
		//}
		//var sftpDriver *server.SftpDriver
		//var sftpServer *sftpd.SftpServer
		//if conf.Conf.SFTP.Listen != "" && conf.Conf.SFTP.Enable {
		//	var err error
		//	sftpDriver, err = server.NewSftpDriver()
		//	if err != nil {
		//		utils.Log.Fatalf("failed to start sftp driver: %s", err.Error())
		//	} else {
		//		utils.Log.Infof("start sftp server on %s", conf.Conf.SFTP.Listen)
		//		go func() {
		//			sftpServer = sftpd.NewSftpServer(sftpDriver)
		//			err = sftpServer.RunServer()
		//			if err != nil {
		//				utils.Log.Fatalf("problem sftp server listening: %s", err.Error())
		//			}
		//		}()
		//	}
		//}
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 1 second.
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscanll.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		utils.Log.Println("Shutdown server...")
		Release()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		var wg sync.WaitGroup
		if config.Conf.Scheme.HttpPort != -1 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := httpSrv.Shutdown(ctx); err != nil {
					utils.Log.Fatal("HTTP server shutdown err: ", err)
				}
			}()
		}
		if config.Conf.Scheme.HttpsPort != -1 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := httpsSrv.Shutdown(ctx); err != nil {
					utils.Log.Fatal("HTTPS server shutdown err: ", err)
				}
			}()
		}
		if config.Conf.Scheme.UnixFile != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := unixSrv.Shutdown(ctx); err != nil {
					utils.Log.Fatal("Unix server shutdown err: ", err)
				}
			}()
		}
		//if config.Conf.FTP.Listen != "" && config.Conf.FTP.Enable && ftpServer != nil && ftpDriver != nil {
		//	wg.Add(1)
		//	go func() {
		//		defer wg.Done()
		//		ftpDriver.Stop()
		//		if err := ftpServer.Stop(); err != nil {
		//			utils.Log.Fatal("FTP server shutdown err: ", err)
		//		}
		//	}()
		//}
		//if config.Conf.SFTP.Listen != "" && config.Conf.SFTP.Enable && sftpServer != nil && sftpDriver != nil {
		//	wg.Add(1)
		//	go func() {
		//		defer wg.Done()
		//		if err := sftpServer.Close(); err != nil {
		//			utils.Log.Fatal("SFTP server shutdown err: ", err)
		//		}
		//	}()
		//}
		wg.Wait()
		utils.Log.Println("Server exit")
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// OutImageBoardInit 暴露用于外部启动server的函数
func OutImageBoardInit() {
	var (
		cmd  *cobra.Command
		args []string
	)
	ServerCmd.Run(cmd, args)
}
