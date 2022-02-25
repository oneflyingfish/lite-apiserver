# Prepare the enviroment

## 0. At the end of each of the following steps, the default is to return to the original project root directory.

## 1. Create CA Certificate For lite-apiserver

    ```bash
    # move lite-apiserver to here
    mkdir -p test/ssl/ca && cd test/ssl/ca/

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
        "CN": "lite-apiserver",
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
    ```

## 2. Create CA Certificate config YAML

    ```bash
    cd test/

    cat > server-ca-cert.yaml <<EOF
    cacert: test/ssl/ca/ca.pem
    cakey: test/ssl/ca/ca-key.pem
    EOF
    ```

## 3. Create config for lite-apiserver

    ```bash
    cd test/

    cat > server.yaml <<EOF
    hostname: $NODE_IP
    port: 20500
    insecure-port: 20501
    ca-tls-configpath: test/server-ca-cert.yaml 
    EOF
    ```

## 4. Create client-TLS Certificate for connect to kubelet

    client-TLS Certificate is issue by another CA-Certificate, which is actually a similar certificate created in [this page](../README.md) for debug.
    
    ```bash
    cd test/

    cat > kubelet-cert.yaml <<EOF
    cert: test/ssl/client.pem
    key: test/ssl/client-key.pem
    EOF
    ```

## 5. Create config for connect to kubelet

    ```bash
    cat > kubelet.yaml <<EOF
    kubelet-hostname: 127.0.0.1
    kubelet-healthzport: 10248
    kubelet-port: 10250
    kubelet-client-cert-config: test/kubelet-cert.yaml
    EOF
    ```

# Run lite-apiserver

* Run by configs modified above
    ```bash
    $ ./lite-apiserver --kubelet-config=./test/kubelet.yaml --config=server.yaml
    ```

* Also, you can give args by command flags. View how to run:
    ```
    $ ./lite-apiserver --help

    '''
    The lite-apiserver is one simplified version of kube-apiserver, which is only service for one node and deal with pods.

    Usage:
        lite-apiserver [flags]

    Lite-apiserver flags:

        --ca-tls-configpath string                                                                                                                                                                                                                           
                    path to config store the X.509 Certificate information for lite-apiserver (default: "")
        --config string                                                                                                                                                                                                                                      
                    config for lite-apiserver (lower priority to flags)
        --hostname string                                                                                                                                                                                                                                    
                    hostname of lite-apiserver (default: 127.0.0.1)
        --insecure-port int                                                                                                                                                                                                                                  
                    http port of lite-apiserver, not secure, set 0 to disable (default: 0)
        --port int                                                                                                                                                                                                                                           
                    https port of lite-apiserver (default: 6500)

    Kubelet flags:

        --kubelet-client-cert-config string                                                                                                                                                                                                                  
                    path to config store the X.509 Certificate information to kubelet (default: "")
        --kubelet-config string                                                                                                                                                                                                                              
                    config for kubelet (lower priority to flags)
        --kubelet-healthzport int                                                                                                                                                                                                                            
                    healthz port of kubelet (default: 10248)
        --kubelet-hostname string                                                                                                                                                                                                                            
                    hostname of kubelet (default: 127.0.0.1)
        --kubelet-port int                                                                                                                                                                                                                                   
                    port of kubelet (default: 10250)
    '''

    # for example:
    $ ./lite-apiserver \
     --hostname=127.0.0.1 \
     --port=20500 \
     --insecure-port=20501 \
     --ca-tls-configpath=/home/aflyingfish/LiteKube/cmd/lite-apiserver/test/server-ca-cert.yaml \
     --kubelet-hostname=localhost \
     --kubelet-healthzport=10248 \
     --kubelet-port=10250 \
     --kubelet-client-cert-config=/home/aflyingfish/LiteKube/cmd/lite-apiserver/test/kubelet-cert.yaml
    ```

* Surely, you can use this two way tother, and flags will be in a higher priority.

    ```bash
    $ ./lite-apiserver --hostname=localhost --kubelet-config=./test/kubelet.yaml --config=server.yaml
    ```