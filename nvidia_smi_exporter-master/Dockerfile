FROM nvidia/cuda
MAINTAINER Kevin Jen <kevin7674@gmail.com>

COPY nvidia_smi_exporter.go /nvidia_smi_exporter.go
RUN apt-get update && apt-get install -y golang-go && go build /nvidia_smi_exporter.go && apt-get remove -y golang-go

EXPOSE 9101:9101

ENTRYPOINT ["/nvidia_smi_exporter"]
