#/bin/bash

pkill -9 -f "unattended-upgrade"
service unattended-upgrades stop
apt-get purge -y unattended-upgrades
while sudo fuser /var/lib/dpkg/lock /var/lib/apt/lists/lock /var/cache/apt/archives/lock >/dev/null 2>&1; do echo 'Waiting for release of dpkg/apt locks'; sleep 5; done;
pkill -9 -f "unattended-upgrade"
apt update
pkill -9 -f "unattended-upgrade"
apt -y install supervisor

mkdir /home/ubuntu/overwatch_slave 
chown -R ubuntu:ubuntu overwatch_slave/
ln -sf /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
cp /home/ubuntu/slave /home/ubuntu/overwatch_slave/slave 

{
	echo "[unix_http_server]"
	echo "file=/var/run/supervisor.sock   ; (the path to the socket file)"
	echo "chmod=0700                       ; sockef file mode (default 0700)"
	echo ""
	echo "[supervisord]"
	echo "logfile=/var/log/supervisor/supervisord.log"
	echo "pidfile=/var/run/supervisord.pid"
	echo "minfds=65535"
	echo "childlogdir=/var/log/supervisor"
	echo ""
	echo "[rpcinterface:supervisor]"
	echo "supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface"
	echo ""
	echo "[supervisorctl]"
	echo "serverurl=unix:///var/run/supervisor.sock"
	echo ""
	echo "[include]"
	echo "files = /etc/supervisor/conf.d/*.conf"
} >> /home/ubuntu/overwatch_slave/supervisord.conf

cp /home/ubuntu/overwatch_slave/supervisord.conf /etc/supervisor/supervisord.conf
service supervisor restart

{
    echo "[program:overwatch_slave]"
    echo "directory=/home/ubuntu/overwatch_slave"
    echo "command=/home/ubuntu/overwatch_slave/slave"
    echo "autostart=true"
    echo "autorestart=true"
    echo "user=ubuntu"
    echo "stderr_logfile=/var/log/overwatch_slave.log"
    echo "stdout_logfile=/var/log/overwatch_slave.log"
} >> /home/ubuntu/overwatch_slave/oslave.conf

cp /home/ubuntu/overwatch_slave/oslave.conf /etc/supervisor/conf.d/oslave.conf

supervisorctl reread
supervisorctl update
