mkdir ./deployments/token
ssh-keygen -t rsa -b 4096 -m PEM -f ./deployments/token/jwtRS256.key
openssl rsa -in ./deployments/token/jwtRS256.key -pubout -outform PEM -out ./deployments/token/jwtRS256.key.pub