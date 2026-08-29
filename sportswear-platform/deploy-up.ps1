Set-Location "D:\code\sportswear-platform"
$log = "D:\code\sportswear-platform\up.log"
"# starting up: $(Get-Date -Format o)" | Out-File -FilePath $log
Start-Process -FilePath "docker" -ArgumentList @("compose","-f","D:\code\sportswear-platform\docker-compose.yml","up","-d","--build") -NoNewWindow -RedirectStandardOutput $log -RedirectStandardError ($log + ".err")
"launched, pid captured"