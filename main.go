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

func getJob(jobName *string, namespaceName *string, clientset *kubernetes.Clientset) (job batch.Job) {
	jobs, err := clientset.BatchV1().Jobs(*namespaceName).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, foundJob := range jobs.Items {
		if *jobName == foundJob.ObjectMeta.GetName() {
			job = foundJob
		}
	}

	return job
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

	log.Print("Looking the job: " + *jobName + " in " + *namespaceName + " NS")
	myJob := getJob(jobName, namespaceName, clientset)

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
		if myJob.Status.Failed == 0 && myJob.Status.Succeeded == 0 && myJob.Status.Active == 0 && *waitToJob == false {
			log.Println("Job not found")
			break
		}
		log.Printf("Active: %d\n", myJob.Status.Active)
		log.Printf("Succeeded: %d\n", myJob.Status.Succeeded)
		log.Printf("Failed: %d\n", myJob.Status.Failed)
		log.Print(".")
		time.Sleep(10 * time.Second)
		myJob = getJob(jobName, namespaceName, clientset)

	}

	if myJob.Status.Succeeded == 1 {
		log.Println("Job success")
	}
}
