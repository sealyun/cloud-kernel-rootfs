package build

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/ecs"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/logger"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/retry"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/utils"
	"github.com/sealyun/cloud-kernel-rootfs/pkg/vars"
	"time"
)

type _package interface {
	InitK8sServer() error
	RunK8sServer() error
	WaitImages() error
	SavePackage() error
}

func Package(k8sVersion string, gc bool) error {

	instance := ecs.NewCloud().New(1, false, true)
	if instance == nil {
		return errors.New("create ecs is error")
	}
	logger.Info("1. begin create ecs")
	var instanceInfo *ecs.CloudInstanceResponse
	if gc {
		defer func() {
			ecs.NewCloud().Delete(instance, 10)
		}()
	} else {
		defer func() {
			logger.Info("end. ecs instanceId: %s", instance)
		}()
	}

	if err := retry.Do(func() error {
		var err error
		logger.Debug("1. retry fetch ecs info " + instance[0])
		instanceInfo, err = ecs.NewCloud().Describe(instance[0])
		if err != nil {
			return err
		}
		if instanceInfo.PublicIP == "" {
			return errors.New("retry error")
		}
		if !instanceInfo.IsOk {
			return errors.New("retry error")
		}
		return nil
	}, 100, 1*time.Second, false); err != nil {
		return utils.ProcessError(err)
	}
	publicIP := instanceInfo.PublicIP
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	logger.Debug("2. connect ssh: " + publicIP)
	if err := retry.Do(func() error {
		var err error
		logger.Debug("2. retry test ecs ssh: " + publicIP)
		_, err = s.CmdAndError(publicIP, "ls /")
		if err != nil {
			return err
		} else {
			return nil
		}
	}, 20, 500*time.Millisecond, true); err != nil {
		return utils.ProcessError(err)
	}
	if err := vars.LoadVars(k8sVersion, publicIP, s); err != nil {
		return err
	}
	var k8s _package
	if utils.For120(k8sVersion) {
		//k8s = NewContainerdK8s(publicIP)
		return fmt.Errorf("当前不支持该版本%s", k8sVersion)
	} else {
		k8s = NewDockerK8s(publicIP, k8sVersion)
	}
	if k8s == nil {
		return utils.ProcessError(errors.New("k8s interface is nil"))
	}
	logger.Info("3. install k8s[ " + k8sVersion + " ] : " + publicIP)
	if err := k8s.InitK8sServer(); err != nil {
		return utils.ProcessError(err)
	}
	if err := k8s.RunK8sServer(); err != nil {
		return utils.ProcessError(err)
	}
	logger.Info("4. wait k8s[ " + k8sVersion + " ] pull all images: " + publicIP)
	if err := checkKubeStatus("4", publicIP, s, false); err != nil {
		return utils.ProcessError(err)
	}
	if err := s.CmdAsync(publicIP, "kubectl get pod -n kube-system"); err != nil {
		return utils.ProcessError(err)
	}
	if err := k8s.WaitImages(); err != nil {
		return utils.ProcessError(err)
	}
	logger.Info("5. k8s[ " + k8sVersion + " ] image save: " + publicIP)
	if err := k8s.SavePackage(); err != nil {
		return utils.ProcessError(err)
	}
	logger.Info("6. k8s[ " + k8sVersion + " ] testing: " + publicIP)
	//if err = test(publicIP, k8sVersion); err != nil {
	//	return utils.ProcessError(err)
	//} else {
	//	logger.Info("6. k8s[ " + k8sVersion + " ] uploading: " + publicIP)
	//	//upload(publicIP, k8sVersion)
	//}
	logger.Info("7. k8s[ " + k8sVersion + " ] finished. " + publicIP)
	return nil
}
