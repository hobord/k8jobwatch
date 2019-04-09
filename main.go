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

	log.Print("Looking the " + *jobName + " Job in " + *namespaceName + " NS")

	found := false
	myJob, err := clientset.BatchV1().Jobs(*namespaceName).Get(*jobName, metav1.GetOptions{})
	if (err != nil) {
		log.Println("Can't get the job:", err)
		if (*waitToJob == false) {
			os.Exit(-1)		
		}
	}

	for {
		if myJob.UID != "" && found == false {
			found = true
			log.Print("Found the job: " + myJob.UID)
		}

		if myJob.Status.Failed > 0 {
			log.Println("Job failed")
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
				log.Println("Can't get the job:", err)
			}
		}
	}

	if myJob.Status.Succeeded == 1 {
		log.Println("Job success")
		if *deleteJob {
			err := clientset.BatchV1().Jobs(*namespaceName).Delete(*jobName, new(metav1.DeleteOptions))
			if (err!=nil) {
				log.Println("Can't delete the job:", err)
			} else {
				log.Println("Job deleted")
			}
		}
	}
}
