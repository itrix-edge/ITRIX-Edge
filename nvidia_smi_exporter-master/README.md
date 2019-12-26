# nvidia_smi_exporter

nvidia-smi metrics exporter for Prometheus

### 1.下載並執行nvidia_smi_exporter
```
git clone https://github.com/kevin7674/nvidia_smi_exporter.git
cd nvidia_smi_exporter
./nvidia_smi_exporter 9101 &
```
### 2.修改 vi /etc/rc.local
```
./nvidia_smi_exporter 9101 &
```




# Build Image
```
> docker build -t="nvidia_smi_exporter:0" .
```

### Docker Run
```
> nvidia-docker run -d --net="host" nvidia_smi_exporter:0 --restart=always
```
Default port is 9101

### curl localhost:9101/metrics
```
temperature_gpu{gpu="Tesla V100-SXM2-16GB[0]"} 34
utilization_gpu{gpu="Tesla V100-SXM2-16GB[0]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[0]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[0]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[0]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[0]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[1]"} 37
utilization_gpu{gpu="Tesla V100-SXM2-16GB[1]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[1]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[1]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[1]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[1]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[2]"} 36
utilization_gpu{gpu="Tesla V100-SXM2-16GB[2]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[2]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[2]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[2]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[2]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[3]"} 33
utilization_gpu{gpu="Tesla V100-SXM2-16GB[3]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[3]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[3]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[3]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[3]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[4]"} 36
utilization_gpu{gpu="Tesla V100-SXM2-16GB[4]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[4]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[4]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[4]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[4]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[5]"} 37
utilization_gpu{gpu="Tesla V100-SXM2-16GB[5]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[5]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[5]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[5]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[5]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[6]"} 39
utilization_gpu{gpu="Tesla V100-SXM2-16GB[6]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[6]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[6]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[6]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[6]"} 10
temperature_gpu{gpu="Tesla V100-SXM2-16GB[7]"} 38
utilization_gpu{gpu="Tesla V100-SXM2-16GB[7]"} 0
utilization_memory{gpu="Tesla V100-SXM2-16GB[7]"} 0
memory_total{gpu="Tesla V100-SXM2-16GB[7]"} 16152
memory_free{gpu="Tesla V100-SXM2-16GB[7]"} 16142
memory_used{gpu="Tesla V100-SXM2-16GB[7]"} 10
```

### Exact command
```
nvidia-smi --query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used --format=csv,noheader,nounits
```

### Prometheus example config

```
- job_name: "gpu_exporter"
  static_configs:
  - targets: ['localhost:9101']
```

