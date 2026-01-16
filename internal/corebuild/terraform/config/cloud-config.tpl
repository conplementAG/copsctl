#cloud-config

bootcmd:
  - mkdir -p /etc/systemd/system/walinuxagent.service.d
  - echo "[Unit]\nAfter=cloud-final.service" > /etc/systemd/system/walinuxagent.service.d/override.conf
  - sed "s/After=multi-user.target//g" /lib/systemd/system/cloud-final.service > /etc/systemd/system/cloud-final.service
  - systemctl daemon-reload

apt:
  sources:
    docker.list:
      source: deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable
      keyid: 9DC858229FC7DD38854AE2D88D81803C0EBFCD88

packages:
  - docker-ce
  - docker-ce-cli

groups:
  - docker

disk_setup:
  ephemeral0:
    table_type: gpt
    layout: [66, [33,82]]
    overwrite: true
%{ if use_data_disk ~}
  /dev/sdc:
    table_type: gpt
    layout: true
    overwrite: true
%{ endif ~}

fs_setup:
  - device: ephemeral0.1
    filesystem: ext4
%{ if use_data_disk ~}
  - device: /dev/sdc
    filesystem: ext4
%{ endif ~}

mounts:
%{ if use_data_disk ~}
  - ["/dev/sdc", "/agent"]
%{ else ~}
  - ["ephemeral0.1", "/agent"]
%{ endif ~}

