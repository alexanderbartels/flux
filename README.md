## Description
Fluent Regular Expressions in golang is a port of https://github.com/selvinortiz/flux *by* [Selvin Ortiz](http://twitter.com/selvinortiz)
   
## Requirements
   * golang
      * Created and Tested with Version 1.4.2
   
## Install 
   Flux is available via `go get github.com/alexanderbartels/flux`

## Example 
This simple example illustrates the way you would use flux and it's fluent interface to build complex patterns.

```go
package main

import (
   "flux"
)

func main() {
   // create the regex 
   regex := flux.NewFlux().StartOfLine().Find("http").Maybe("s").Then("://").Maybe("www.").AnythingBut(".").Either(".co", ".com", ".de").IgnoreCase().EndOfLine()
   
   // print the created regex
   fmt.Println("Regex: ", regex.String())
   // or
   fmt.Println("Regex: ", regex)
   
   // The subject string (URL)
   subject := 'http://www.selvinortiz.com';
   
   // Match 
   match, _ := regex.Match("http://selvinortiz.com")
   if match {
      fmt.Println("Matched")
   }
   
   // Replace
   repl := regex.Replace("http://selvinortiz.com", "$5$6")
   if repl == "selvinortiz.com" {
      fmt.Println("Replacement ok")
   }
}
```
For other examples, please see the flux_test.go file.

## License
Flux is released under the [MIT license](http://opensource.org/licenses/MIT)


## TODO 
 * Add CI Builds
 * API Documentation
 * SourceCode Documentation
 * Branching, Changelog, Versioning / using [semver](http://semver.org/)
