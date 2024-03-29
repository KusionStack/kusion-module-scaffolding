# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import kam.v1.workload as wl
import kam.v1.workload.container as c
import kawesome.v1.kawesome as ks

kawesome: ac.AppConfiguration {
    # Declare the workload configurations. 
    workload: wl.Service {
        containers: {
            kawesome: c.Container {
                image: "hashicorp/http-echo"
                env: {
                    "ECHO_TEXT": "$(KUSION_KAWESOME_RANDOM_PASSWORD)"
                }
            }
        }
    }
    # Declare the kawesome module configurations. 
    accessories: {
        "kawesome": ks.Kawesome {
            service: ks.Service{
                port: 5678
            }
            randomPassword: ks.RandomPassword {
                length: 20
            }
        }
    }
}
