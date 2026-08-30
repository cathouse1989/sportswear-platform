# ============================================================
#  Sportswear 一键更新脚本（用于后续代码更新）
#  用法：右键 -> "使用 PowerShell 运行"，或在终端执行：
#      powershell -ExecutionPolicy Bypass -File update.ps1
# ============================================================

$ErrorActionPreference = "Stop"

# ------------------------------------------------------------
# 兼容 Windows PowerShell 5.1：
# 5.1 会把外部命令(git/docker 等)写入 stderr 的正常输出也看作
# "错误"，在 $ErrorActionPreference="Stop" 下会直接终止脚本并报
# NativeCommandError。下面两个函数用于保证外部命令正常运行：
#   - Invoke-Native      ：执行外部命令，期间临时放宽错误策略，
#                          把 stderr 作为普通文本显示，返回退出码
#   - Get-ContainerHealth：只读回容器的健康状态(丢弃 stderr)
# ------------------------------------------------------------
function Invoke-Native {
    param([scriptblock]$ScriptBlock)
    $savedEAP = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
        & $ScriptBlock 2>&1 | ForEach-Object { "$_" } | Write-Host
        return $LASTEXITCODE
    }
    finally {
        $ErrorActionPreference = $savedEAP
    }
}

function Get-ContainerHealth([string]$name) {
    $savedEAP = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    $out = (& docker inspect --format='{{.State.Health.Status}}' $name 2>$null) -join ''
    $ErrorActionPreference = $savedEAP
    return $out.Trim()
}

$projectRoot = $PSScriptRoot
$composeFile  = "$projectRoot\docker-compose.yml"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Sportswear 一键更新部署"                -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 拉取最新代码 -------------------------------------------------
Write-Host "[1/4] 拉取最新代码..." -ForegroundColor Yellow
Set-Location $projectRoot
$gitExitCode = Invoke-Native { & git pull origin main }
if ($gitExitCode -ne 0) {
    Write-Host "[WARNING] git pull 失败，将使用本地代码继续" -ForegroundColor DarkYellow
} else {
    Write-Host "  代码已更新" -ForegroundColor Green
}

# 2. 重新构建并启动所有服务 ---------------------------------------
Write-Host ""
Write-Host "[2/4] 重新构建镜像并启动容器..." -ForegroundColor Yellow

$dcExitCode = Invoke-Native { & docker compose -p sportswear -f $composeFile up -d --build }
if ($dcExitCode -ne 0) {
    Write-Host "[ERROR] 更新失败，请检查错误信息" -ForegroundColor Red
    pause
    exit 1
}
Write-Host "  构建+启动完成" -ForegroundColor Green

# 3. 等待健康检查 -------------------------------------------------
Write-Host ""
Write-Host "[3/4] 等待服务就绪..." -ForegroundColor Yellow
$retry = 0
do {
    Start-Sleep -Seconds 2
    $healthy = Get-ContainerHealth "sportswear-postgres"
    $retry++
} while ($healthy -ne 'healthy' -and $retry -lt 15)
if ($healthy -eq 'healthy') {
    Write-Host "  PostgreSQL 就绪" -ForegroundColor Green
}

# 4. 显示运行状态 -------------------------------------------------
Write-Host ""
Write-Host "[4/4] 容器运行状态：" -ForegroundColor Yellow
Write-Host ""
docker ps -a --filter "name=sportswear" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# 访问地址 --------------------------------------------------------
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  更新完成！访问地址："                   -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  后端 API   : http://localhost:8080"    -ForegroundColor White
Write-Host "  管理后台   : http://localhost:8090"    -ForegroundColor White
Write-Host "  门户前端   : http://localhost:3000"    -ForegroundColor White
Write-Host "  MinIO      : http://localhost:9001"    -ForegroundColor White
Write-Host ""
Write-Host "  默认管理员: admin@sportswear.com / Admin@123456" -ForegroundColor Gray
Write-Host "  MinIO账号 : minioadmin / minioadmin"              -ForegroundColor Gray
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "提示：如需查看日志，运行：" -ForegroundColor DarkYellow
Write-Host "  docker compose -p sportswear -f `"$composeFile`" logs -f" -ForegroundColor DarkGray
Write-Host ""

pause
