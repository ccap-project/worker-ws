roles-ws
=========

API that interfaces with gitlab, control and serve roles (ansible) informations



Endpoints
----------------
/api/v1/roles
/api/v1/roles/:id/meta/:version
/api/v1/roles/:id/params/:version



Endpoint Response Examples
----------------

$ curl http://127.0.0.1:8000/api/v1/roles
```
[
  {
    "ID": 36,
    "Name": "vmware-collector-ansible",
    "Url": "https://gitlab.uoldiveo.intranet/ansible/vmware-collector-ansible.git",
    "Versions": null
  },
  {
    "ID": 31,
    "Name": "checklist-users",
    "Url": "https://gitlab.uoldiveo.intranet/ansible/checklist-users.git",
    "Versions": null
  },
  {
    "ID": 24,
    "Name": "syslog-ng",
    "Url": "https://gitlab.uoldiveo.intranet/ansible/syslog-ng.git",
    "Versions": null
  },
  {
    "ID": 21,
    "Name": "rsyslog",
    "Url": "https://gitlab.uoldiveo.intranet/ansible/rsyslog.git",
    "Versions": [
      "v2.0",
      "v1.0"
    ]
  },
  {
    "ID": 16,
    "Name": "diamond",
    "Url": "https://gitlab.uoldiveo.intranet/ansible/diamond.git",
    "Versions": [
      "1.00"
    ]
  }
]
```

$ curl http://127.0.0.1:8000/api/v1/roles/16/meta/1.00 | jq

```{
  "author": "EngCloud",
  "company": "UOLDIVEO",
  "dependencies": [],
  "description": "install Diamond and configure its collectors",
  "galaxy_tags": [
    "diamond"
  ],
  "issue_tracker_url": "https://gitlab.uoldiveo.intranet/ansible/diamond/issues",
  "license": "BSD",
  "min_ansible_version": 1.2
}
```

$ curl http://127.0.0.1:8000/api/v1/roles/16/params/1.00
```
{
  "diamond_collector_CPU_byte_unit": "byte",
  "diamond_collector_CPU_measure_collector_time": false,
  "diamond_collector_CPU_normalize": true,
  "diamond_collector_CPU_path_suffix": "",
  "diamond_collector_CPU_percore": false,
  "diamond_collector_CPU_simple": false,
  "diamond_collector_CPU_ttl_multiplier": 2,
  "diamond_collector_CephStats_byte_unit": "byte",
  "diamond_collector_CephStats_ceph_binary": "/usr/bin/ceph",
  "diamond_collector_CephStats_enabled": true,
  "diamond_collector_CephStats_measure_collector_time": false,
  "diamond_collector_CephStats_metrics_whitelist": [],
  "diamond_collector_CephStats_path": "",
  "diamond_collector_CephStats_path_suffix": "",
  "diamond_collector_CephStats_socket_ext": "asok",
  "diamond_collector_CephStats_socket_path": "/var/run/ceph",
  "diamond_collector_CephStats_socket_prefix": "ceph-",
  "diamond_collector_CephStats_ttl_multiplier": 2,
  "diamond_collector_Ceph_byte_unit": "byte",
  "diamond_collector_Ceph_ceph_binary": "/usr/bin/ceph",
  "diamond_collector_Ceph_enabled": true,
  "diamond_collector_Ceph_measure_collector_time": false,
  "diamond_collector_Ceph_metrics_whitelist": [
    "filestore.journal_latency",
    "filestore.journal_queue_bytes",
    "filestore.journal_queue_ops",
    "osd.op_latency",
    "osd.op_r",
    "osd.op_w",
    "osd.recovery_ops"
  ],
  "diamond_collector_Ceph_path": "",
  "diamond_collector_Ceph_path_suffix": "",
  "diamond_collector_Ceph_socket_ext": "asok",
  "diamond_collector_Ceph_socket_path": "/var/run/ceph",
  "diamond_collector_Ceph_socket_prefix": "ceph-",
  "diamond_collector_Ceph_ttl_multiplier": 2,
  "diamond_collector_DiskSpace_byte_unit": "byte",
  "diamond_collector_DiskSpace_exclude_filters": [
    "^/export/home"
  ],
  "diamond_collector_DiskSpace_filesystems": [
    "xfs",
    "ext2",
    "ext3",
    "ext4",
    "nfs"
  ],
  "diamond_collector_DiskSpace_measure_collector_time": false,
  "diamond_collector_DiskSpace_path_suffix": "",
  "diamond_collector_DiskSpace_ttl_multiplier": 2,
  "diamond_collector_DiskUsage_attach_device_label": false,
  "diamond_collector_DiskUsage_byte_unit": "byte",
  "diamond_collector_DiskUsage_devices": "PhysicalDrive[0-9]+$|md[0-9]+$|sd[a-z]+[0-9]*$|x?vd[a-z]+[0-9]*$|disk[0-9]+$",
  "diamond_collector_DiskUsage_measure_collector_time": false,
  "diamond_collector_DiskUsage_metrics_sector_size": 512,
  "diamond_collector_DiskUsage_metrics_send_zero": true,
  "diamond_collector_DiskUsage_metrics_whitelist": [
    "average_queue_length",
    "average_request_size_byte",
    "await",
    "iops",
    "read_kilobyte_per_second",
    "reads_per_second",
    "service_time",
    "util_percentage",
    "write_kilobyte_per_second",
    "writes_per_second"
  ],
  "diamond_collector_DiskUsage_path_suffix": "",
  "diamond_collector_DiskUsage_ttl_multiplier": 2,
  "diamond_collector_HAProxy_byte_unit": "byte",
  "diamond_collector_HAProxy_enabled": true,
  "diamond_collector_HAProxy_ignore_servers": false,
  "diamond_collector_HAProxy_measure_collector_time": false,
  "diamond_collector_HAProxy_method": "http",
  "diamond_collector_HAProxy_metrics_blacklist": [],
  "diamond_collector_HAProxy_metrics_whitelist": [],
  "diamond_collector_HAProxy_pass": "password",
  "diamond_collector_HAProxy_sock": "/var/run/haproxy.sock",
  "diamond_collector_HAProxy_url": "http://localhost/haproxy?stats;csv",
  "diamond_collector_HAProxy_user": "admin",
  "diamond_collector_Httpd_byte_unit": "byte",
  "diamond_collector_Httpd_enabled": true,
  "diamond_collector_Httpd_measure_collector_time": false,
  "diamond_collector_Httpd_metrics_blacklist": [],
  "diamond_collector_Httpd_metrics_whitelist": [],
  "diamond_collector_Httpd_urls": [],
  "diamond_collector_KSM_byte_unit": "byte",
  "diamond_collector_KSM_enabled": true,
  "diamond_collector_KSM_measure_collector_time": false,
  "diamond_collector_KSM_metrics_whitelist": [],
  "diamond_collector_LibvirtKVM_byte_unit": "byte",
  "diamond_collector_LibvirtKVM_count_provisioned_vcpus": true,
  "diamond_collector_LibvirtKVM_count_vms_by_state": true,
  "diamond_collector_LibvirtKVM_enabled": true,
  "diamond_collector_LibvirtKVM_format_name_using_metadata": "${owner_project}_${owner_project_uuid}.${instance}_${instance_uuid}",
  "diamond_collector_LibvirtKVM_measure_collector_time": true,
  "diamond_collector_LibvirtKVM_path": ".",
  "diamond_collector_LibvirtKVM_ttl_multiplier": 2,
  "diamond_collector_LoadAverage_byte_unit": "byte",
  "diamond_collector_LoadAverage_measure_collector_time": false,
  "diamond_collector_LoadAverage_path_suffix": "",
  "diamond_collector_LoadAverage_simple": false,
  "diamond_collector_LoadAverage_ttl_multiplier": 2,
  "diamond_collector_Memcached_byte_unit": "byte",
  "diamond_collector_Memcached_enabled": true,
  "diamond_collector_Memcached_hosts": "localhost:11211",
  "diamond_collector_Memcached_measure_collector_time": false,
  "diamond_collector_Memcached_metrics_blacklist": [],
  "diamond_collector_Memcached_metrics_whitelist": [],
  "diamond_collector_Memcached_publish": "",
  "diamond_collector_Memory_byte_unit": "byte",
  "diamond_collector_Memory_detailed": "",
  "diamond_collector_Memory_measure_collector_time": false,
  "diamond_collector_Memory_metrics_whitelist": [
    "Active",
    "Buffers",
    "Cached",
    "MemAvailable",
    "MemTotal",
    "MemFree",
    "Committed_AS"
  ],
  "diamond_collector_Memory_path_suffix": "",
  "diamond_collector_Memory_ttl_multiplier": 2,
  "diamond_collector_Network_byte_unit": "byte",
  "diamond_collector_Network_greedy": true,
  "diamond_collector_Network_interfaces": [
    "eth",
    "bond",
    "em",
    "p1p",
    "eno",
    "enp",
    "ens",
    "enx"
  ],
  "diamond_collector_Network_measure_collector_time": false,
  "diamond_collector_Network_metrics_whitelist": [
    "rx_byte",
    "rx_drop",
    "rx_errors",
    "rx_packets",
    "tx_byte",
    "tx_drop",
    "tx_errors",
    "tx_packets"
  ],
  "diamond_collector_Network_path_suffix": "",
  "diamond_collector_Network_ttl_multiplier": 2,
  "diamond_collector_RabbitMQ_byte_unit": "byte",
  "diamond_collector_RabbitMQ_cluster": false,
  "diamond_collector_RabbitMQ_enabled": true,
  "diamond_collector_RabbitMQ_host": "localhost:5672",
  "diamond_collector_RabbitMQ_measure_collector_time": false,
  "diamond_collector_RabbitMQ_metrics_blacklist": [],
  "diamond_collector_RabbitMQ_metrics_whitelist": [],
  "diamond_collector_RabbitMQ_password": "guest",
  "diamond_collector_RabbitMQ_queues": "",
  "diamond_collector_RabbitMQ_queues_ignored": "",
  "diamond_collector_RabbitMQ_replace_dot": false,
  "diamond_collector_RabbitMQ_replace_slash": false,
  "diamond_collector_RabbitMQ_user": "guest",
  "diamond_collector_RabbitMQ_vhosts": "",
  "diamond_collector_Redis_auth": "",
  "diamond_collector_Redis_byte_unit": "byte",
  "diamond_collector_Redis_databases": 16,
  "diamond_collector_Redis_db": 0,
  "diamond_collector_Redis_enabled": true,
  "diamond_collector_Redis_host": "localhost",
  "diamond_collector_Redis_instances": [],
  "diamond_collector_Redis_measure_collector_time": false,
  "diamond_collector_Redis_metrics_blacklist": [],
  "diamond_collector_Redis_metrics_whitelist": [],
  "diamond_collector_Redis_port": 6379,
  "diamond_collector_Redis_timeout": 5,
  "diamond_collector_VMStat_byte_unit": "byte",
  "diamond_collector_VMStat_measure_collector_time": false,
  "diamond_collector_VMStat_path_suffix": "",
  "diamond_collector_VMStat_ttl_multiplier": 2,
  "diamond_collectors_conf_dir": "{{ diamond_conf_dir }}/collectors",
  "diamond_collectors_path": "/usr/share/diamond/collectors/",
  "diamond_conf_dir": "/etc/diamond",
  "diamond_custom_collectors": [],
  "diamond_default_collectors": [
    "CPU",
    "DiskSpace",
    "DiskUsage",
    "LoadAverage",
    "Network",
    "Memory",
    "VMStat"
  ],
  "diamond_default_interval": 20,
  "diamond_graphite_batch_size": 1,
  "diamond_graphite_host": "localhost",
  "diamond_graphite_pickle_batch_size": 512,
  "diamond_graphite_pickle_host": "localhost",
  "diamond_graphite_pickle_port": 2004,
  "diamond_graphite_pickle_timeout": 15,
  "diamond_graphite_port": 2003,
  "diamond_graphite_timeout": 15,
  "diamond_instance_prefix": "",
  "diamond_manage_service": true,
  "diamond_path_prefix": "servers",
  "diamond_path_suffix": "",
  "diamond_pid_dir": "/var/run"
}
```
