package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/shubhamagarwal19/pod-crash-notification-controller/pkg/slack"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type ReconcilePod struct {
	Client client.Client
}

var (
	restartList map[string]int32
)

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	log.Println("NameSpace: ", viper.GetString("namespace"))
	log.Println("Slack Webhook: ", viper.GetString("slack-webhook"))
	restartList = make(map[string]int32)
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	mgr, err := manager.New(config, manager.Options{Namespace: viper.GetString("namespace")})
	if err != nil {
		log.Fatal(err)
	}

	ctrl, err := controller.New("pod-crash-notification-controller", mgr, controller.Options{
		Reconciler: &ReconcilePod{
			Client: mgr.GetClient(),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching kubernetes Pods .... ")
	if err = ctrl.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting the POD CRASH NOTIFICATION Controller ....")
	if err = mgr.Start(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

func (r *ReconcilePod) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	pod := &corev1.Pod{}
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Name}, pod)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("failed to retrieve pod, might be deleted")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	for i := range pod.Status.ContainerStatuses {
		container := pod.Status.ContainerStatuses[i].Name
		restartCount := pod.Status.ContainerStatuses[i].RestartCount
		identifier := pod.Name + pod.Status.ContainerStatuses[i].Name

		if _, ok := restartList[identifier]; !ok {
			restartList[identifier] = restartCount
		} else if restartList[identifier] < restartCount {
			log.Println("Reconciling .... ")
			log.Printf("Received event for container[%v] pod[%v]", container, pod.Name)
			msg := slack.SlackRequestBody {
				Text: fmt.Sprintf("Container[%s] inside pod[%s] crashed ... ", container, pod.Name),
			}
			// send Slack notification
			err = slack.SendSlackNotification(viper.GetString("slack-webhook"), msg)
			if err != nil {
				log.Println("Error: ", err)
			} else {
				restartList[identifier] = restartCount
			}
		}
	}

	return reconcile.Result{}, nil
}