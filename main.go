package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/coreos/clair/api/v1"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/rest"
)

var (
	tokenPath       = flag.String("t", "/var/run/secrets/kubernetes.io/serviceaccount/token", "token path")
	server          = flag.String("s", "https://kubernetes.default", "server name")
	namespace       = flag.String("n", "default", "namespace")
	refreshSeconds  = flag.Int("r", 60, "refresh pods loop in seconds")
	refreshDuration = flag.Int("e", 10800, "rescan image after in seconds")
	listener        = flag.Bool("l", true, "start listener")
	clairTest       = flag.Bool("c", false, "clair test")
	clairTestImage  = flag.String("cc", "nginx", "clair test image")
	registryFilter  = flag.String("rr", "", "registry filter")
	noop            = flag.Bool("o", false, "no op")
	images          map[string]*containerScan
	mutex           = &sync.Mutex{}
	jobs            chan string
	mockScanner     bool
)

type containerScanResult struct {
	Containers []*containerScan
}

type containerScan struct {
	LastCheck       time.Time
	Vulnerabilities []v1.Vulnerability
	Image           string
	ScanStarted     bool
}

func main() {
	flag.Parse()
	ipAddress, err := GetDefaultIP()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Current IP Address %s\n", ipAddress)

	if *clairTest {
		data, err := scanContainer(*clairTestImage)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range data {
			log.Printf("Severity: %s Name: %s\n\t\t Description %s\n", v.Severity, v.Name, v.Description)
		}

		return
	}

	token, err := ioutil.ReadFile(*tokenPath)
	if err != nil {
		log.Fatal(err)
	}

	config := &rest.Config{
		Host:        *server,
		BearerToken: string(token),
		Insecure:    true,
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
			pods, err := clientset.Core().Pods(*namespace).List(api.ListOptions{})
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
					s.Vulnerabilities = result
					if err != nil {
						log.Println(err)
					}

					scans <- s
				}(image, scan)
			}
		}
	}
}

func scanContainer(image string) ([]v1.Vulnerability, error) {
	mutex.Lock()
	log.Printf("ScanningContainer %s\n", image)
	defer mutex.Unlock()

	var err error
	var data []v1.Vulnerability

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

	out, err = exec.Command("analyze-local-images", "-json", "-minimum-severity", "High", image).Output()
	log.Println(string(out))
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(out, &data)

	return data, err
}
