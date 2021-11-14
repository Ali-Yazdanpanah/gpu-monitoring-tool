import pynvml

try:
    pynvml.nvmlInit()
except pynvml.NVMLError as error:
    print(error) 
    # Driver Not Loaded      驱动加载失败(没装驱动或者驱动有问题)
    # Insufficent Permission 没有以管理员权限运行  pynvml.NVMLError_DriverNotLoaded: Driver Not Loaded
    exit()

try:
    print(pynvml.nvmlDeviceGetCount())
except pynvml.NVMLError as error:
    print(error)

print(pynvml.nvmlDeviceGetCount())# total gpu count = 1
print(pynvml.nvmlSystemGetDriverVersion()) # 396.54

GPU_ID = 0
handle = pynvml.nvmlDeviceGetHandleByIndex(GPU_ID)
print(pynvml.nvmlDeviceGetName(handle)) # GeForce GTX 1060

meminfo = pynvml.nvmlDeviceGetMemoryInfo(handle)
use = pynvml.nvmlDeviceGetUtilizationRates(handle)
ratio = 0.5*(float(use.gpu+float(use.memory)))
MB_SIZE = 1024*1024
print(meminfo.total/MB_SIZE) # 6078 MB
print(meminfo.used/MB_SIZE)  # 531 MB
print(meminfo.free/MB_SIZE)  # 5546 MB
print(ratio)

pynvml.nvmlShutdown()

