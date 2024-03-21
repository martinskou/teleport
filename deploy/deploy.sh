cd ..
go build .
cd deploy
systemctl stop teleport.service
cp teleport.service /etc/systemd/system
systemctl daemon-reload
systemctl enable teleport.service
systemctl start teleport.service
echo "teleport.service started"
journalctl -u teleport.service
