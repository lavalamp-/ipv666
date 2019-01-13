package splash

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

var title = `
  /\\\\\\\\\\\  /\\\\\\\\\\\\\                            /\\\\\            /\\\\\            /\\\\\        
  \/////\\\///  \/\\\/////////\\\                      /\\\\////         /\\\\////         /\\\\////        
       \/\\\     \/\\\       \/\\\                   /\\\///           /\\\///           /\\\///            
        \/\\\     \/\\\\\\\\\\\\\/   /\\\    /\\\   /\\\\\\\\\\\      /\\\\\\\\\\\      /\\\\\\\\\\\        
         \/\\\     \/\\\/////////    \//\\\  /\\\   /\\\\///////\\\   /\\\\///////\\\   /\\\\///////\\\     
          \/\\\     \/\\\              \//\\\/\\\   \/\\\      \//\\\ \/\\\      \//\\\ \/\\\      \//\\\   
           \/\\\     \/\\\               \//\\\\\    \//\\\      /\\\  \//\\\      /\\\  \//\\\      /\\\   
         /\\\\\\\\\\\ \/\\\                \//\\\      \///\\\\\\\\\/    \///\\\\\\\\\/    \///\\\\\\\\\/   
         \///////////  \///                  \///         \/////////        \/////////        \/////////    
`
var backslashColor = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
var forwardslashColor = color.New(color.FgRed).SprintFunc()

func getPrintableSplash() string {
	toReturn := strings.Replace(title, "/", forwardslashColor("/"), -1)
	return strings.Replace(toReturn, "\\", backslashColor("\\"), -1)
}

func PrintSplash() {
	fmt.Println(getPrintableSplash())
}
