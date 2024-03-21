systemctl stop teleport.service
cd ..
go build .
cp teleport.service /etc/systemd/system
systemctl daemon-reload
systemctl enable teleport.service
systemctl start teleport.service
echo "teleport.service started"
journalctl -u teleport.service
