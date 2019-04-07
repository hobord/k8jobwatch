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
		if *waitToJob == false {
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

	// fmt.Println(myJob.UID)
	// fmt.Println(myJob.Status.Active)
	// fmt.Println(myJob.Status.Succeeded)
	// fmt.Println(myJob.Status.Failed)
	// fmt.Println(myJob.Status.String())

	if myJob.UID != "" {
		log.Print("Found: " + myJob.UID)
	}

	for {
		if myJob.Status.Failed == 1 {
			panic("Job failed")
		}
		if myJob.Status.Succeeded == 1 {
			break
		}
		log.Printf("Active: %d\n", myJob.Status.Active)
		log.Printf("Succeeded: %d\n", myJob.Status.Succeeded)
		log.Printf("Failed: %d\n", myJob.Status.Failed)
		log.Print(".")
		time.Sleep(10 * time.Second)
		myJob = getJob(jobName, namespaceName, clientset, waitToJob)

	}

	if myJob.Status.Succeeded == 1 {
		log.Println("Job success")
	}
}
