# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import kawesome.v1.kawesome

kawesome: ac.AppConfiguration {
    # Declare the workload configurations. 
    workload: service.Service {
        containers: {
            kawesome: c.Container {
                image: "hashicorp/http-echo"
                env: {
                    "ECHO_TEXT": "$(KUSION_KAWESOME_RANDOM_PASSWORD)"
                }
            }
        }
        replicas: 1
    }
    # Declare the kawesome module configurations. 
    accessories: {
        "kawesome": kawesome.Kawesome {
            service: kawesome.Service{
                port: 5678
            }
            randomPassword: kawesome.RandomPassword {
                length: 20
            }
        }
    }
}
