# obfuscate-xor

A little toy to encode strings with other strings. The point of this is to hide a string "under" an expected string in a binary or script.

## Usage
```
Usage of ./obfuscate-xor:
  -language string
    	Language for output (default "go")
  -matches int
    	Number of matches (default 10)
  -plain string
    	Plain text to hide (default "AmsiScanBuffer")
  -wordlist string
    	Wordlist for cover text (default "words.txt")
```

## Go Example
```
$ ./obfuscate-xor
// Helper xor function:
func xor(input string, key []byte) (output string) {
    for i := 0; i < len(input); i++ {
        output += string(input[i] ^ key[i%len(key)])
    }
    return output
}
// Covers:
AmsiScanBuffer := xor("sysblocktraced", []byte{0x32, 0x14, 0x0, 0xb, 0x3f, 0xc, 0x2, 0x5, 0x36, 0x7, 0x7, 0x5, 0x0, 0x16}) // AmsiScanBuffer
```

## Powershell Example
```
$ ./obfuscate-xor -language powershell
# Helper xor function:
function xor {
    Param(
        [Parameter(Position = 0, Mandatory = $True)] $plain,
        [Parameter(Position = 1, Mandatory = $True)] $key
    )
    $r = ""
    $plaintext.ToCharArray() | foreach-object -process {
        $r += [char]([byte][char]$_ -bxor $key[$r.Length % $key.Length])
    }
    return $r
}
# Covers:
$AmsiScanBuffer = xor -plain "sysblocktraced" -key @(0x32,0x14,0x0,0xb,0x3f,0xc,0x2,0x5,0x36,0x7,0x7,0x5,0x0,0x16)
```
