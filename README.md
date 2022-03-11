# LiteKube
Aims to build a container deployment monitoring system for edge weak configuration scenarios, and stay the same call-api with K8S.

# Develop Enviroment

|  core     | version |
|   :-      |   :-    |
| Arch      |  x86_64 |
| Ubuntu    |  18.04  |
| golang    | v1.17.5
| kubelet   | v1.23.1 |

# Build Develop-Environment

## 1. prepare environment to run kubelet

See this [blog](https://www.aflyingfish.top/articles/205801b55ca4/) post for more details.  *Notice: it's not need to install any K8S-related components!*

* stop firewalld
    ```bash
    systemctl stop firewalld
    systemctl disable firewalld
    ```
* close free-swap
    ```bash
    sed -i 's/enforcing/disabled/' /etc/selinux/config
    ```
* install docker
* set `--exec-opt native.cgroupdriver=systemd` for docker
    ```bash
    vim /usr/lib/systemd/system/docker.service

    -------------- detail-----------
    [Unit]
    Description=Docker Application Container Engine
    # .....

    [Service]
    Type=notify
    # the default is not to use systemd for cgroups because the delegate issues still
    # exists and systemd currently does not support the cgroup feature set required
    # for containers run by docker
    ExecStart=/usr/bin/dockerd  --exec-opt native.cgroupdriver=systemd -H fd:// --containerd=/run/containerd/containerd.sock
    ExecReload=/bin/kill -s HUP $MAINPID
    
    # .....
    ```
* reboot OS

## 2. Install Components

```bash
mkdir -p ~/litekube/downloads && cd ~/litekube/downloads
wget https://dl.k8s.io/v1.23.1/kubernetes-node-linux-amd64.tar.gz && tar -xavf kubernetes-node-linux-amd64.tar.gz && chmod +x kubernetes/node/bin/*

cd ~/litekube/ 
cp downloads/kubernetes/node/bin/kubelet ./

# kube-proxy and kubectl-convert is not used at this time.
cp downloads/kubernetes/node/bin/kube-proxy ./
cp downloads/kubernetes/node/bin/kubectl-convert ./

# install debug-tools
mkdir -p ~/litekube/downloads && cd ~/litekube/downloads
wget https://github.com/cyberark/kubeletctl/releases/download/v1.8/kubeletctl_linux_amd64 && chmod +x kubeletctl_linux_amd64

cd ~/litekube/ 
cp downloads/kubeletctl_linux_amd64 ./kubeletctl

# add soft-link to /usr/bin/
sudo ln -s $HOME/litekube/kube* /usr/bin/
```

## 3. Prepare to run

```bash
cd ~/litekube/ && mkdir -p config/kubelet/ssl/ca && cd config/kubelet/ssl/ca/

# generate for CA
cat > ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "87600h"
        },
        "profiles": {
            "kubernetes": {
                "expiry": "87600h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth",
                    "client auth"
                ]
            }
        }
    }
}
EOF

cat > ca-csr.json <<EOF
{
    "CN": "kubernetes",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "Beijing",
            "ST": "Beijing",
            "O": "k8s",
            "OU": "System"
        }
    ]
}
EOF

# need to install golang-cfssl.
cfssl gencert -initca ca-csr.json | cfssljson -bare ca - 

# for client to visit node:10250
cd ~/litekube/ && mkdir -p config/kubelet/ssl/client && cd config/kubelet/ssl/client/
cat > client-csr.json <<EOF
{
    "CN": "kubernetes",
    "hosts": [
        "127.0.0.1",
        "kubernetes",
        "kubernetes.default",
        "kubernetes.default.svc",
        "kubernetes.default.svc.cluster",
        "kubernetes.default.svc.cluster.local"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BeiJing",
            "ST": "BeiJing",
            "O": "k8s",
            "OU": "System"
        }
    ]
}
EOF

cfssl gencert -ca=../ca/ca.pem -ca-key=../ca/ca-key.pem -config=../ca/ca-config.json -profile=kubernetes client-csr.json | cfssljson -bare client

```

## 4. Run Kubelet with Standalone-Mode

* prepare enviroment
```bash
cd ~/litekube/
mkdir -p manifests && mkdir -p config/kubelet/pki && cd config/kubelet/

cat > kubelet <<EOF
export KUBELET_OPT=" \
--v=4 \
--hostname-override=127.0.0.1 \
--config=/home/aflyingfish/litekube/config/kubelet/kubelet.yaml \
--cert-dir=/home/aflyingfish/litekube/config/kubelet/pki/ssl \
--cgroup-driver=systemd \
--file-check-frequency=20s \
--max-pods=30 \
--pod-manifest-path=/home/aflyingfish/litekube/manifests \
--pod-infra-container-image=registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0"

export CERT="--cacert /home/aflyingfish/litekube/config/kubelet/ssl/ca/ca.pem --cert /home/aflyingfish/litekube/config/kubelet/ssl/client/client.pem --key /home/aflyingfish/litekube/config/kubelet/ssl/client/client-key.pem"
EOF

cat > kubelet.yaml <<EOF
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: true
  webhook:
    cacheTTL: 0s
    enabled: false
  x509:
    clientCAFile: /home/aflyingfish/litekube/config/kubelet/ssl/ca/ca.pem
authorization:
  mode: AlwaysAllow
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s
clusterDNS:
- 10.96.0.10
clusterDomain: cluster.local
cpuManagerReconcilePeriod: 0s
evictionPressureTransitionPeriod: 0s
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 0s
imageMinimumGCAge: 0s
kind: KubeletConfiguration
nodeStatusReportFrequency: 0s
nodeStatusUpdateFrequency: 0s
rotateCertificates: true
runtimeRequestTimeout: 0s
staticPodPath: /home/aflyingfish/litekube/manifests
streamingConnectionIdleTimeout: 0s
syncFrequency: 0s
volumeStatsAggPeriod: 0s
file-check-frequency: 20s
EOF
```

* run kubelet with Standalone-mode
```bash
# run only once at one login-in
cd ~/litekube/config/kubelet && source ./kubelet

# start to run
sudo kubelet $KUBELET_OPT
```

* run one pod (static-pod)
```bash
cat > ~/litekube/manifests <<EOF
apiVersion: v1
kind: Pod
metadata:
 name: static-web
 labels:
  name: static-web
spec:
 containers:
 - name: static-web
   image: nginx
   ports:
   - name: web
     containerPort: 80
     hostPort: 80
EOF
```
now, you can view `nginx hello` by `curl http://127.0.0.1`

## How to debug directly

note: if authentication.anonymous.enabled=true in kubelet.yaml, certificate can be ignored.


* By `curl`
```bash
curl -k $CERT https://127.0.0.1:10250/pods      # -k means uncheck certificate of kubelet itself
```

* By `kubeletctl`

```bash
cd ~/litekube/config/kubelet/ssl/
kubeletctl pods -s 127.0.0.1 $CERT      
```

## How to use LiteKube
[run lite-apiserver](docs/READNE.md)

## fix port already in use
```shell
ps -ef | grep kubelet
sudo kill -9 $PID

sudo netstat -an -p | grep $PID # view port used by this process
```

curl -k $CERT https://127.0.0.1:10250/exec/default/{podID}/{containerName}?command={command}/&input=1&output=1&tty=1"