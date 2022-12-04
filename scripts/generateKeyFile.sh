mkdir ./deployments/token
mkdir ./deployments/token/privat
mkdir ./deployments/token/public
ssh-keygen -t rsa -P "" -b 4096 -m PEM -f ./deployments/token/privat/jwtRS256.key
ssh-keygen -e -m PEM -f ./deployments/token/privat/jwtRS256.key > ./deployments/token/public/jwtRS256.key.pub