schema Kawesome: 
    """ Kawesome is a sample module schema consisting of Service
    and RandomPassword

    Attributes
    ----------
    service: Service, default is Undefined, required. 
        The exposed port of Workload, which will be generated into Kubernetes Service. 
    randomPassword: RandomPassword, default is Undefined, required. 
        The sensitive random string, which will be generated into Terraform random_password. 
    
    Examples
    --------
    import kawesome.v1.kawesome as ks

    accessories: {
        "kusionstack/kawesome@v0.1.0": ks.Kawesome {
            service: ks.Service {
                port: 8080
            }
            randomPassword: ks.RandomPassword {
                length: 20
            }
        }
    }
    """

    # The exposed port of Workload, which will be generated into Kubernetes Service. 
    service:                    Service

    # The sensitive random string, which will be generated into Terraform random_password. 
    randomPassword:             RandomPassword

schema Service: 
    """ Service defines the exposed port of Workload, and will be generated 
    into a Kubernetes Service resource. 

    Attributes
    ----------
    port: int, default is 80, required. 
        The exposed port of the Workload. 
    targetPort: int, default is Undefined, optional. 
        The backend container port. If empty, it will be set the same as the port. 
    protocol: "TCP" | "UDP", default is "TCP", required. 
        The protocol to access the port. 

    Examples
    --------
    import kawesome_schema as ks

    svc = ks.Service {
        port: 80
        targetPort: 8080
        protocol: "TCP"
    }
    """

    # The exposed port of the Service. 
    port:                       int = 80

    # The backend container port. 
    targetPort?:                int

    # The protocol of port. 
    protocol:                   "TCP" | "UDP" = "TCP"

    check: 
        1 <= port <= 65535, "port must be between 1 and 65535, inclusive"
        1 <= targetPort <= 65535 if targetPort, "targetPort must be between 1 and 65535, inclusive"

schema RandomPassword: 
    """ RandomPassword defines a sensitive random string with specified length, 
    and will be generated into a Terraform random_password resource. 

    Attributes
    ----------
    length: int, default is 16, required. 
        The specified length of the random password string. 
    
    Examples
    --------
    import kawesome_schema as ks

    rp = ks.RandomPassword {
        length: 20
    }
    """

    # The specified length of the random password string. 
    length:                     int = 16

    check: 
        length >= 1, "length must be more than 0"
