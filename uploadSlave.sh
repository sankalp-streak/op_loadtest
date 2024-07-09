CGO_ENABLED=0 GOOS=linux go build -o slave ./slave/
echo build done
mv slave/slave slave/loadTestSlave/home/ubuntu/loadTest
echo move done
cd slave
dpkg-deb --build loadTestSlave
echo packagedone
sudo scp -i /mnt/c/Workspace/pem/streak-mumbai-deployment.pem -r /mnt/c/Workspace/op_loadtest/slave/loadTestSlave.deb  ubuntu@52.66.84.248:/home/ubuntu/services/loadTestSlave/