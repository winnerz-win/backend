@echo off

set port=60002

title ETH_TX_SIGN_SERVER [%port%]

call tx_signer_eth.exe --port %port% --config ./config.yaml

pause