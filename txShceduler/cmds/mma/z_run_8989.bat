@echo off


set port=8989
title TXSCHEDULER_MMA V1.0.0  [ %port% ]

call txmMma.exe port=%port% start

pause