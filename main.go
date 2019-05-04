package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	namespace       = flag.String("n", "default", "namespace")
	refreshSeconds  = flag.Int("r", 60, "refresh pods loop in seconds")
	refreshDuration = flag.Int("e", 10800, "rescan image after in seconds")
	listener        = flag.Bool("l", true, "start listener")
	registryFilter  = flag.String("rr", "", "registry filter")
	clairLocation   = flag.String("c", "clair", "clair endpoint")
	noop            = flag.Bool("o", false, "no op")
	images          map[string]*containerScan
	mutex           = &sync.Mutex{}
	jobs            chan string
	mockScanner     bool
	ipAddress       string
)

type containerScanResult struct {
	Containers []*containerScan
}

type containerScan struct {
	LastCheck   time.Time
	Image       string
	ScanStarted bool
	Report      containerVulnerabilityReport
}

type containerVulnerabilityReport struct {
	Unapproved      []string                     `json:"unapproved"`
	Vulnerabilities []containerVulnerabilityInfo `json:"vulnerabilities"`
}

type containerVulnerabilityInfo struct {
	FeatureName    string `json:"featurename"`
	FeatureVersion string `json:"featureversion"`
	Vulnerability  string `json:"vulnerability"`
	Namespace      string `json:"namespace"`
	Description    string `json:"description"`
	Link           string `json:"link"`
	Severity       string `json:"severity"`
	FixedBy        string `json:"fixedby"`
}

func main() {
	flag.Parse()
	log.SetFlags(log.Llongfile)

	var err error
	ipAddress, err = GetDefaultIP()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Current IP Address %s\n", ipAddress)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// dirty read
		var containerImages []*containerScan
		for _, v := range images {
			containerImages = append(containerImages, v)
		}

		imagesBytes, err := json.Marshal(containerScanResult{Containers: containerImages})
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(imagesBytes)
	})

	if *listener {
		go func() {
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()
	}

	initScanWorker()
	for {
		if !*noop {
			pods, err := clientset.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
			if err != nil {
				log.Println(err.Error())
				continue
			}

			for _, pod := range pods.Items {
				for _, container := range pod.Spec.Containers {
					jobs <- container.Image
				}
			}
		} else {
			log.Println("no op")
		}

		time.Sleep(time.Duration(*refreshSeconds) * time.Second)
	}
}

func initScanWorker() {
	jobs = make(chan string)
	go scanWorker()

	images = make(map[string]*containerScan)
}

func scanWorker() {
	scans := make(chan *containerScan)
	for {
		select {
		case scan := <-scans:
			scan.ScanStarted = false
			scan.LastCheck = time.Now().UTC()

			log.Printf("ContainerScan %s updated\n", scan.Image)

		case image := <-jobs:
			if images[image] == nil {
				images[image] = &containerScan{Image: image}
				log.Printf("Created new ContainerScan %s\n", image)
			}

			scan := images[image]
			timeResult := scan.LastCheck.IsZero() || time.Now().UTC().After(scan.LastCheck.Add(time.Duration(*refreshDuration)*time.Second))
			if mockScanner {
				log.Printf("Time result for %s - %v", image, timeResult)
			}
			if timeResult && !scan.ScanStarted {
				scan.ScanStarted = true

				go func(i string, s *containerScan) {
					log.Printf("Starting ContainerScan Job %s\n", i)
					result, err := scanContainer(i)
					s.Report = result
					if err != nil {
						log.Println(err)
					}

					scans <- s
				}(image, scan)
			}
		}
	}
}

func scanContainer(image string) (containerVulnerabilityReport, error) {
	mutex.Lock()
	log.Printf("ScanningContainer %s\n", image)
	defer mutex.Unlock()

	var err error
	var data containerVulnerabilityReport

	if len(*registryFilter) > 0 {
		if !strings.HasPrefix(image, *registryFilter) {
			log.Println("Skipping security scan")
			return data, err
		}
	}

	if mockScanner {
		time.Sleep(2 * time.Second)
		return data, err
	}

	out, err := exec.Command("docker", "pull", image).Output()
	log.Println(string(out))
	if err != nil {
		return data, err
	}

	if _, err := os.Stat("report.json"); err == nil {
		_ = os.Remove("report.json")
	}

	out, err = exec.Command("/usr/local/bin/clair-scanner", "-c", *clairLocation, "--ip", ipAddress, "-r", "report.json", "-t", "Medium", image).Output()
	log.Println(string(out))
	if _, err := os.Stat("report.json"); os.IsNotExist(err) {
		return data, fmt.Errorf("no report found for image %s\n", image)
	}

	out, err = ioutil.ReadFile("report.json")
	if err != nil {
		log.Println(err)
		return data, err
	}

	if err = json.Unmarshal(out, &data); err != nil {
		log.Println(err)
	}

	return data, err
}
