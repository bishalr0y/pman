# From repo root:
#   go build -o bin/pman .
#   vhs scripts/screenshot.tape
Output assets/screenshot.gif
Set Shell "bash"
Set FontFamily "Berkeley Mono"
Set FontSize 16
Set Columns 120
Set Rows 34
Set Theme "Catppuccin Mocha"
Set WindowBar Colorful
Set Padding 20
Set TypingSpeed 1ms

Env TERM "xterm-256color"
Env COLORTERM "truecolor"

Hide
Type "./bin/pman"
Enter
Sleep 5s
Show
Sleep 2s
Screenshot assets/screenshot.png
Type "q"
Enter
