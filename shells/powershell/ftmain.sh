
function ft {
    param (
        [string[]]$args
    )
                    
    $temp_output = New-TemporaryFile
                        
    & $env:FT_EXE_PATH $args | Tee-Object -FilePath $temp_output
                            
    $output = Get-Content $temp_output | Select-Object -Last 1
                                
    if (Test-Path -Path $output -PathType Container) {
        Set-Location -Path $output
    }
                                                    
    Remove-Item $temp_output
}

