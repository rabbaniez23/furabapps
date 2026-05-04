$services = Get-ChildItem -Path d:\Pekerjaan\furabapps\furab-backend\services -Directory
$failedCount = 0

foreach ($s in $services) {
    Write-Host "`n=================================================="
    Write-Host "Running tests for $($s.Name)..."
    Set-Location $s.FullName

    # Run Unit Tests
    Write-Host ">>> Unit Tests:"
    $unitTestResult = go test ./test/unit/... -v 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[FAIL] Unit tests failed in $($s.Name)" -ForegroundColor Red
        $failedCount++
    } else {
        Write-Host "[PASS] Unit tests passed" -ForegroundColor Green
    }

    # Run Functional Tests
    if (Test-Path "test/functional") {
        Write-Host ">>> Functional Tests:"
        $funcTestResult = go test ./test/functional/... -v -tags=functional 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Host "[FAIL] Functional tests failed in $($s.Name)" -ForegroundColor Red
            $failedCount++
        } else {
            Write-Host "[PASS] Functional tests passed" -ForegroundColor Green
        }
    }
}

Set-Location d:\Pekerjaan\furabapps\furab-backend
if ($failedCount -eq 0) {
    Write-Host "`n[SUCCESS] ALL SERVICES PASSED PERFECTLY! (100%)" -ForegroundColor Green
} else {
    Write-Host "`n[ERROR] $failedCount test suites failed." -ForegroundColor Red
}
