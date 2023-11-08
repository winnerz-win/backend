@echo off


set port=80
title TXSCHEDULER_MMA V1.0.0  [ %port% ]

call txmMma.exe port=%port% start

pause