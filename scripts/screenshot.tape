# Run from repo root: go build -o bin/pman . && vhs scripts/screenshot.tape
Output assets/screenshot.gif
Set Shell "bash"
Set FontSize 14
Set Width 1280
Set Height 720
Set Theme "Catppuccin Mocha"

Type "./bin/pman"
Enter
Sleep 4s
Screenshot assets/screenshot.png
Type "q"
Enter