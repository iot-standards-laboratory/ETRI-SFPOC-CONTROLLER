import os 
import sys

os.system('go build -o controller-app main.go')
serviceTemplate = '''[Unit]
Description=Run Smart Farm PoC Controller
After=multi-user.target

[Service]
Type=simple
WorkingDirectory={}
ExecStart={}/controller-app
Restart=always
RestartSec=10s

[Install]
WantedBy=multi-user.target
'''

curDir = os.getcwd()
service = serviceTemplate.format(curDir, curDir)

with open('sfpoc.service', 'w') as f:
    f.write(service)

os.system('sudo chmod 644 sfpoc.service')
os.system('sudo mv sfpoc.service /lib/systemd/system/')
os.system('sudo systemctl daemon-reload')
os.system('sudo systemctl enable sfpoc.service')


