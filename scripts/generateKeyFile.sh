mkdir ./deployments/token
mkdir ./deployments/token/privat
mkdir ./deployments/token/public
ssh-keygen -t rsa -b 4096 -m PEM -f ./deployments/token/privat/jwtRS256.key
openssl rsa -in ./deployments/token/privat/jwtRS256.key -pubout -outform PEM -out ./deployments/token/public/jwtRS256.key.pub