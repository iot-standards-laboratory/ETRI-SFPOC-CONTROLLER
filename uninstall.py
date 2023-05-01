
import os

os.system('sudo systemctl disable sfpoc.service')
os.system('sudo rm /lib/systemd/system/sfpoc.service')