package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hobord/k8jobwatch/kubeclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	jobName := flag.String("j", "", "job name")
	namespaceName := flag.String("ns", "", "namespace")
	waitToJob := flag.Bool("w", false, "wait to job if (job is not exists)")
	deleteJob := flag.Bool("d", false, "delete job when exit")
	flag.Parse()
	if *jobName == "" {
		fmt.Println("Job name parameter is mandatory! Use -j=jobname ")
		os.Exit(2)
	}

	clientset := kubeclient.GetClientset()

	log.Printf("Looking the %s Job in %s NS\n", *jobName, *namespaceName)

	found := false
	myJob, err := clientset.BatchV1().Jobs(*namespaceName).Get(*jobName, metav1.GetOptions{})
	if (err != nil) {
		log.Printf("Can't get the job: %v\n", err)
		if (*waitToJob == false) {
			os.Exit(-1)		
		}
	}

	for {
		if myJob.UID != "" && found == false {
			found = true
			log.Printf("Found the job: %s\n", myJob.UID)
		}

		if myJob.Status.Failed > 0 {
			log.Printf("Job failed")
			os.Exit(-1)
		}
		if myJob.Status.Succeeded > 0 {
			break
		}
		if found {
			log.Printf("Active: %d, Succeeded: %d, Failed: %d", myJob.Status.Active, myJob.Status.Succeeded, myJob.Status.Failed)
		} 
		time.Sleep(10 * time.Second)
		myJob, err = clientset.BatchV1().Jobs(*namespaceName).Get(*jobName, metav1.GetOptions{})
		if (err != nil) {
			if (*waitToJob == false) {
				os.Exit(-1)		
			} else {
				log.Printf("Can't get the job: %v\n", err)
			}
		}
	}

	if myJob.Status.Succeeded == 1 {
		log.Printf("Job success")
		if *deleteJob {
			err := clientset.BatchV1().Jobs(*namespaceName).Delete(*jobName, new(metav1.DeleteOptions))
			if (err!=nil) {
				log.Printf("Can't delete the job: %v\n", err)
			} else {
				log.Printf("Job deleted")
			}
		}
	}
}
