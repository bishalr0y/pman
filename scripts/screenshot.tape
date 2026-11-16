# From repo root:
#   go build -ldflags "-X github.com/bishalr0y/pman/cmd.Version=6.0.0" -o bin/pman .
#   vhs scripts/screenshot.tape
Output assets/screenshot.gif
Set Shell "bash"
Set FontSize 18
Set Width 1920
Set Height 1080
Set Theme "Catppuccin Mocha"
Set WindowBar Colorful
Set TypingSpeed 50ms

Type "./bin/pman"
Enter
Sleep 5s
# Toggle UDP twice so the status line shows "Showing TCP only"
Type "u"
Sleep 400ms
Type "u"
Sleep 400ms
Screenshot assets/screenshot.png
Type "q"
Enter