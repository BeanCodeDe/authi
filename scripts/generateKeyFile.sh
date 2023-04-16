ssh-keygen -t rsa -b 4096 -m PEM -f ./deployments/data/token/jwtRS256.key
openssl rsa -in ./deployments/data/token/jwtRS256.key -pubout -outform PEM -out ./deployments/data/token/jwtRS256.key.pub