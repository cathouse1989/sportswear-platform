@echo off
setlocal
cd /d D:\code\sportswear-platform
echo START %DATE% %TIME% > up.log
docker compose -f docker-compose.yml up -d --force-recreate >> up.log 2>&1
echo EXITCODE:%ERRORLEVEL% >> up.log
docker compose -f docker-compose.yml ps >> up.log 2>&1
echo DONE > built.flag