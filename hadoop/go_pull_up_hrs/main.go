package main

import (
	"bytes"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"go_pull_up_hrs/pkg/sshclient"
	"sync"
	"time"
)

var (
	cfg      *ProgramConfig
	pool, _  = ants.NewPool(100)
	log      = logging.Logger("pull_up_hrs")
	interval time.Duration
)

func (h HostConfig) isHMasterRunning() bool {
	if h.Master.SSHPort == "" {
		log.Debug("master ssh port undefined, set 22 to default")
		h.Master.SSHPort = "22"
	}
	HMasterStatus := false
	for _, addr := range h.Master.Targets {
		client, err := sshclient.DialWithPasswd(addr+":"+h.Master.SSHPort, h.Master.User, h.Master.Passwd)
		if err != nil {
			log.Errorf("dial %s+ failed: %s", addr, err)
			continue
		}
		cmd := `ps -ef | grep org.apache.hadoop.hbase.master.HMaster | grep -v grep | wc -l`
		data, err := client.Cmd(cmd).Output()
		if err != nil {
			log.Errorf("invoke (%s) failed, %s", cmd, err)
			_ = client.Close()
			continue
		}
		retVal := string(bytes.TrimSpace(data))
		HMasterStatus = HMasterStatus || retVal == "1"
		_ = client.Close()
	}
	return HMasterStatus
}

func (h HostConfig) checkOrPullUpHRegion() {
	if h.Slave.SSHPort == "" {
		log.Debug("slave ssh port undefined, set 22 to default")
		h.Slave.SSHPort = "22"
	}
	var wg sync.WaitGroup
	for _, addr := range h.Slave.Targets {
		if addr == "" {
			continue
		}
		wg.Add(1)
		addr := addr
		_func := func() {
			defer wg.Done()
			client, err := sshclient.DialWithPasswd(addr+":"+h.Slave.SSHPort, h.Slave.User, h.Slave.Passwd)
			if err != nil {
				log.Errorf("dial %s+ failed: %s", addr, err)
				return
			}
			defer func(client *sshclient.Client) {
				_ = client.Close()
			}(client)
			if err := hasHRegionProcess(client); err == nil {
				return
			}
			wait := interval / 2
			msg := fmt.Sprintf("检测到 %s 的HRegionServer离线, %s 后将尝试拉起", addr, time.Duration(wait))
			//sendToFeishuBot(msg)
			log.Error(msg)
			// try to pull up hbase region server
			time.Sleep(wait)
			// try again check hbase region server process
			if err := hasHRegionProcess(client); err == nil {
				// hbase region server is resolved; skipped...
				//sendToFeishuBot(fmt.Sprintf("%s HRegionServer 已恢复, 取消拉起操作", addr))
				log.Infof("%s HRegionServer is recoverd", addr)
				return
			}
			if err := startHRegionServer(client, h.HBaseDaemonScript); err != nil {
				sendToFeishuBot(fmt.Sprintf("%s HRegionServer 拉起失败, 原因: %s", addr, err))
				log.Error(err)
				return
			}
			//sendToFeishuBot(fmt.Sprintf("%s HRegionServer 拉起成功", addr))
		}
		if err := pool.Submit(_func); err != nil {
			log.Errorf("_func submit failed, %s", err)
		}
	}
	wg.Wait()
}

func startHRegionServer(_ssh *sshclient.Client, bootScript string) error {
	if _ssh == nil || bootScript == "" {
		return errors.New("_ssh or boot_script is nil")
	}
	cmd := fmt.Sprintf("%s start regionserver", bootScript)
	_, err := _ssh.Cmd(cmd).Output()
	if err != nil {
		return errors.Errorf("start hbase region server failed, %s", err)
	}
	return nil
}

func hasHRegionProcess(_ssh *sshclient.Client) error {
	if _ssh == nil {
		return errors.New("_ssh is nil")
	}
	cmd := `ps -ef | grep -v grep | grep org.apache.hadoop.hbase.regionserver.HRegionServer | wc -l`
	data, err := _ssh.Cmd(cmd).Output()
	if err != nil {
		return errors.Errorf("find hbase region server process error, %s", err)
	}
	val := string(bytes.TrimSpace(data))
	if val == "1" {
		return nil
	}
	if val == "0" {
		return errors.New("HRegionServer process not found")
	}
	return nil
}

func main() {
	defer pool.Release()
	logging.SetAllLoggers(logging.LevelDebug)

	for {
		for name, host := range cfg.CheckOrPullUp.HbaseRegionServer {
			if !host.isHMasterRunning() {
				log.Infof("%s    HMaster 进程处于离线状态, 跳过HRegionServer的检查", name)
				continue
			}
			log.Infof("%s    HMaster 进程处于在线状态, 进行HRegionServer进程检查", name)
			host.checkOrPullUpHRegion()
		}
		time.Sleep(interval)
	}
}
