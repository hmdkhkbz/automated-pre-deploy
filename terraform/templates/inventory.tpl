all:
  vars:
    ansible_ssh_port: 22
    ansible_ssh_private_key_file: /root/ary.pem
    populate_inventory_to_hosts_file: true
  children:
    all:
      children:
        control:
          hosts:
%{ for host in controllers ~}
            ${host.name}:
              ansible_host: ${host.ansible_host}
              access_ip: ${host.access_ip}
              ansible_hostname: ${host.name}
%{ endfor ~}
        compute:
          hosts:
%{ for host in computes ~}
            ${host.name}:
              ansible_host: ${host.ansible_host}
              access_ip: ${host.access_ip}
              ansible_hostname: ${host.name}
%{ endfor ~}
        network:
          hosts:
%{ for host in networks ~}
            ${host.name}:
              ansible_host: ${host.ansible_host}
              access_ip: ${host.access_ip}
              ansible_hostname: ${host.name}
%{ endfor ~}
