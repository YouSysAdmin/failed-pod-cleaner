package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yousysadmin/failed-pod-cleaner/pkg/kube"
	"github.com/yousysadmin/failed-pod-cleaner/pkg/logging"
	"github.com/yousysadmin/failed-pod-cleaner/pkg/osenv"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	log     = logging.New()
	timeout = parseDuration(osenv.GetEnv("TIMEOUT", "60"))
)

func main() {
	// connect to the cluster
	kc, err := kube.Connect()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// make context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// make signal channel for handling gracefully stopping
	closing := make(chan bool)

	// make wait group
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context, closing chan bool) {
		defer wg.Done()

		// initialise timer
		ticker := time.NewTicker(timeout)

		// logic
		for {
			select {
			case <-ticker.C:
				pods, err := kube.ListPods(kc, "") // get all pods from all namespace
				if err != nil {
					log.Error(err.Error())
				}
				cleanup(kc, pods)
			case <-closing:
				ticker.Stop()
				return
			}
		}
	}(ctx, closing)

	log.Info("pod cleaner is running.")

	// make channel for handling os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// wait for stop signal
	<-c

	// initiate graceful shutdown
	close(closing)

	// wait until workers are done
	wg.Wait()

	log.Info("quit")
}

// Clearing dead Containers
func cleanup(clientset *kubernetes.Clientset, podsList *v1.PodList) {
	for _, p := range podsList.Items {
		if p.Status.Phase == "Failed" { // only pods with "Failed" status
			var result string
			err := kube.DeletePod(clientset, p.Namespace, p.Name)
			if err != nil {
				result = err.Error()
			} else {
				result = "removed"
			}
			// output log
			log.WithFields(logrus.Fields{
				"pod_name":      p.Name,
				"pod_namespace": p.Namespace,
				"pod_status":    p.Status.Phase,
				"pod_reason":    p.Status.Reason,
				"pod_node":      p.Spec.NodeName,
				"pod_message":   p.Status.Message,
			}).Info(result)
		}
	}
}

// parse timeout value from the environments
func parseDuration(value string) time.Duration {
	v, _ := strconv.Atoi(value)
	return time.Duration(v) * time.Minute
}
