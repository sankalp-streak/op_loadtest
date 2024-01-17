#/bin/bash

while sudo fuser /var/lib/dpkg/lock /var/lib/apt/lists/lock /var/cache/apt/archives/lock >/dev/null 2>&1; do echo 'Waiting for release of dpkg/apt locks'; sleep 5; done;
apt update
while sudo fuser /var/lib/dpkg/lock /var/lib/apt/lists/lock /var/cache/apt/archives/lock >/dev/null 2>&1; do echo 'Waiting for release of dpkg/apt locks'; sleep 5; done;
apt -y install supervisor

dpkg --add-architecture amd64
dpkg --add-architecture i386
while sudo fuser /var/lib/dpkg/lock /var/lib/apt/lists/lock /var/cache/apt/archives/lock >/dev/null 2>&1; do echo 'Waiting for release of dpkg/apt locks'; sleep 5; done;
apt update
rm -rf /home/ubuntu/loadTest/
dpkg --unpack /home/ubuntu/loadTestSlave.deb
chown -R ubuntu:ubuntu /home/ubuntu/loadTest/
ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
rm /home/ubuntu/loadTest/config.json
mkdir /home/ubuntu/loadTest/slave
SLAVE_IP=$(curl https://ipinfo.io/ip)
{
    echo "{\"master_path\": \"rch.streak.tech:6379\", \"master_redis_auth\": \"b626aaa7a231aabf6aa3df5fc5caa847c202dc3bd0d361a1291bb2f855afe8ba7727d3e4273e9d2a4ad5020c8695bbaf9b8d051c5bf43b868795a130d482acae\", \"slave_port\": \"9093\", \"slave_ip\": \""$SLAVE_IP"\", \"slave_ip_type\": \"private\"}"
} >> /home/ubuntu/loadTest/config.json

{
    echo "[program:slave]"
    echo "directory=/home/ubuntu/loadTest"
    echo "command=/home/ubuntu/loadTest/slave"
    echo "autostart=true"
    echo "autorestart=true"
    echo "user=ubuntu"
    echo "stderr_logfile=/var/log/slave.log"
    echo "stdout_logfile=/var/log/slave.log"
} >> /home/ubuntu/loadTest/slave.conf

cp /home/ubuntu/loadTest/slave.conf /etc/supervisor/conf.d/slave.conf

supervisorctl reread
supervisorctl update
supervisorctl restart slave