# Blackbox arp exporter 
The blackbox exporter allows blackbox probing of endpoints over
ARP

### Build:
```
git clone https://git.digikala.com/dkdevops/exporters/blackbox-arp-exporter.git
cd PROJECT_DIR
go build -o blackbox-arp-exporter
./blackbox-arp-exporter -help
```

### How to run?

Then:

    ./blackbox-arp-exporter -interface {INTERFACE_TO_SEND_ARP_REQUESTS} -listen-address {IP:PORT}

Example:

    ./blackbox-arp-exporter -interface eth0 -listen-address :9185

### Using the docker image

    docker run --rm \
      --network host \
      --cap-add CAP_NET_RAW
      --name blackbox_arp_exporter \
      local_image
      
### Using the docker compose:

    docker compose up -d
