# Kate, a friend of Clair
## Summary
> [CoreOS Clair](https://github.com/coreos/clair) is an open source project for the static analysis of vulnerabilities in application containers (currently including appc and docker)

Turns out if you throw [CoreOS Clair](https://github.com/coreos/clair) into your Kubernetes cluster, with the help of a friend, Kate will automatically scan all newly launched containers.

Kate will also rescan all the images every couple of hours just to let you know if the CVE situation has changed.

This allows you to identify vulnerabilities that exist in production, as apposed to fixes that may exit on your upstream platforms.

## Example Deployment
Create a dedicated namespace for clair and switch context to the clair namespace ([Helm Clair Chart](https://github.com/coreos/clair/blob/master/Documentation/running-clair.md#kubernetes-helm))

Swith to namespace where you want to deploy kate and modify the rbac and arguments for kate ([./example/kate.yaml](./example/kate.yaml))
```bash
  kubectl create -f ./example/kate.yaml
```

- Start an example NGINX container
```bash
kubectl run nginx --image=nginx --restart=Always
```

- Find the pod name for Clair
```bash
watch kubectl get pods -o wide
```

- Monitor and wait for the Clair CVE crawl to finish. Look for something like an 'updater: update finished' message.
```bash
kubectl logs kate-66ayv -f
```

- Access your kate service on http://127.0.0.1:8080
```bash
kubectl port-forward kate-1jzhx 8080
```

- Integrate the endpoint into a UI or your monitoring tools

## TODO
- Build an additional UI Visualization Dashboard
- Prometheous integration (for alert manager)

## API output http://127.0.0.1:8080/
```json
{
  "Containers": [
    {
      "LastCheck": "2017-01-26T19:48:23.757871838Z",
      "Vulnerabilities": [],
      "Image": "andrewwebber/kate",
      "ScanStarted": false
    },
    {
      "LastCheck": "2017-01-26T19:47:28.267409008Z",
      "Vulnerabilities": [
        {
          "Name": "CVE-2015-8964",
          "NamespaceName": "debian:8",
          "Description": "The tty_set_termios_ldisc function in drivers/tty/tty_ldisc.c in the Linux kernel before 4.5 allows local users to obtain sensitive information from kernel memory by reading a tty data structure.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2015-8964",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.1,
                "Vectors": "AV:N/AC:M/Au:N/C:C/I:N"
              }
            }
          },
          "FixedBy": "3.16.39-1"
        },
        {
          "Name": "CVE-2016-0758",
          "NamespaceName": "debian:8",
          "Description": "Integer overflow in lib/asn1_decoder.c in the Linux kernel before 4.6 allows local users to gain privileges via crafted ASN.1 data.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-0758",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.2,
                "Vectors": "AV:L/AC:L/Au:N/C:C/I:C"
              }
            }
          }
        },
        {
          "Name": "CVE-2016-7543",
          "NamespaceName": "debian:8",
          "Description": "Bash before 4.4 allows local users to execute arbitrary commands with root privileges via crafted SHELLOPTS and PS4 environment variables.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-7543",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.2,
                "Vectors": "AV:L/AC:L/Au:N/C:C/I:C"
              }
            }
          },
          "FixedBy": "4.3-11+deb8u1"
        },
        {
          "Name": "CVE-2016-8399",
          "NamespaceName": "debian:8",
          "Description": "An elevation of privilege vulnerability in the kernel networking subsystem could enable a local malicious application to execute arbitrary code within the context of the kernel. This issue is rated as Moderate because it first requires compromising a privileged process and current compiler optimizations restrict access to the vulnerable code. Product: Android. Versions: Kernel-3.10, Kernel-3.18. Android ID: A-31349935.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-8399",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.6,
                "Vectors": "AV:N/AC:H/Au:N/C:C/I:C"
              }
            }
          },
          "FixedBy": "3.16.39-1"
        },
        {
          "Name": "CVE-2016-7910",
          "NamespaceName": "debian:8",
          "Description": "Use-after-free vulnerability in the disk_seqf_stop function in block/genhd.c in the Linux kernel before 4.7.1 allows local users to gain privileges by leveraging the execution of a certain stop operation even if the corresponding start operation had failed.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-7910",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 9.3,
                "Vectors": "AV:N/AC:M/Au:N/C:C/I:C"
              }
            }
          },
          "FixedBy": "3.16.39-1"
        },
        {
          "Name": "CVE-2016-7911",
          "NamespaceName": "debian:8",
          "Description": "Race condition in the get_task_ioprio function in block/ioprio.c in the Linux kernel before 4.6.6 allows local users to gain privileges or cause a denial of service (use-after-free) via a crafted ioprio_get system call.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-7911",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 9.3,
                "Vectors": "AV:N/AC:M/Au:N/C:C/I:C"
              }
            }
          },
          "FixedBy": "3.16.39-1"
        },
        {
          "Name": "CVE-2016-2090",
          "NamespaceName": "debian:8",
          "Description": "Off-by-one vulnerability in the fgetwln function in libbsd before 0.8.2 allows attackers to have unspecified impact via unknown vectors, which trigger a heap-based buffer overflow.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-2090",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.5,
                "Vectors": "AV:N/AC:L/Au:N/C:P/I:P"
              }
            }
          }
        },
        {
          "Name": "CVE-2016-0494",
          "NamespaceName": "debian:8",
          "Description": "Unspecified vulnerability in the Java SE and Java SE Embedded components in Oracle Java SE 6u105, 7u91, and 8u66 and Java SE Embedded 8u65 allows remote attackers to affect confidentiality, integrity, and availability via unknown vectors related to 2D.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-0494",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 10,
                "Vectors": "AV:N/AC:L/Au:N/C:C/I:C"
              }
            }
          },
          "FixedBy": "52.1-8+deb8u4"
        },
        {
          "Name": "CVE-2016-6293",
          "NamespaceName": "debian:8",
          "Description": "The uloc_acceptLanguageFromHTTP function in common/uloc.cpp in International Components for Unicode (ICU) through 57.1 for C/C++ does not ensure that there is a '\\0' character at the end of a certain temporary array, which allows remote attackers to cause a denial of service (out-of-bounds read) or possibly have unspecified other impact via a call with a long httpAcceptLanguage argument.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-6293",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.5,
                "Vectors": "AV:N/AC:L/Au:N/C:P/I:P"
              }
            }
          },
          "FixedBy": "52.1-8+deb8u4"
        }
      ],
      "Image": "quay.io/coreos/clair:v1.2.6",
      "ScanStarted": true
    },
    {
      "LastCheck": "2017-01-26T19:48:31.247263772Z",
      "Vulnerabilities": [
        {
          "Name": "CVE-2016-2090",
          "NamespaceName": "debian:8",
          "Description": "Off-by-one vulnerability in the fgetwln function in libbsd before 0.8.2 allows attackers to have unspecified impact via unknown vectors, which trigger a heap-based buffer overflow.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-2090",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.5,
                "Vectors": "AV:N/AC:L/Au:N/C:P/I:P"
              }
            }
          }
        }
      ],
      "Image": "postgres:latest",
      "ScanStarted": false
    },
    {
      "LastCheck": "2017-01-26T19:48:08.535599784Z",
      "Vulnerabilities": [
        {
          "Name": "CVE-2017-5225",
          "NamespaceName": "debian:8",
          "Description": "LibTIFF version 4.0.7 is vulnerable to a heap buffer overflow in the tools/tiffcp resulting in DoS or code execution via a crafted BitsPerSample value.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2017-5225",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.5,
                "Vectors": "AV:N/AC:L/Au:N/C:P/I:P"
              }
            }
          }
        },
        {
          "Name": "CVE-2016-9535",
          "NamespaceName": "debian:8",
          "Description": "tif_predict.h and tif_predict.c in libtiff 4.0.6 have assertions that can lead to assertion failures in debug mode, or buffer overflows in release mode, when dealing with unusual tile size like YCbCr with subsampling. Reported as MSVR 35105, aka \"Predictor heap-buffer-overflow.\"",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2016-9535",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.5,
                "Vectors": "AV:N/AC:L/Au:N/C:P/I:P"
              }
            }
          }
        },
        {
          "Name": "CVE-2015-1258",
          "NamespaceName": "debian:8",
          "Description": "Google Chrome before 43.0.2357.65 relies on libvpx code that was not built with an appropriate --size-limit value, which allows remote attackers to trigger a negative value for a size field, and consequently cause a denial of service or possibly have unspecified other impact, via a crafted frame size in VP9 video data.",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2015-1258",
          "Severity": "High",
          "Metadata": {
            "NVD": {
              "CVSSv2": {
                "Score": 7.5,
                "Vectors": "AV:N/AC:L/Au:N/C:P/I:P"
              }
            }
          }
        }
      ],
      "Image": "nginx",
      "ScanStarted": false
    }
  ]
}
```
