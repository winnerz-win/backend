@echo off


title ETH_TX_SIGN_SERVER_BUILD

echo build try...

go build -o tx_signer_eth.exe main.go

echo build end...