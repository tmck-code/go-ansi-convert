[![.github/workflows/test.yml](https://github.com/tmck-code/go-ansi-convert/actions/workflows/test.yml/badge.svg)](https://github.com/tmck-code/go-ansi-convert/actions/workflows/test.yml)

# go-ansi-convert

A tool to convert ANSI images & files in various ways.

- Sanitize
- Justify
- Flip

https://github.com/user-attachments/assets/06123c36-fbbf-44f0-b852-c628d9c69aef

```
Usage: ansi-flip [-dhjsx] [--display-separator value] [--display-separator-width value] [-f value] [-i value] [-o value] [parameters ...]
 -d, --display      Display original and flipped side-by-side in terminal
     --display-separator=value
                    Separator string between original and flipped when
                    displaying [ ]
     --display-separator-width=value
                    Width of separator between original and flipped when
                    displaying [1]
 -f, --flip=value   Flip horizontally (h), vertically (v), or both (h,v or v,h)
                    {operation}
 -h, --help         display this help message {operation}
 -i, --input=value  Input file path (default: stdin)
 -j, --justify      Justify lines to the same length (sanitise mode only)
 -o, --output=value
                    Output file path (default: stdout)
 -s, --sanitise     Sanitise ANSI lines, ensuring that each line ends with a
                    reset code {operation}
 -x, --display-swapped
                    When displaying, reverse the order of original and flipped
```
