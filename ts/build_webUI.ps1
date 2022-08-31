
$webUIHeader = @"
package src

var webUI = make(map[string]interface{})

func loadHTMLMap() {

"@

Remove-Item '../src/webUI.go'

Set-Content '../src/webUI.go' $webUIHeader

Get-ChildItem -Path "../html" -Recurse | 

Foreach-Object {
	if ($_.PsIsContainer -eq $False) {
		$fullPath = $_.FullName
		$htmlIndex = $_.FullName.Replace("\", "/").IndexOf("/html") + 1
		$file = $_.FullName.Replace("\", "/").Substring($htmlIndex)
		$fileContentBytes = Get-Content $fullPath -Encoding byte -Raw
		$fileContentEncoded = [System.Convert]::ToBase64String($fileContentBytes)
		$output = "`t" + 'webUI["' + $file + '"] = "' + $fileContentEncoded + '"'
		Add-Content '../src/webUI.go' $output
	}
}

Add-Content '../src/webUI.go' "`n}"
