package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hobord/k8jobwatch/kubeclient"
	batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getJob(jobName *string, namespaceName *string, clientset *kubernetes.Clientset, waitToJob *bool) (job batch.Job) {
	log.Print("Looking the job: " + *jobName + " in " + *namespaceName + " NS")
	for {
		jobs, err := clientset.BatchV1().Jobs(*namespaceName).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, foundJob := range jobs.Items {
			if *jobName == foundJob.ObjectMeta.GetName() {
				job = foundJob
				*waitToJob = false
			}
		}
		if !*waitToJob {
			return
		}
		log.Print(".")
		time.Sleep(10 * time.Second)
	}
}

func main() {
	jobName := flag.String("j", "", "job name")
	namespaceName := flag.String("n", "", "namespace")
	waitToJob := flag.Bool("w", false, "wait to job if (job is not exists)")
	flag.Parse()
	if *jobName == "" {
		fmt.Println("Job name parameter is mandatory! Use -j=jobname ")
		os.Exit(2)
	}

	clientset := kubeclient.GetClientset()

	myJob := getJob(jobName, namespaceName, clientset, waitToJob)

	if myJob.Status.Succeeded != 0 {
		log.Print("Found: " + myJob.UID)
		log.Print("Waiting to job success.")
		if myJob.Status.Failed == 1 {
			panic("Job failed")
		}
		for myJob.Status.Succeeded != 1 {
			if myJob.Status.Failed == 1 {
				panic("Job failed")
			}
			log.Print(".")
			time.Sleep(10 * time.Second)
		}
		log.Print("Success.")
	} else {
		fmt.Println("Job not found")
	}
}
