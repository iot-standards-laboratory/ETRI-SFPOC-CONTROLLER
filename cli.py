import os 

os.system("rm *.db")
os.system("rm config.properties")

os.system("rm main")
os.system("go build main.go")
os.system("sudo chmod +s main")

os.system("./main -init")
